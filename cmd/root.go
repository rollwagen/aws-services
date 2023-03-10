package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/briandowns/spinner"

	"github.com/fatih/color"
	"github.com/rollwagen/aws-services/pkg/prompter"
	"github.com/rollwagen/aws-services/pkg/service"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-services",
	Short: "Print all regions with information if selected service is available",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[24], 70*time.Millisecond)
		_ = s.Color("yellow", "bold")
		s.Suffix = " Retrieving list of services..."
		s.Start()
		services, err := service.Services()
		if err != nil {
			s.Stop()
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		s.Stop()

		prompt := prompter.New()
		idx, err := prompt.Select("Select service to query", "", services)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Aborted. Exiting.")
			os.Exit(1)
		}

		s.Start()

		pollProgress := func(ch <-chan string) {
			for regionName := range ch {
				fgYellow := color.New(color.FgYellow).SprintFunc()
				s.Suffix = fmt.Sprintf(" Retrieving services for region %s ...", fgYellow(regionName))
			}
		}
		regionProgressChannel := make(chan string)
		go pollProgress(regionProgressChannel)
		serviceRegions, _ := service.AvailabilityPerRegion(services[idx], regionProgressChannel)

		s.Stop()

		var regions []string
		for k := range serviceRegions {
			regions = append(regions, k)
		}
		sort.Strings(regions)

		availabilitySign := make(map[bool]string)
		availabilitySign[true] = "✔"
		availabilitySign[false] = "✖"

		for _, region := range regions {
			// for region, available := range serviceRegions {
			c := color.New(color.FgGreen)
			available := serviceRegions[region]
			if !available {
				c = color.New(color.FgRed)
			}
			_, _ = c.Printf("%s %s\n", availabilitySign[available], region)
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
