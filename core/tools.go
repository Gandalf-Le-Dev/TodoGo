package todogo

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
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

func ListTasks() ([]Task, error) {
	filePath, err := FindTodoFile("../", 0)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to find todo.txt file")
		return nil, err
	}
	if filePath == "" {
		log.Warn().Msg("Unable to find todo.txt file")
		return nil, nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to open %s", filePath)
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

func parseLine(line string) Task {
	var task Task
	words := strings.Fields(line)

	// Parsing for completion and priority
	if len(words) > 0 {
		if words[0] == "x" {
			task.Completed = true
			words = words[1:] // Remove the 'x'
		}

		if len(words) > 0 && isPriority(words[0]) {
			task.Priority = words[0]
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
			task.Project = append(task.Project, word)
		} else if strings.HasPrefix(word, "@") {
			task.Context = append(task.Context, word)
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
