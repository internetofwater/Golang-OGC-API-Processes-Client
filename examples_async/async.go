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

	response, err := client.ExecuteAsync("config-store", map[string]any{"name": "my_test_config"})
	if err != nil {
		panic("This currently doesn't seem to work in pygeoapi for some reason TODO: fix this")
	}
	fmt.Println("Job URL:", response)
}
