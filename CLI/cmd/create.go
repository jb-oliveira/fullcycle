/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jb-oliveira/fullcycle/CLI/internal/database"
	"github.com/spf13/cobra"
)

func newCreateCmd(categoryDB *database.Category) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "A brief description of your command",
		Long:  `A longer description that spans multiple lines and likely contains examples`,
		RunE:  runCreate(categoryDB),
	}
}

func runCreate(categoryDB *database.Category) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		err := categoryDB.Create(name, description)
		if err != nil {
			return err
		}
		fmt.Printf("Category created: %s\n", name)
		return nil
	}
}

func init() {
	createCmd := newCreateCmd(database.NewCategory(GetDB()))
	categoryCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("name", "n", "", "Category name")
	createCmd.Flags().StringP("description", "d", "", "Category description")
	// createCmd.MarkFlagRequired("name")
	// createCmd.MarkFlagRequired("description")
	// isso é mesma coisa dos dois anteriores
	createCmd.MarkFlagsRequiredTogether("name", "description")
}
