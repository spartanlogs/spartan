package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spartanlogs/spartan/config/parser"
	"github.com/spartanlogs/spartan/event"
	"github.com/spartanlogs/spartan/filters"
	"github.com/spartanlogs/spartan/inputs"
	"github.com/spartanlogs/spartan/outputs"
)

var (
	configFile       string
	filtersPath      string
	verFlag          bool
	testConfig       bool
	testFilterConfig bool

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
	flag.BoolVar(&testFilterConfig, "configtest", false, "Test filter configuration")
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

	if testFilterConfig {
		return
	}

	inputMods, err := inputs.CreateFromDefs(parsed.Inputs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filterCont, err := filters.GeneratePipeline(parsed.Filters, 100, 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	outputPipeline, err := outputs.GeneratePipeline(parsed.Outputs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	outputCont := outputs.NewOutputController(outputPipeline, 100)

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
