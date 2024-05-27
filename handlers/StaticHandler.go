package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
)

// A StaticHandler is meant to handle requests that require an entire page to be loaded. These pages will be
// constructed by the html/template package and returned to the browser.
type StaticHandler struct {
	staticPath string
	template   string
	error      string
	logger     *slog.Logger
}

// NewStaticHandler creates a StaticHandler struct and returns a pointer to it.
func NewStaticHandler(staticPath string, template string, error string, logger *slog.Logger) *StaticHandler {
	return &StaticHandler{
		staticPath: staticPath,
		template:   template,
		error:      error,
		logger:     logger.With("Location", "StaticHandler"),
	}
}

// Render method attempts to write a templated view to the http.ResponseWriter. It takes a string as the path to the page
// to load into the template, and attempts to load that template into the StaticHandler.template field.
func (sh *StaticHandler) Render(w http.ResponseWriter, page string, data interface{}) {
	page = filepath.Join(sh.staticPath, page)

	err := sh.tryServePage(page, w, data)
	if err == nil {
		return
	}

	sh.logger.Warn("Error Rendering Page:", "Error", err)

	page = sh.nextPage(page, "index.html")

	if err := sh.tryServePage(page, w, data); err == nil {
		return
	}

	sh.logger.Info(fmt.Sprintf("Could not find page. Loading error view. Error: %s", err))
	errorPage := filepath.Join(sh.staticPath, "error.html")
	err = sh.tryServePage(errorPage, w, nil)
	if err != nil {
		sh.logger.Error(fmt.Sprintf("Could not find error page. Error: '%s'.", err))
		sh.logger.Error(fmt.Sprintf("Attempted to find template page at: %s. Attempted to find error page at %s.", sh.template, sh.error))
		return
	}
	return
}

// tryServePage attempts to execute the template based on the page it is provided.
// Returns early with an error if there is any trouble loading the page.
func (sh *StaticHandler) tryServePage(page string, w http.ResponseWriter, Data interface{}) error {
	tmpl, err := template.New("content").ParseFiles(page)
	if err != nil {
		sh.logger.Error("Error parsing content template.", "Error", err)
		return err
	}

	tmpl, err = tmpl.ParseFiles(sh.template)
	if err != nil {
		sh.logger.Error("Error parsing layout template.", "Error", err)
		return err
	}

	fmt.Printf("Data is: %+v\n", Data)
	if err := tmpl.ExecuteTemplate(w, "layout", Data); err != nil {
		sh.logger.Error("Error executing layout template.", "Error", err)
		return err
	}

	sh.logger.Info(fmt.Sprintf("Now serving page: %s", page))
	return nil
}

// nextPage adds a specified suffix to the current page and returns it.
// Also logs that the page could not be found and is attempting to find it at another location.
func (sh *StaticHandler) nextPage(page, nextLocation string) string {
	message := "Could not find page: " + page + ". "
	page = path.Join(page, nextLocation)
	message += "Checking: " + page

	sh.logger.Info(message)
	return page
}

// HandleRequests is the default request handler for static file handling.
// It takes the path from the url and checks sends it to the Render method.
func (sh *StaticHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := filepath.Clean(r.URL.String())
	sh.Render(w, url, nil)
}
