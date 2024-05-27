package helpers

import (
	"html/template"
	"log/slog"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	logger := slog.With("Location", "RenderTemplate")
	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		logger.Error("Could not load template files.", "Error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		logger.Error("Could not execute template.", "Error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
