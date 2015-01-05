package main

import (
	"github.com/jbdalido/todo/Godeps/_workspace/src/github.com/codegangsta/cli"
	"log"
	"os"
)

var (
	e   *Engine
	err error
)

func main() {
	// Setup config file
	//config := engine.NewConfig()
	app := cli.App{
		Name:    "Todo",
		Usage:   "TODO application, for fun and no profit",
		Author:  "Jean-Baptiste Dalido",
		Version: "1.0.0",
	}

	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "list",
			Usage: "list tasks in your todos",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ordered, o",
					Usage: "order tasks by dates",
				},
				cli.BoolFlag{
					Name:  "all, A",
					Usage: "display all tasks including (n) dones (A=n)",
				},
				cli.BoolFlag{
					Name:  "done, d",
					Usage: "display only (d=n) dones tasks",
				},
				cli.StringFlag{
					Name:  "category, c",
					Value: "",
					Usage: "filter by categories",
				},
			},
			Action: List,
		},
		cli.Command{
			Name:   "init",
			Usage:  "init a new todos folder",
			Action: Init,
		},
		cli.Command{
			Name:   "create",
			Usage:  "create a new todo category",
			Action: Create,
		},
		cli.Command{
			Name:  "add",
			Usage: "add a new task to the chosen or default todo",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "message, m",
					Value: "",
					Usage: "task message, what else?",
				},
				cli.StringFlag{
					Name:  "category, c",
					Value: "",
					Usage: "category to add task to",
				},
			},
			Action: Add,
		},
		cli.Command{
			Name:  "rm",
			Usage: "rm a task by its id",
		},
		cli.Command{
			Name:   "done",
			Usage:  "Set a task done by its id",
			Action: Done,
		},
	}

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}

	app.Run(os.Args)
}

func Pre() {
	e, err = NewEngine()
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func Init(c *cli.Context) {
	err := InitEngine(c.Args().First(), false)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func List(c *cli.Context) {
	Pre()
	err := e.List(c.Bool("ordered"), c.Bool("all"), c.Bool("done"), c.String("category"))
	if err != nil {
		log.Fatal(err)
	}
}

func Create(c *cli.Context) {
	Pre()
	err := e.Create(c.Args().First())
	if err != nil {
		log.Fatal(err)
	}
}

func Add(c *cli.Context) {
	Pre()
	err := e.Add(c.String("message"), c.String("category"))
	if err != nil {
		log.Fatal(err)
	}
}

func Done(c *cli.Context) {
	Pre()
	err := e.Done(c.Args().First())
	if err != nil {
		log.Fatal(err)
	}
}
