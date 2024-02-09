package cmd

import (
	"fmt"

	"github.com/rollwagen/aws-services/pkg/service"
	"github.com/spf13/cobra"
)

// regionsCmd represents the regions command.
var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "Print all regions",

	Run: func(_ *cobra.Command, _ []string) {
		regions, _ := service.Regions()
		for _, region := range regions {
			fmt.Println(region)
		}
	},
}

func init() {
	listCmd.AddCommand(regionsCmd)
}
