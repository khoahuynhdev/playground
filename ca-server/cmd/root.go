package cmd

import (
	"ca-server/cmd/cli"
	"ca-server/cmd/server"
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xpki",
	Short: "eXtensible PKI (xpki) is a CLI application to manage your pki",
	Long:  `eXtensible PKI (xpki) is a command-line application CLI application to manage your pki`,
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
	rootCmd.AddCommand(cli.SetupCommand)
}
