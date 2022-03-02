package main

import (
	"net/http"

	glv "github.com/goliveview/controller"
)

type SimpleView struct {
	glv.DefaultView
}

func (s *SimpleView) Content() string {
	return `{{define "content"}}<div>world</div>{{ end }}`
}

func (s *SimpleView) Layout() string {
	return `<div>Hello: {{template "content" .}}</div>`
}

func main() {
	glvc := glv.Websocket("goliveview-simple", glv.DevelopmentMode(true))
	http.Handle("/", glvc.Handler(&SimpleView{}))
	http.ListenAndServe(":9867", nil)
}
