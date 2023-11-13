package delete

import (
	"strconv"

	todogo "github.com/Gandalf-Le-Dev/TodoGo/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// delete/deleteCmd represents the delete/delete command
var DeleteCmd = &cobra.Command{
	Use:   "rm [id]",
	Short: "Delete a task by its ID",
	Long:  `Delete a task from the todo list by specifying its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Invalid task ID")
			return
		}

		err = todogo.DeleteTask(id, viper.GetString("path")) // Assuming you have a DeleteTask function
		if err != nil {
			log.Error().Err(err).Msg("Unable to delete task")
			return
		}

		log.Info().Msg("Task deleted successfully")
	},
}
