package cmd

import (
	"fmt"
	"github.com/rollwagen/qrs/pkg/prompter"
	"github.com/rollwagen/qrs/pkg/ssm"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qrs",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example: ...`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[28], 70*time.Millisecond)
		_ = s.Color("yellow", "bold")
		s.Suffix = " Retrieving list of services..."
		s.Start()
		services, err := ssm.Services()
		if err != nil {
			s.Stop()
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		s.Stop()
		p := prompter.New()
		idx, err := p.Select("Select service to query", "", services)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Aborted. Exiting.")
			os.Exit(1)
		}
		s.UpdateCharSet(spinner.CharSets[24])
		s.Start()
		servicePerRegion, _ := ssm.ServiceAvailabilityPerRegion(services[idx])
		s.Stop()
		// iterate over map
		for region, available := range servicePerRegion {
			fmt.Printf("%s %t\n", region, available)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
