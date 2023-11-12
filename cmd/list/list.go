package list

import (
	todogo "github.com/Gandalf-Le-Dev/TodoGo/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("list called")
		tasks, err := todogo.ListTasks()
		if err != nil {
			log.Error().Msg(err.Error())
		}
		for _, task := range tasks {
			log.Info().Msg(task.Print())
		}
	},
}
