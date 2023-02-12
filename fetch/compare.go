package fetch

import "github.com/jessjenkins/octoplay/octopus"

type Compare struct {
	API               *octopus.API
	ElectricityMPAN   string
	ElectricitySerial string
}

type Result struct {
}

func (c *Compare) Run() (*Result, error) {
	gsp, err := c.API.GetGSP(c.ElectricityMPAN)
	if err != nil {
		return nil, err
	}

	_ = gsp

	return &Result{}, nil
}
