package main

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID               string // String base64 of creationDate + rand()
	Category         string // Category as plain text (match a folder)
	Author           string // Contains author as plain test
	Message          string // Message as plain text
	CreationDate     int64  // Unix Timestamp
	ModificationDate int64  // Unix Timestamp
	CompletionDate   int64  // Unix Timestamp
	DeletionDate     int64  // Unix Timestamp
}

// NewTask initialize a New task and setup a new TaskMessage
func NewTask(message, author, category string) (*Task, error) {
	if msg == "" {
		return fmt.Errorf("Task can't be null")
	}
	time := time.Unix()

	return &Task{
		ID:               getID(msg, time),
		Category:         strings.ToLower(category),
		Author:           author,
		Message:          message,
		CreationDate:     time,
		ModificationDate: time,
		CompletionDate:   0,
		DeletionDate:     0,
	}
}

// EditMessage is used to edit the message of a task
func (t *Task) EditMessage(msg, author string) error {
	if msg == "" {
		return fmt.Errorf("Task can't be null")
	}
	t.Message = msg
	t.ModificationDate = time.Unix()

	return nil
}

// encodeMessage is encoding a string to base64 to obtain a gitcommit-like id
func getID(msg string, key int64) string {
	return base64.StdEncoding.EncodeToString([]byte(msg + string(key)))
}
