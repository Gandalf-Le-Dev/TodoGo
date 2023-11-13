package todogo

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	Id             int
	Completed      bool
	Priority       string
	CompletionDate string
	CreationDate   string
	Description    string
	Project        []string
	Context        []string
	DueDate        string
}

func (task Task) Format() string {
	var builder strings.Builder

	// Include ID
	builder.WriteString(fmt.Sprintf("#%d ", task.Id))

	// Priority
	if task.Priority != "" {
		builder.WriteString("(" + task.Priority + ") ")
	}

	// Creation Date
	if task.CreationDate == "" {
		task.CreationDate = time.Now().Format("2006-01-02")
	}
	builder.WriteString(task.CreationDate + " ")

	// Description
	builder.WriteString(task.Description + " ")

	// Projects
	for _, project := range task.Project {
		builder.WriteString("+" + project + " ")
	}

	// Contexts
	for _, context := range task.Context {
		builder.WriteString("@" + context + " ")
	}

	// Due Date
	if task.DueDate != "" {
		builder.WriteString("due:" + task.DueDate + " ")
	}

	return strings.TrimSpace(builder.String())
}
