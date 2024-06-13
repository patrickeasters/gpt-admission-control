package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/patrickeasters/gpt-admission-webhook/gpt"
	"github.com/sashabaranov/go-openai"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Handlers struct {
	OpenAIClient *openai.Client
}

func (h *Handlers) AdmissionWebhook(c echo.Context) error {
	var req v1.AdmissionReview
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.Request == nil {
		c.Logger().Info("No admission request present")
		return c.NoContent(http.StatusBadRequest)
	}

	out := v1.AdmissionReview{
		TypeMeta: req.TypeMeta,
		Response: &v1.AdmissionResponse{
			UID:    req.Request.UID,
			Result: &metav1.Status{},
		},
	}

	// what kind of day is it?
	serializedObject, err := json.Marshal(req.Request.Object)
	if err != nil {
		c.Logger().Errorf("Failed to serialize object from admission review: %s", err)
	}
	decision, err := gpt.Decide(c.Request().Context(), h.OpenAIClient, string(serializedObject))
	if err != nil {
		out.Response.Allowed = false
		out.Response.Result.Code = http.StatusInternalServerError
		out.Response.Result.Message = "Sorry! The AI wasn't cooperating."
	} else if decision.Admitted {
		out.Response.Allowed = true
		out.Response.Result.Code = http.StatusOK
		out.Response.Result.Message = decision.Reason
	} else {
		out.Response.Allowed = false
		out.Response.Result.Code = http.StatusForbidden
		out.Response.Result.Message = decision.Reason
	}

	return c.JSON(http.StatusOK, out)
}
