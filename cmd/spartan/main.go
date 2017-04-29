package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lfkeitel/spartan/config/parser"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/filters"
	"github.com/lfkeitel/spartan/inputs"
	"github.com/lfkeitel/spartan/outputs"
)

var (
	configFile  string
	filtersPath string
	verFlag     bool
	testConfig  bool

	version   = ""
	buildTime = ""
	builder   = ""
	goversion = ""
)

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration file path")
	flag.StringVar(&filtersPath, "f", "", "Filter path, can be a file or directory")
	flag.BoolVar(&verFlag, "v", false, "Display version information")
	flag.BoolVar(&testConfig, "t", false, "Test main configuration")
}

func main() {
	flag.Parse()

	if verFlag {
		displayVersionInfo()
		return
	}

	if testConfig {
		testMainConfig()
		return
	}

	parsed, err := parser.ParseFile(filtersPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Filter path %s doesn't exist", filtersPath)
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	inputMods, err := inputs.CreateFromDefs(parsed.Inputs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filterPipeline, err := filters.GeneratePipeline(parsed.Filters)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	outputPipeline, err := outputs.GeneratePipeline(parsed.Outputs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filterCont := filters.NewFilterController(filterPipeline, 10)
	outputCont := outputs.NewOutputController(outputPipeline, 10)

	// Communication channels
	inputChan := make(chan *event.Event)
	outputChan := make(chan *event.Event)

	// Start everything
	fmt.Println("Starting outputs")
	outputCont.Start(outputChan)

	fmt.Println("Starting filters")
	filterCont.Start(inputChan, outputChan)

	fmt.Println("Starting inputs")
	for _, input := range inputMods {
		input.Start(inputChan)
	}

	// Wait for Ctrl+C
	fmt.Println("Waiting for signal")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan

	//Shutdown
	fmt.Println("Shutting down inputs")
	for _, input := range inputMods {
		input.Close()
	}

	fmt.Println("Shutting down filters")
	filterCont.Close()

	fmt.Println("Shutting down outputs")
	outputCont.Close()
}

func displayVersionInfo() {
	fmt.Printf(`Spartan - (C) 2017 Lee Keitel <lee@onesimussystems.com>
Version:     %s
Built:       %s
Compiled by: %s
Go version:  %s
`, version, buildTime, builder, goversion)
}

func testMainConfig() {
	// _, err := utils.NewConfig(configFile)
	// if err != nil {
	// 	fmt.Printf("Error loading configuration: %v\n", err)
	// 	os.Exit(1)
	// }
	fmt.Println("Configuration looks good")
}
