package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Todo struct {
	Name     string
	Tasks    []*Task
	WorkPath string
	Path     string
}

// NewTodo returns a Todo and can create it
func NewTodo(name, wp string, create bool) (*Todo, error) {
	p := path.Clean(wp) + "/" + name
	if _, err := os.Stat(p); err != nil {
		if !create {
			return nil, fmt.Errorf("Category %s does not exist", name)
		}
		err := os.Mkdir(p, 0744)
		if err != nil {
			return nil, fmt.Errorf("Can't create folder %s", p)
		}
		log.Printf("Category %s created", name)
	}

	return &Todo{
		Name:     name,
		WorkPath: wp,
		Path:     p,
	}, nil
}

func (t *Todo) FindById(id string) (*Task, error) {
	files, err := ioutil.ReadDir(t.Path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {

		seek := f.Name()
		if len(id) != len(seek) {
			seek = f.Name()[0:len(id)]
		}

		if seek == id {
			data, err := OpenAndReadFile(t.Path + "/" + f.Name())
			if err != nil {
				return nil, err
			}

			t, err := NewTaskFromJson(data)
			if err != nil {
				return nil, err
			}

			return t, nil
		}
	}
	return nil, fmt.Errorf("Task not found")
}

// LoadTasks read the filesystem to find the tasks
func (t *Todo) LoadTasks(all, onlyDones bool) (int, int, error) {

	var dones, actives int

	files, err := ioutil.ReadDir(t.Path)
	if err != nil {
		return 0, 0, err
	}

	for _, f := range files {
		data, err := OpenAndReadFile(t.Path + "/" + f.Name())
		if err != nil {
			return 0, 0, err
		}
		// TODO : Here deal with a NewTaskFromByte
		// to load either yaml or json or whatever
		task, err := NewTaskFromJson(data)
		if err != nil {
			return 0, 0, err
		}

		if !all && !onlyDones {
			if task.CompletionTime != 0 {
				continue
			}
		}

		if onlyDones {
			if task.CompletionTime == 0 {
				continue
			}
		}

		if task.CompletionTime != 0 {
			dones += 1
		} else {
			actives += 1
		}

		t.Tasks = append(t.Tasks, task)
	}

	return actives, dones, nil
}
