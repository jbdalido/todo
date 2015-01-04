package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

type Engine struct {
	ConfigPath      string
	Path            string                  `json:"path"`
	DefaultCategory string                  `json:"default"`
	Index           map[string]IndexedTasks `json:"index"`
	OrderedTasks    []*Task
	Todos           []*Todo
	DisplayAuthors  bool
	Authors         []string
}

type IndexedTasks struct {
	ID       string
	Category string
	Done     bool
	Path     string
}

// NewEngine returns a new engine based on the readability of /etc/todo.conf (json)
func NewEngine() (*Engine, error) {
	// Setup a new engine
	e := &Engine{}
	e.loadEnv()

	return e, nil
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
		return fmt.Errorf("No Todo path has been setup yet, \nor %s has been deleted\nor denied access\nRun todo init -p PATH to rebuild it.", e.ConfigPath)
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

	err = ioutil.WriteFile(e.ConfigPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) InitEngine(path string, reload bool) error {
	if _, err := os.Stat(path); err != nil {
		err := os.Mkdir(path, 0644)
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
	//e.IndexTask(t, 1)

	return nil
}

func (e *Engine) Edit(id, message string) error {

	return nil
}

func (e *Engine) Delete(id string) error {

	//e.IndexTask(t, 3)
	return nil
}

func (e *Engine) Done(id string) error {

	//e.IndexTask(t, 2)
	return nil
}

func (e *Engine) findTasks() error {
	for _, t := range e.Todos {
		err := t.LoadTasks()
		if err != nil {
			return err
		}
		e.OrderedTasks = append(e.OrderedTasks, t.Tasks...)
	}
	return nil
}

func (e *Engine) findTodos() error {
	files, err := ioutil.ReadDir(e.Path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			todo, err := NewTodo(f.Name(), e.Path, false)
			if err != nil {
				return err
			}
			e.Todos = append(e.Todos, todo)
		}
	}

	return nil
}

func (e *Engine) List(ordered bool, all bool, onlyDones bool) error {
	e.findTasks()
	for i, t := range e.OrderedTasks {
		fmt.Printf("\t| %d#[%s] %s\n", i, t.Category, t.Message)
	}
	return nil
}
