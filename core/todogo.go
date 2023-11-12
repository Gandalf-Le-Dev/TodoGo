package todogo

type Task struct {
	Completed      bool
	Priority       string
	CompletionDate string
	CreationDate   string
	Description    string
	Project        []string
	Context        []string
	DueDate        string
}

func (t Task) Print() string {
	var output string

	if t.Completed {
		output += "x "
	}

	if t.Priority != "" {
		output += t.Priority + " "
	}

	if t.CompletionDate != "" {
		output += t.CompletionDate + " "
	}

	if t.CreationDate != "" {
		output += t.CreationDate + " "
	}

	output += t.Description

	for _, project := range t.Project {
		output += " " + project
	}

	for _, context := range t.Context {
		output += " " + context
	}

	if t.DueDate != "" {
		output += " due:" + t.DueDate
	}

	return output
}
