package main

import (
	"log"
	"net/http"
	"time"

	"github.com/timshannon/bolthold"

	glv "github.com/goliveview/controller"
)

type Todo struct {
	ID        uint64 `json:"id" boltholdKey:"ID"`
	Text      string `json:"text"`
	Done      bool   `json:"done"`
	CreatedAt time.Time
}

type TodosView struct {
	glv.DefaultView
	db *bolthold.Store
}

func (c *TodosView) Content() string {
	return "app.html"
}

func (c *TodosView) Partials() []string {
	return []string{"todos.html"}
}

func (c *TodosView) OnMount(w http.ResponseWriter, r *http.Request) (glv.Status, glv.M) {
	var todos []Todo
	if err := c.db.Find(&todos, &bolthold.Query{}); err != nil {
		return glv.Status{
			Code: 200,
		}, nil
	}
	return glv.Status{Code: 200}, glv.M{"todos": todos}
}

func (c *TodosView) OnEvent(ctx glv.Context) error {
	switch ctx.Event().ID {
	case "todos/new":
		var todo Todo
		if err := ctx.Event().DecodeParams(&todo); err != nil {
			return err
		}
		if err := c.db.Insert(bolthold.NextSequence(), &todo); err != nil {
			return err
		}
		// list updated todos
		var todos []Todo
		if err := c.db.Find(&todos, &bolthold.Query{}); err != nil {
			return err
		}
		ctx.DOM().Morph("#todos", "todos", glv.M{"todos": todos})
	case "todos/del":
		var todo Todo
		if err := ctx.Event().DecodeParams(&todo); err != nil {
			return err
		}
		if err := c.db.Delete(todo.ID, &todo); err != nil {
			return err
		}
		// list updated todos
		var todos []Todo
		if err := c.db.Find(&todos, &bolthold.Query{}); err != nil {
			return err
		}
		ctx.DOM().Morph("#todos", "todos", glv.M{"todos": todos})

	default:
		log.Printf("warning:handler not found for event => \n %+v\n", ctx.Event())
	}
	return nil
}

func main() {
	db, err := bolthold.Open("todos.db", 0666, nil)
	if err != nil {
		panic(err)
	}
	glvc := glv.Websocket("goliveview-counter", glv.DevelopmentMode(true))
	http.Handle("/", glvc.Handler(&TodosView{
		db: db,
	}))
	log.Println("listening on http://localhost:9867")
	http.ListenAndServe(":9867", nil)
}
