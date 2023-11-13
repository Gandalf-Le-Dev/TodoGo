package cmd

import (
	"os"

	"github.com/Gandalf-Le-Dev/TodoGo/cmd/add"
	"github.com/Gandalf-Le-Dev/TodoGo/cmd/complete"
	"github.com/Gandalf-Le-Dev/TodoGo/cmd/delete"
	"github.com/Gandalf-Le-Dev/TodoGo/cmd/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Verbose bool
var Debug bool
var Path string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "TodoGo",
	Short: "A CLI tool to manage you todo.txt file.",
	Long:  `A CLI tool to manage you todo.txt file.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(complete.CompleteCmd)
	rootCmd.AddCommand(delete.DeleteCmd)

	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Path flag for todo.txt file
	rootCmd.PersistentFlags().StringP("path", "p", "", "Path to the todo.txt file. (default: $HOME/todo.txt)") //TODO: Add default value
	viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path"))
}
