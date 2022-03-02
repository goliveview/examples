package main

import (
	"log"
	"net/http"

	glv "github.com/goliveview/controller"
)

type LayoutView struct {
	glv.DefaultView
}

func (l *LayoutView) Layout() string {
	return "templates/layouts/index.html"
}

type HomeView struct {
	LayoutView
}

func (h *HomeView) Content() string {
	return "templates/views/home"
}

type HelpView struct {
	LayoutView
}

func (h *HelpView) Content() string {
	return "templates/views/help"
}

type SettingsView struct {
	LayoutView
}

func (h *SettingsView) Content() string {
	return "templates/views/settings"
}

func main() {
	glvc := glv.Websocket("goliveview-layout", glv.DevelopmentMode(true))
	http.Handle("/", glvc.Handler(&HomeView{}))
	http.Handle("/help", glvc.Handler(&HelpView{}))
	http.Handle("/settings", glvc.Handler(&SettingsView{}))
	log.Println("listening on http://localhost:9867")
	http.ListenAndServe(":9867", nil)
}
