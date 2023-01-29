package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/samber/lo"
)

type regionServices struct {
	Metadata struct {
		Copyright     string `json:"copyright"`
		Disclaimer    string `json:"disclaimer"`
		FormatVersion string `json:"format:version"`
		SourceVersion string `json:"source:version"`
	} `json:"metadata"`
	Prices []struct {
		Attributes struct {
			AwsRegion      string `json:"aws:region"`
			AwsServiceName string `json:"aws:serviceName"`
			AwsServiceURL  string `json:"aws:serviceUrl"`
		} `json:"attributes"`
		ID string `json:"id"`
	} `json:"prices"`
}

func Regions() ([]string, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := ssm.NewFromConfig(cfg)

	input := &ssm.GetParametersByPathInput{
		Path: aws.String("/aws/service/global-infrastructure/regions/"),
	}

	var regions []string

	paginator := ssm.NewGetParametersByPathPaginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, p := range output.Parameters {
			name := *p.Name
			regions = append(regions, name[strings.LastIndex(name, "/")+1:])
		}
	}

	sort.Strings(regions)

	return regions, nil
}

func ServiceAvailabilityPerRegion(service string, regionProgress chan<- string) (map[string]bool, error) {
	// a map of a string to boolean
	serviceAvailability := make(map[string]bool)

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := ssm.NewFromConfig(cfg)

	regions, _ := Regions()
	for _, region := range regions {

		regionProgress <- region // update channel which region is being processed

		serviceFoundForRegion := false

		path := fmt.Sprintf("/aws/service/global-infrastructure/regions/%s/services/", region)

		input := &ssm.GetParametersByPathInput{
			Path: aws.String(path),
		}

		paginator := ssm.NewGetParametersByPathPaginator(client, input)
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}
			for _, p := range output.Parameters {
				name := *p.Name
				s := name[strings.LastIndex(name, "/")+1:]

				if service == s {
					serviceAvailability[region] = true
					serviceFoundForRegion = true
					break
				}
			}
			if serviceFoundForRegion {
				break
			}
		}
		if serviceFoundForRegion {
			continue
		}
		serviceAvailability[region] = false
	}
	return serviceAvailability, nil
}

func Services() ([]string, error) {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := ssm.NewFromConfig(cfg)

	input := &ssm.GetParametersByPathInput{
		Path: aws.String("/aws/service/global-infrastructure/regions/us-east-1/services/"),
	}

	var services []string

	paginator := ssm.NewGetParametersByPathPaginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, p := range output.Parameters {
			name := *p.Name
			services = append(services, name[strings.LastIndex(name, "/")+1:])
		}
	}

	sort.Strings(services)

	return services, nil
}

func ServiceNames() ([]string, error) {
	resp, err := http.Get("https://api.regional-table.region-services.aws.a2z.com/index.json")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var rs regionServices
	err = json.NewDecoder(resp.Body).Decode(&rs)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var services []string
	for _, p := range rs.Prices {
		serviceName := p.Attributes.AwsServiceName
		present := lo.Contains[string](services, serviceName)
		if !present {
			services = append(services, serviceName)
		}
	}

	sort.Strings(services)
	return services, nil
}
