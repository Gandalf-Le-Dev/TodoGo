package todogo

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// FindTodoFile searches for todo.txt in the current directory and subdirectories up to a specified depth.
// depth of 0 means search only in the current directory, 1 means search in subdirectories one level deep, and so on.
// Returns the path to the file if found, otherwise returns an empty string.
func FindTodoFile(startDir string, maxDepth int) (string, error) {
	var foundPath string

	err := filepath.WalkDir(startDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // propagate the error
		}

		if d.IsDir() {
			depth := getDepth(startDir, path)
			if depth > maxDepth {
				return filepath.SkipDir
			}
		} else if d.Name() == "todo.txt" {
			foundPath = path
			log.Debug().Msgf("Found todo.txt at %s", path)
			return filepath.SkipDir // Found the file, no need to search further
		}

		return nil
	})

	return foundPath, err
}

// getDepth calculates the depth of a subdirectory relative to the start directory.
func getDepth(startDir, dir string) int {
	relativePath, err := filepath.Rel(startDir, dir)
	if err != nil {
		log.Warn().Err(err).Msgf("Unable to calculate depth of %s", dir)
		return 0 // if unable to calculate, treat as same level
	}
	return len(filepath.SplitList(relativePath)) - 1
}

// ListTasks returns a slice of Task and an error. It reads the todo.txt file at the provided path and parses each line into a Task struct.
// If the file is not found, it returns a nil slice and an error. If there is an error while reading the file, it returns an error.
func ListTasks(fileDirPath ...string) ([]Task, error) {
	if len(fileDirPath) == 0 || fileDirPath[0] == "" {
		fileDirPath[0] = "./" // Default path
	}

	filePath, err := FindTodoFile(fileDirPath[0], 0)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to find todo.txt file at \"%s\". Please create one in the current or given directory or specify a path with the --path flag.", fileDirPath[0])
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to open %s.", filePath)
		return nil, err
	}
	defer file.Close()

	var tasks []Task
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		task := parseLine(line)
		tasks = append(tasks, task)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msgf("Error reading %s", filePath)
		return nil, err
	}

	return tasks, nil
}

func GetTask(id int, fileDirPath ...string) (*Task, error) {
	tasks, err := ListTasks(fileDirPath...)
	if err != nil {
		log.Error().Err(err).Msgf("Error while getting task with id %d", id)
		return nil, err
	}

	for _, task := range tasks {
		if task.Id == id {
			return &task, nil
		}
	}

	return nil, nil
}

// AddTask adds a new task to the todo.txt file.
func AddTask(task Task, fileDirPath ...string) error {
	if len(fileDirPath) == 0 || fileDirPath[0] == "" {
		fileDirPath[0] = "./" // Default path
	}

	filePath, err := FindTodoFile(fileDirPath[0], 0)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to find todo.txt file at \"%s\". Please create one in the current or given directory or specify a path with the --path flag.", fileDirPath[0])
		return err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to open %s.", filePath)
		return err
	}

	nbLines, err := getNumberOfLines(filePath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to get number of lines in %s.", filePath)
		return err
	}
	defer file.Close()

	// Format the task line
	task.Id = nbLines + 1
	taskLine := task.Format()

	// Write the task line to the file
	_, err = file.WriteString(taskLine + "\n")
	return err
}

func getNumberOfLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

func parseLine(line string) Task {
	var task Task
	words := strings.Fields(line)

	// Parsing for ID
	if len(words) > 0 && strings.HasPrefix(words[0], "#") {
		idStr := strings.TrimPrefix(words[0], "#")
		id, err := strconv.Atoi(idStr)
		if err == nil {
			task.Id = id
		}
		words = words[1:] // Remove the ID
	}

	// Parsing for completion and priority
	if len(words) > 0 {
		if words[0] == "x" {
			task.Completed = true
			words = words[1:] // Remove the 'x'
		}

		if len(words) > 0 && isPriority(words[0]) {
			priority := strings.Trim(words[0], "()")
			task.Priority = priority
			words = words[1:] // Remove the priority
		}
	}

	// Parsing for dates
	for i, word := range words {
		if isDate(word) {
			if task.Completed && task.CompletionDate == "" {
				task.CompletionDate = word
			} else if !task.Completed && task.CreationDate == "" {
				task.CreationDate = word
			}
			words = append(words[:i], words[i+1:]...)
			break // Stop after finding the first date
		}
	}

	// Parsing for projects, contexts, and due date
	for _, word := range words {
		if strings.HasPrefix(word, "+") {
			project := strings.TrimPrefix(word, "+")
			task.Project = append(task.Project, project)
		} else if strings.HasPrefix(word, "@") {
			context := strings.TrimPrefix(word, "@")
			task.Context = append(task.Context, context)
		} else if strings.HasPrefix(word, "due:") {
			task.DueDate = strings.TrimPrefix(word, "due:")
		} else {
			task.Description += word + " "
		}
	}

	task.Description = strings.TrimSpace(task.Description)

	return task
}

func isPriority(word string) bool {
	matched, _ := regexp.MatchString(`^\([A-Z]\)$`, word)
	return matched
}

func isDate(word string) bool {
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, word)
	return matched
}

// extractProjects finds +Project tags in the input string.
func ExtractProjects(input string) []string {
	return ExtractByPrefix(input, `+\S+`)
}

// extractContexts finds @Context tags in the input string.
func ExtractContexts(input string) []string {
	return ExtractByPrefix(input, `@\S+`)
}

// extractByPrefix uses a regular expression to find words with a specific prefix.
func ExtractByPrefix(input, pattern string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}

	matches := re.FindAllString(input, -1)
	for i, match := range matches {
		matches[i] = strings.TrimPrefix(match, string(match[0]))
	}
	return matches
}

// extractDescription removes project and context tags from the input string.
func ExtractDescription(input string) string {
	re := regexp.MustCompile(`(\+\S+|@\S+)`)
	return strings.TrimSpace(re.ReplaceAllString(input, ""))
}
