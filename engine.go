package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"sort"
	"time"
)

type Engine struct {
	ConfigPath      string
	Path            string `json:"path"`
	DefaultCategory string `json:"default"`
	OrderedTasks    []*Task
	Todos           []*Todo
	DisplayAuthors  bool
	Actives         int
	Dones           int
	Authors         []string
}

type Tasks []*Task

func (t Tasks) Len() int      { return len(t) }
func (t Tasks) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Date ordered Sort
type ByDateOrdered struct{ Tasks }

func (t ByDateOrdered) Less(i, j int) bool {
	return t.Tasks[i].CreationTime > t.Tasks[j].CreationTime
}

// Date inverted sort
type ByDateInverted struct{ Tasks }

func (t ByDateInverted) Less(i, j int) bool {
	return t.Tasks[i].CreationTime > t.Tasks[j].CreationTime
}

// NewEngine returns a new engine based on the readability of /etc/todo.conf (json)
func NewEngine() (*Engine, error) {
	// Setup a new engine
	e := &Engine{}

	err = e.loadEnv()
	if err != nil {
		return nil, err
	}

	err := e.loadConfig()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func InitEngine(path string, reload bool) error {
	e := &Engine{}

	err = e.loadEnv()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		err := os.Mkdir(path, 0744)
		if err != nil {
			return fmt.Errorf("Can't create folder %s", path)
		}

	} else if !reload {
		return fmt.Errorf("Folder already exists at %s, use -r to reload an existing todo list", path)
	}

	e.Path = path

	err := e.save()
	if err != nil {
		return fmt.Errorf("Initialization failed %s %s", e.ConfigPath, err)
	}

	return nil
}

func (e *Engine) loadEnv() error {
	home := os.Getenv("HOME")
	if home == "" {
		return fmt.Errorf("Cant read home variable")
	}
	e.ConfigPath = home + "/.todo.conf"

	return nil
}

func (e *Engine) loadConfig() error {
	// Load env variable
	data, err := OpenAndReadFile(e.ConfigPath)
	if err != nil {
		return fmt.Errorf("No Todo path has been setup yet at %s", e.ConfigPath)
	}

	err = json.Unmarshal(data, e)
	if err != nil {
		return fmt.Errorf("Your config file at /etc/todo.conf is not readable")
	}

	return nil
}

// save is saving datas as /etc/todo.conf
func (e *Engine) save() error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(e.ConfigPath, data, 0744)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) SetDefault(category string) error {
	_, err := NewTodo(category, e.Path, false)
	if err != nil {
		return fmt.Errorf("Can't set Default. %s", err)
	}

	e.DefaultCategory = category

	err = e.save()
	if err != nil {
		return fmt.Errorf("Initialization failed %s", err)
	}

	return nil
}

func (e *Engine) Create(category string) error {
	_, err := NewTodo(category, e.Path, true)
	if err != nil {
		return fmt.Errorf("Can't Create category %s. %s", category, err)
	}
	return nil
}

func (e *Engine) Add(message, category string) error {

	c := category
	if c == "" {
		c = e.DefaultCategory
	}

	_, err := NewTodo(category, e.Path, true)
	if err != nil {
		return fmt.Errorf("Can't Create category %s. %s", category, err)
	}

	u, err := user.Current()
	if err != nil {
		u.Username = "nil"
	}
	t, err := NewTask(message, u.Username, c)
	if err != nil {
		return err
	}

	err = t.Save(e.Path + "/" + c)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) Delete(id string) error {

	//e.IndexTask(t, 3)
	return nil
}

func (e *Engine) findTodos(category string) error {
	files, err := ioutil.ReadDir(e.Path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if category != "" {
				if category != f.Name() {
					continue
				}
			}
			todo, err := NewTodo(f.Name(), e.Path, false)
			if err != nil {
				return err
			}
			e.Todos = append(e.Todos, todo)
		}
	}

	return nil
}

func (e *Engine) findTasks(all, onlyDones bool) error {
	for _, t := range e.Todos {
		actives, dones, err := t.LoadTasks(all, onlyDones)
		if err != nil {
			return err
		}

		e.Actives += actives
		e.Dones += dones

		e.OrderedTasks = append(e.OrderedTasks, t.Tasks...)
	}

	return nil
}

func (e *Engine) findTask(id string) (*Task, error) {
	for _, t := range e.Todos {
		t, _ := t.FindById(id)
		if t != nil {
			return t, nil
		}
	}
	return nil, fmt.Errorf("Task %s not found", id)
}

func (e *Engine) Done(id string) error {
	err := e.findTodos("")
	if err != nil {
		return err
	}

	task, err := e.findTask(id)
	if err != nil {
		return err
	}

	err = task.IsDone()
	if err != nil {
		return err
	}

	err = task.Save(e.Path + "/" + task.Category)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) List(ordered, all, onlyDones bool, category string) error {

	var dones []*Task

	err := e.findTodos(category)
	if err != nil {
		return err
	}

	err = e.findTasks(all, onlyDones)
	if err != nil {
		return err
	}

	if len(e.OrderedTasks) <= 0 {
		log.Printf("No tasks found")
		return nil
	}

	sort.Sort(ByDateOrdered{e.OrderedTasks})

	lastdate := ""

	if e.Actives > 0 {

		fmt.Printf("#TODOS(%d)\n", e.Actives)

		for _, t := range e.OrderedTasks {

			if ordered {
				tm := time.Unix(t.CreationTime, 0)
				d := tm.Format("Jan2 2006")

				if d != lastdate {
					fmt.Printf("- %s\n", d)
				}
				lastdate = d
			}

			if t.CompletionTime != 0 {
				dones = append(dones, t)
			} else {
				fmt.Printf("    | %s - [%s] %s\n", t.ID, t.Category, t.Message)
			}

		}
	}

	if len(dones) > 0 {

		lastdate = ""
		fmt.Printf("#Done(%d)\n", e.Dones)

		for _, t := range dones {

			tm := time.Unix(t.CompletionTime, 0)
			d := tm.Format("Jan2 2006")

			if d != lastdate {
				fmt.Printf("- %s\n", d)
			}
			lastdate = d

			fmt.Printf("    | %s - [%s] %s\n", t.ID, t.Category, t.Message)
		}
	}
	return nil
}
