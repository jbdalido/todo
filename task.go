package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jbdalido/todo/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"
)

type Task struct {
	ID             string `json:"id" yaml:"id"`                         // String base64 of CreationTime + rand()
	Category       string `json:"category" yaml:"category"`             // Category as plain text (match a folder)
	Author         string `json:"author" yaml:"author"`                 // Contains author as plain test
	Message        string `json:"message" yaml:"message"`               // Message as plain text
	CreationTime   int64  `json:"creationTime" yaml:"creationTime"`     // Unix Timestamp
	UpdateTime     int64  `json:"updateTime" yaml:"updateTime"`         // Unix Timestamp
	CompletionTime int64  `json:"completionTime" yaml:"completionTime"` // Unix Timestamp
	DeletionTime   int64  `json:"deletionTime" yaml:"deletionTime"`     // Unix Timestamp
}

// NewTask initialize a New task and setup a new TaskMessage
func NewTask(msg, author, category string) (*Task, error) {
	if msg == "" {
		return nil, fmt.Errorf("Task can't be null")
	}
	time := time.Now().Unix()

	return &Task{
		ID:             getID(msg, time),
		Category:       strings.ToLower(category),
		Author:         author,
		Message:        msg,
		CreationTime:   time,
		UpdateTime:     time,
		CompletionTime: 0,
		DeletionTime:   0,
	}, nil
}

func NewTaskFromJson(data []byte) (*Task, error) {
	t := &Task{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	err = t.Validate()
	if err != nil {
		return nil, err
	}

	return t, nil
}
func NewTaskFromYaml(data []byte) (*Task, error) {
	t := &Task{}
	err := yaml.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	err = t.Validate()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Task) Validate() error {
	if t.Message == "" {
		return fmt.Errorf("invalid format, no message")
	}

	if t.ID == "" {
		return fmt.Errorf("invalid format, no ID")
	}

	if t.CreationTime == 0 {
		return fmt.Errorf("invalid format no Creation Date")
	}

	return nil
}

// encodeMessage is encoding a string to md5 to obtain a gitcommit-like id
func getID(msg string, key int64) string {
	sum := md5.Sum([]byte(msg + string(key)))
	return fmt.Sprintf("%x", sum)
}

// EditMessage is used to edit the message of a task
func (t *Task) EditMessage(msg, author string) error {
	if msg == "" {
		return fmt.Errorf("Task message can't be null")
	}
	t.Message = msg
	t.UpdateTime = time.Now().Unix()

	return nil
}

func (t *Task) Save(path string) error {
	data, err := t.ToJson()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path+"/"+t.ID, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

// ToJson returns Task struct as json text
func (t *Task) ToJson() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToYaml returns Task struct as json text
func (t *Task) ToYaml() (string, error) {
	data, err := yaml.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
