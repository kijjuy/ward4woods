package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
)

type StaticHandler struct {
	staticPath string
	template   string
	error      string
	logger     *slog.Logger
}

func NewStaticHandler(staticPath string, template string, error string, logger *slog.Logger) *StaticHandler {
	return &StaticHandler{
		staticPath: staticPath,
		template:   template,
		error:      error,
		logger:     logger.With("Location", "StaticHandler"),
	}
}

func (sh *StaticHandler) Render(w http.ResponseWriter, page string) {
	page = filepath.Join(sh.staticPath, page)
	tmpl, err := template.ParseFiles(sh.template, page)
	if err == nil {
		tmpl.Execute(w, nil)
		sh.logger.Info(fmt.Sprintf("Now serving page: %s", page))
		return
	}

	message := "Could not find page: " + page + ". "
	page = path.Join(page, "index.html")
	message += "Checking: " + page
	sh.logger.Info(fmt.Sprintf("%s Error: %s", message, err))

	tmpl, err = template.ParseFiles(sh.template, page)
	if err == nil {
		tmpl.Execute(w, nil)
		sh.logger.Info(fmt.Sprintf("Now serving page: %s", page))
		return
	}

	sh.logger.Info(fmt.Sprintf("Could not find page. Loading error view. Error: %s", err))
	tmpl, err = tmpl.ParseFiles(sh.template, sh.error)
	if err != nil {
		sh.logger.Error(fmt.Sprintf("Could not find error page. Error: '%s'.", err))
		sh.logger.Error(fmt.Sprintf("Attempted to find template page at: %s. Attempted to find error page at %s.", sh.template, sh.error))
		return
	}
	tmpl.Execute(w, nil)
	return
}

func (sh *StaticHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := filepath.Clean(r.URL.String())
	sh.Render(w, url)
}
