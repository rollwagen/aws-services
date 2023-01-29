package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"

	"github.com/rollwagen/aws-services/pkg/service"
	"github.com/spf13/cobra"
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Print all possible services",

	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[24], 70*time.Millisecond)
		_ = s.Color("yellow", "bold")
		s.Suffix = " Retrieving list of services..."
		s.Start()

		services, err := service.Services()

		s.Stop()

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
