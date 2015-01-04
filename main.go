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
				cli.IntFlag{
					Name:  "all, A",
					Value: 10,
					Usage: "display all tasks including (n) dones (A=n)",
				},
				cli.IntFlag{
					Name:  "done, d",
					Value: 10,
					Usage: "display only (d=n) dones tasks",
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
		},
		cli.Command{
			Name:  "rm",
			Usage: "rm a task by its id",
		},
		cli.Command{
			Name:  "done",
			Usage: "Set a task done by its id",
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
	Pre()
	err := e.InitEngine(c.Args().First(), false)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func List(c *cli.Context) {
	Pre()
	e.List(c.Bool("ordered"), c.Bool("all"), c.Bool("done"))
}

func Create(c *cli.Context) {
	Pre()
	err := e.Create(c.Args().First())
	if err != nil {
		log.Fatal(err)
	}
}
