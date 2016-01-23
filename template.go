package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
)

func add(a, b int) int {
	return a + b
}

func float64ToTime(v float64) time.Time {
	return time.Unix(int64(v), 0)
}

func float64ToHuman(v float64) string {
	return humanize.Time(float64ToTime(v))
}

var funcMap = template.FuncMap{
	"add":            add,
	"float64ToTime":  float64ToTime,
	"float64ToHuman": float64ToHuman,
}

type AppTemplate struct {
	t *template.Template
}

func NewAppTemplate(files ...string) *AppTemplate {
	base := template.Must(template.New("base").Funcs(funcMap).ParseFiles("base.html"))
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
