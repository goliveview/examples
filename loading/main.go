package main

import (
	"log"
	"net/http"
	"time"

	glv "github.com/goliveview/controller"
)

type Loading struct {
	glv.DefaultView
}

func (l *Loading) Content() string {
	return "app.html"
}

func (l *Loading) OnLiveEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "show_loader":
		ctx.DOM().AddClass("#loading-modal", "is-active")
		defer func() {
			ctx.DOM().RemoveClass("#loading-modal", "is-active")
		}()

		// some work
		time.Sleep(time.Second * 2)
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func main() {
	glvc := glv.Websocket("goliveview-counter", glv.DevelopmentMode(true))
	http.Handle("/", glvc.Handler(&Loading{}))
	log.Println("listening on http://localhost:9867")
	http.ListenAndServe(":9867", nil)
}
