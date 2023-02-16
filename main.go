package main

import (
	"fmt"
	"github.com/jessjenkins/octoplay/fetch"
	"github.com/jessjenkins/octoplay/octopus"
	"log"
	"os"
)

const api_url = "https://api.octopus.energy/"

func main() {
	fmt.Println("Hello octopus!")

	//TODO, better config
	octopusAPIKey := os.Getenv("OCTOPUS_API_KEY")
	electricityMPAN := os.Getenv("ELECTRICITY_MPAN")
	electricitySerial := os.Getenv("ELECTRICITY_SERIAL")

	api, err := octopus.New(api_url, octopusAPIKey)
	if err != nil {
		log.Fatalf("fatal run error occured : %v", err)
	}

	compare := fetch.Compare{
		API:               api,
		ElectricityMPAN:   electricityMPAN,
		ElectricitySerial: electricitySerial,
	}

	result, err := compare.Run()
	if err != nil {
		log.Fatalf("fatal run error occured : %v", err)
	}
	fmt.Printf("Result: %+v\n", *result)
}
