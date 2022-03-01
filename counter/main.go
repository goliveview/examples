package main

import (
	"fmt"
	"log"
	"net/http"

	glv "github.com/goliveview/controller"
)

type Counter struct {
	glv.DefaultView
}

func (c *Counter) Content() string {
	return "app.html"
}

func (c *Counter) Partials() []string {
	return []string{"count.html"}
}

func (c *Counter) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	return glv.Status{Code: 200}, glv.M{
		"val": 0,
	}
}

func (c *Counter) OnEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "inc":
		var val int
		ctx.Store().Get("val", &val)
		val += 1
		ctx.Store().Put(glv.M{"val": val})
		ctx.DOM().SetInnerHTML("#count", fmt.Sprintf(`<p>%v</p>`, val))
	case "dec":
		var val int
		ctx.Store().Get("val", &val)
		val -= 1
		ctx.Store().Put(glv.M{"val": val})
		ctx.DOM().SetInnerHTML("#count", fmt.Sprintf(`<p>%v</p>`, val))
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func main() {
	glvc := glv.Websocket("goliveview-counter", glv.DevelopmentMode(true))
	http.Handle("/", glvc.Handler(&Counter{}))
	log.Println("listening on http://localhost:9867")
	http.ListenAndServe(":9867", nil)
}
