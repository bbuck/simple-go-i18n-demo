package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dannyvankooten/extemplate"
	"github.com/go-chi/chi"
)

var en = map[string]string{
	"Index Page":     "Index Page",
	"Home Page":      "Home Page",
	"Hello, {name}!": "Hello, %s!",
}

var es = map[string]string{
	"Hello, {name}!": "¡Hola, %s!",
}

var localeMap = map[string]map[string]string{
	"en": en,
	"es": es,
}

type templateKey int

var templateDataKey = templateKey(1)

var xt *extemplate.Extemplate

// TemplateData is data for the templates, containing vars and translation
// mappings.
type TemplateData struct {
	Vars map[string]string
	T    map[string]string
}

func init() {
	xt = extemplate.New().Funcs(template.FuncMap{
		"format": fmt.Sprintf,
	})
	err := xt.ParseDir("templates/", []string{".tmpl"})
	if err != nil {
		panic(err)
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(InitTemplateData)
	r.Use(InitI18N)
	r.Use(SetContentType("text/html"))

	r.Get("/", TemplateHandler("index.tmpl"))
	r.Get("/home", TemplateHandler("views/home.tmpl"))
	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		data := getTemplateData(r.Context())
		data.Vars["name"] = chi.URLParam(r, "name")

		renderTemplate("views/hello.tmpl", w, data)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

// InitI18N populates templateData.T with values
func InitI18N(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := getTemplateData(r.Context())
		// base language data
		mergeMaps(data.T, en)
		locale := r.URL.Query().Get("locale")
		if locale != "en" {
			if translations, ok := localeMap[locale]; ok {
				mergeMaps(data.T, translations)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// SetContentType is a middleware to set the content type of requests nested
// under it's router.
func SetContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", contentType)
			next.ServeHTTP(w, r)
		})
	}
}

// InitTemplateData ensures that an empty mpa is set to template data for all
// requests.
func InitTemplateData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := &TemplateData{
			Vars: make(map[string]string),
			T:    make(map[string]string),
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, templateDataKey, data)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

// TemplateHandler is a handler that simply renders a template with current
// contextual template data (usually just containing translations)
func TemplateHandler(templateName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		templateData := getTemplateData(r.Context())

		renderTemplate(templateName, w, templateData)
	}
}

func renderTemplate(name string, w http.ResponseWriter, data *TemplateData) {
	err := xt.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Printf("[ERROR] %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getTemplateData(ctx context.Context) *TemplateData {
	return ctx.Value(templateDataKey).(*TemplateData)
}

func mergeMaps(dst, src map[string]string) {
	for key, val := range src {
		dst[key] = val
	}
}
