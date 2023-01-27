/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/rollwagen/qrs/pkg/ssm"
	"github.com/spf13/cobra"
	"os"
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example: ...`,

	Run: func(cmd *cobra.Command, args []string) {
		services, err := ssm.Services()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, s := range services {
			fmt.Println(s)
		}
	},
}

func init() {
	listCmd.AddCommand(servicesCmd)
}
