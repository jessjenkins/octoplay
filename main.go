package main

import (
	"fmt"
	"github.com/jessjenkins/octoplay/fetch"
	"github.com/jessjenkins/octoplay/octopus"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello octopus!")

	//TODO, better config
	octopusAPIKey := os.Getenv("OCTOPUS_API_KEY")
	electricityMPAN := os.Getenv("ELECTRICITY_MPAN")
	electricitySerial := os.Getenv("ELECTRICITY_SERIAL")

	compare := fetch.Compare{
		API:               &octopus.API{OctopusAPIKey: octopusAPIKey},
		ElectricityMPAN:   electricityMPAN,
		ElectricitySerial: electricitySerial,
	}

	result, err := compare.Run()
	if err != nil {
		log.Fatalf("fatal run error occured : %v", err)
	}
	fmt.Printf("Result: %v", *result)
}
