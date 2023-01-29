package cmd

import (
	"fmt"
	"os"

	"github.com/rollwagen/aws-services/pkg/service"
	"github.com/spf13/cobra"
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List all services by name",

	Run: func(cmd *cobra.Command, args []string) {
		services, err := service.Services()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, s := range services {
			_, _ = fmt.Fprintln(os.Stdout, s)
		}
	},
}

func init() {
	listCmd.AddCommand(servicesCmd)
}
