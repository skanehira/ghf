package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghf",
	Short: "ghf is cli which manage file in GitHub repository",
}

func Execute() {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		rootCmd.Help()
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
