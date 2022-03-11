package main

import (
	"log"
	"net/http"
	"time"

	glv "github.com/goliveview/controller"
)

func NewTimer() *Timer {
	timerCh := make(chan glv.Event)
	ticker := time.NewTicker(time.Second)
	go func() {
		for ; true; <-ticker.C {
			timerCh <- glv.Event{ID: "tick"}
		}
	}()
	return &Timer{ch: timerCh}
}

type Timer struct {
	glv.DefaultView
	ch chan glv.Event
}

func (t *Timer) Content() string {
	return "app.html"
}

func (t *Timer) Partials() []string {
	return []string{"time.html"}
}

func (t *Timer) OnMount(ctx glv.Context) (glv.Status, glv.M) {
	return glv.Status{Code: 200}, glv.M{
		"val": time.Now().String(),
	}
}

func (t *Timer) OnLiveEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "tick":
		ctx.DOM().Morph("#time", "time", glv.M{
			"val": time.Now().String(),
		})
		return nil
	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func (t *Timer) LiveEventReceiver() <-chan glv.Event {
	return t.ch
}

func main() {
	glvc := glv.Websocket("goliveview-timer", glv.DevelopmentMode(true))
	http.Handle("/", glvc.Handler(NewTimer()))
	log.Println("listening on http://localhost:9867")
	http.ListenAndServe(":9867", nil)
}
