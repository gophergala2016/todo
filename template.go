package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type AppTemplate struct {
	t *template.Template
}

func NewAppTemplate(files ...string) *AppTemplate {
	base := template.Must(template.New("base").ParseFiles("base.html"))
	template.Must(base.ParseFiles(files...))
	return &AppTemplate{base}
}

func (a AppTemplate) Execute(w http.ResponseWriter, data interface{}) *appError {
	d := struct {
		Data interface{}
	}{
		Data: data,
	}
	if err := a.t.Execute(w, d); err != nil {
		return InternalServerError(fmt.Errorf("execute template: %v", err))
	}
	return nil
}
