package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
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
	ctx := context.Background()

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

func AvailabilityPerRegion(service string, regionProgress chan<- string) (map[string]bool, error) {
	// a map of a string (region) to boolean (available true/false)
	serviceAvailability := sync.Map{}

	const concurrency = 3
	concurrencyPool := pool.New().WithMaxGoroutines(concurrency)

	regions, _ := Regions()
	for _, r := range regions {
		region := r

		regionProgress <- region // update channel which region is being processed

		concurrencyPool.Go(func() {
			isAvailable(region, service, &serviceAvailability)
		})
	}

	concurrencyPool.Wait()

	m := make(map[string]bool) // map string (region) -> boolean (available true/false)
	serviceAvailability.Range(func(k any, v any) bool { m[k.(string)] = v.(bool); return true })

	return m, nil
}

// isAvailable check if service is available in region, and updates the given availability map
func isAvailable(region string, service string, availability *sync.Map) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	client := ssm.NewFromConfig(cfg)

	path := fmt.Sprintf("/aws/service/global-infrastructure/regions/%s/services/", region)

	input := &ssm.GetParametersByPathInput{
		Path: aws.String(path),
	}

	paginator := ssm.NewGetParametersByPathPaginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			panic(err)
		}
		for _, p := range output.Parameters {
			name := *p.Name
			s := name[strings.LastIndex(name, "/")+1:]

			if service == s {
				availability.Store(region, true)
				return
			}
		}
	}
	availability.Store(region, false)
}

func Services() ([]string, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	ssmClient := ssm.NewFromConfig(cfg)

	input := &ssm.GetParametersByPathInput{
		// us-east-1 picked as reference region as tend to have most services, and
		// have new services first
		Path:       aws.String("/aws/service/global-infrastructure/regions/us-east-1/services/"),
		MaxResults: aws.Int32(10), // ten is max
	}

	var services []string

	paginator := ssm.NewGetParametersByPathPaginator(ssmClient, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, p := range output.Parameters {
			name := *p.Name
			serviceName := name[strings.LastIndex(name, "/")+1:]
			services = append(services, serviceName)
		}
	}

	sort.Strings(services)

	return services, nil
}

func Names() ([]string, error) {
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
		return nil, err
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
