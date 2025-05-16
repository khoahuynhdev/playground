package cmd

import (
	"ca-server/cmd/server"
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ca-server",
	Short: "CA Server is a CLI application to manage HTTP servers",
	Long:  `CA Server is a command-line application that can start HTTP servers and execute various commands.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to run comnad: %v", err)
	}
	return nil
}

func init() {
	// add subcommands here
	rootCmd.AddCommand(server.ServerCommand)
}
