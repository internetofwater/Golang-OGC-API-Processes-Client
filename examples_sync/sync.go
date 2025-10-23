package main

import (
	"fmt"

	"github.com/internetofwater/go-ogc-api-process-client/pkg"
)

func main() {

	client, err := pkg.NewProcessesClient("https://asu-awo-pygeoapi-864861257574.us-south1.run.app/")
	if err != nil {
		panic(err)
	}

	response, err := client.ExecuteSync("config-store", map[string]any{"name": "my_test_config"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Job URL:", response.JobUrl)
	fmt.Println("Outputs:", response.Outputs)

	status, err := client.JobStatus(response.JobUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Job Status:", status.Status)

	results, err := client.GetJobResults(response.JobUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Job Results:", results)
}
