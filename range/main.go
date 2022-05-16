package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	glv "github.com/goliveview/controller"
)

type CountRequest struct {
	Count string `json:"count"`
}

type Range struct {
	glv.DefaultView
}

func (r *Range) Content() string {
	return "app.html"
}

func (r *Range) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	return glv.Status{Code: 200}, glv.M{
		"val": 0,
	}
}

func (r *Range) OnLiveEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "update_item_count":
		req := new(CountRequest)
		if err := ctx.Event().DecodeParams(req); err != nil {
			return err
		}
		ctx.DOM().SetInnerHTML("#count", fmt.Sprintf(`<p>You have %v items</p>`, req.Count))
		count, err := strconv.Atoi(req.Count)
		if err != nil {
			return err
		}
		ctx.DOM().SetInnerHTML("#total", fmt.Sprintf(`<p>$ %v</p>`, count*10))
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func main() {
	glvc := glv.Websocket("goliveview-range", glv.DevelopmentMode(true))
	http.Handle("/", glvc.Handler(&Range{}))
	log.Println("listening on http://localhost:9867")
	http.ListenAndServe(":9867", nil)
}
