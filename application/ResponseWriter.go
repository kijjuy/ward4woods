package application

import (
	"fmt"
	"html/template"
	"net/http"
)

type ResponseWriter func(http.ResponseWriter, interface{}) error

type ErrorWriter func(http.ResponseWriter, error) error

type TemplateWriter struct {
	templatePath string
	errorPath    string
}

func NewApiTemplateWriter(templatePath string) *TemplateWriter {
	return &TemplateWriter{templatePath: templatePath}
}

func NewTemplateWriter(templatePath, errorPath string) *TemplateWriter {
	return &TemplateWriter{templatePath: templatePath, errorPath: errorPath}
}

func writeTemplate(w http.ResponseWriter, templatePath string, data interface{}) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

func (tw *TemplateWriter) WriteTemplate(w http.ResponseWriter, data interface{}) error {
	return writeTemplate(w, tw.templatePath, data)
}

func (tw *TemplateWriter) WriteErrorTemplate(w http.ResponseWriter, err error) error {
	if tw.errorPath == "" {
		return fmt.Errorf("Template Writer error path not set.")
	}
	err = writeTemplate(w, tw.errorPath, err)
	return err
}

func WriteServerError(w http.ResponseWriter, err error) error {
	fmt.Fprintf(w, "Internal Server Error")
	return nil
}
