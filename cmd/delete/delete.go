/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package delete

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// delete/deleteCmd represents the delete/delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long:  `A longer description.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("delete called")
	},
}
