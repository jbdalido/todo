package main

import (
	"fmt"
	"io/ioutil"
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
		err := os.Mkdir(p, 0644)
		if err != nil {
			return nil, fmt.Errorf("Can't create folder %s", p)
		}
	}

	return &Todo{
		Name:     name,
		WorkPath: wp,
		Path:     p,
	}, nil
}

// LoadTasks read the filesystem to find the tasks
func (t *Todo) LoadTasks() error {

	files, err := ioutil.ReadDir(t.Path)
	if err != nil {
		return err
	}

	for _, f := range files {
		data, err := OpenAndReadFile(t.Path + "/" + f.Name())
		if err != nil {
			return err
		}
		// TODO : Here deal with a NewTaskFromByte
		// to load either yaml or json or whatever
		task, err := NewTaskFromJson(data)
		if err != nil {
			return err
		}
		t.Tasks = append(t.Tasks, task)
	}

	return nil
}
