package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "redash-visualizer",
		Short: "redash-visualizer",
		Long:  `redash-visualizer`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}
