package frontend

import (
	"fileupbackendv2/config"
	"fileupbackendv2/endpoints"
	"fileupbackendv2/internal/storage"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// type Template struct {
// templates *template.Template
// }?

func renderTemplate(w io.Writer, name string, data interface{}, c echo.Context, only bool) error {
	baseTemps := []string{
		"frontend/templates/base.html",
		"frontend/templates/sidebar.html",
	}
	if only {
		baseTemps = []string{
			"frontend/templates/" + name + ".html",
		}
	} else {
		baseTemps = append(baseTemps, "frontend/templates/"+name+".html")
	}
	tmpl, err := template.ParseFiles(baseTemps...)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func StartEchoServer() {
	e := echo.New()
	subgroup := e.Group("/app")

	e.Debug = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(CookieMiddleware)

	subgroup.GET("/hello/", hello)
	subgroup.GET(endpoints.FsPaths.ApiKeyForm, apiKeyForm)
	subgroup.POST(endpoints.FsPaths.ApiKeyVerify, verifyApiKey)
	subgroup.GET("/uploadfiles/", uploadFileForm)
	subgroup.GET("/addfolder/", addFolderForm)
	subgroup.GET("/", listFiles)
	e.Static("/static/", "frontend/static")
	e.Logger.Fatal(e.Start(":" + config.FrontEndPort))
}

func hello(c echo.Context) error {
	return renderTemplate(c.Response(), "hello", nil, c, false)
}

func apiKeyForm(c echo.Context) error {
	return renderTemplate(c.Response(), "apikey", nil, c, true)
}

func verifyApiKey(c echo.Context) error {
	apiKey := c.FormValue("apikey")
	_, err := storage.GetUserByKey(apiKey)
	if err != nil {
		return c.Render(http.StatusBadRequest, "apikey", err.Error())
	}
	cookie := &http.Cookie{
		Name:     "api-key",
		Value:    apiKey,
		Expires:  time.Now().Add(24 * 10 * time.Hour),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	return c.HTML(http.StatusOK, `<p>✅ Login successful!</p> <script>window.location.href='/app/hello/'</script>`)
}
