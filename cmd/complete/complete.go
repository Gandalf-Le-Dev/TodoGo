package complete

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// complete represents the complete command
var CompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "A brief description of your command",
	Long:  `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("complete called")
	},
}
