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
	"Application":    "Application",
	"Home":           "Home",
	"No content.":    "No content.",
	"Index Page":     "Index Page",
	"Home Page":      "Home Page",
	"Hello":          "Hello",
	"Hello, {name}!": "Hello, %s!",
}

var es = map[string]string{
	"Home Page":      "Página de inicio",
	"Hello":          "Hola",
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
	vars map[string]string
	t    map[string]string
}

// Var returns a variaable from the template data by name.
func (td *TemplateData) Var(key string) string {
	return td.vars[key]
}

// SetVar updates the value of a variable in the template.
func (td *TemplateData) SetVar(key, value string) {
	td.vars[key] = value
}

// MergeTranslations adds a new translation to this template data map.
func (td *TemplateData) MergeTranslations(tmap map[string]string) {
	mergeMaps(td.t, tmap)
}

// T fetches a translation by key from the map
func (td *TemplateData) T(key string) string {
	return td.t[key]
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
		data.SetVar("name", chi.URLParam(r, "name"))

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
		data.MergeTranslations(en)
		locale := r.URL.Query().Get("locale")
		if locale != "en" {
			if translations, ok := localeMap[locale]; ok {
				data.MergeTranslations(translations)
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
			vars: make(map[string]string),
			t:    make(map[string]string),
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
