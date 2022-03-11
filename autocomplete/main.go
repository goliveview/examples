package main

import (
	"log"
	"net/http"
	"strings"

	glv "github.com/goliveview/controller"
)

var cities = []string{
	"Paris",
	"Amsterdam",
	"Berlin",
	"New York",
	"Delhi",
	"Beijing",
	"London",
	"Rome",
	"Athens",
	"Seoul",
}

func getCities(str string) []string {
	if str == "" {
		return nil
	}
	var result []string
	for _, city := range cities {
		if strings.HasPrefix(strings.ToLower(city), strings.ToLower(str)) {
			result = append(result, city)
		}
	}
	return result
}

type QueryRequest struct {
	Query string `json:"query"`
}

type Search struct {
	glv.DefaultView
}

func (s *Search) Content() string {
	return "app.html"
}

func (s *Search) Partials() []string {
	return []string{"cities.html"}
}

func (s *Search) OnMount(_ glv.Context) (glv.Status, glv.M) {
	return glv.Status{Code: 200}, nil
}

func (s *Search) OnLiveEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "autocomplete":
		req := new(QueryRequest)
		if err := ctx.Event().DecodeParams(req); err != nil {
			return err
		}
		ctx.DOM().Morph("#list_cities", "cities", glv.M{
			"cities": getCities(req.Query),
		})
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func main() {
	glvc := glv.Websocket("goliveview-autocomplete", glv.DevelopmentMode(true))
	http.Handle("/", glvc.Handler(&Search{}))
	log.Println("listening on http://localhost:9867")
	http.ListenAndServe(":9867", nil)
}
