/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// categoryCmd represents the category command
var categoryCmd = &cobra.Command{
	Use:   "category",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	// PreRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("pre run env: %s\n", enviroment)
	// },
	// PostRun: func(cmd *cobra.Command, args []string) {
	// 	fmt.Printf("post run env: %s\n", enviroment)
	// },
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	fmt.Printf("run error env: %s\n", enviroment)
	// 	return nil
	// },
}

func init() {
	rootCmd.AddCommand(categoryCmd)
}
