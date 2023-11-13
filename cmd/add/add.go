package add

import (
	"strings"

	todogo "github.com/Gandalf-Le-Dev/TodoGo/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var priority string
var dueDate string

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add [description]",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := strings.Join(args, " ")
		task := todogo.Task{
			Description: todogo.ExtractDescription(input),
			Priority:    priority,
			DueDate:     dueDate,
			Project:     todogo.ExtractProjects(input),
			Context:     todogo.ExtractContexts(input),
		}

		// Assuming todo.txt file is in the current directory
		err := todogo.AddTask(task, viper.GetString("path"))
		if err != nil {
			log.Error().Err(err).Msg("Error adding task")
			return
		}
		log.Info().Msg("Task added successfully.")
	},
}

func init() {
	AddCmd.Flags().StringVarP(&priority, "priority", "", "", "Set the priority of the task")
	AddCmd.Flags().StringVar(&dueDate, "due", "", "Set the due date of the task")
}
