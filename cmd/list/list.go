package list

import (
	todogo "github.com/Gandalf-Le-Dev/TodoGo/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("list called")

		id := viper.GetInt("id")

		if id != 0 {
			log.Debug().Msgf("Getting task with id %d", id)
			task, err := todogo.GetTask(id, viper.GetString("path"))
			if err != nil {
				log.Error().Msgf("Error while getting task with id %d: %s", id, err.Error())
				return
			} else if task == nil {
				log.Error().Msgf("No task with id %d", id)
				return
			}
			log.Trace().Msg(task.Format())
		} else {
			tasks, err := todogo.ListTasks(viper.GetString("path"))
			if err != nil {
				log.Error().Msgf("Error while listing tasks: %s", err.Error())
				return
			}
			for _, task := range tasks {
				log.Trace().Msg(task.Format())
			}
		}

	},
}

func init() {
	ListCmd.Flags().IntP("id", "i", 0, "List only tasks with id equal to value of flag. Usage: --id 1")
	viper.BindPFlag("id", ListCmd.Flags().Lookup("id"))
}
