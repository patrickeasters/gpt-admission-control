package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickeasters/gpt-admission-webhook/handlers"
	"github.com/sashabaranov/go-openai"
)

func main() {
	e := echo.New()

	key := os.Getenv("OPENAI_API_KEY")
	if len(key) == 0 {
		e.Logger.Fatal("no OPENAI_API_KEY provided")
	}

	c := openai.NewClient(key)
	h := handlers.Handlers{
		OpenAIClient: c,
	}

	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🛂")
	})
	e.POST("/admission", h.AdmissionWebhook)

	if useTLS, _ := strconv.ParseBool(os.Getenv("TLS_ENABLED")); useTLS {
		e.Logger.Fatal(e.StartTLS(":8443", "/tls/server.crt", "/tls/server.key"))
	} else {
		e.Logger.Fatal(e.Start(":3000"))
	}

}
