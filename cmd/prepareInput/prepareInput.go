package main

import (
	"github.com/projectdiscovery/gologger"
	"github.com/secinto/prepareInput/prepare"
)

func main() {
	// Parse the command line flags and read config files
	options := prepare.ParseOptions()

	newPreparer, err := prepare.NewPreparer(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create preparer: %s\n", err)
	}

	err = newPreparer.Prepare()
	if err != nil {
		gologger.Fatal().Msgf("Could not prepare: %s\n", err)
	}
}
