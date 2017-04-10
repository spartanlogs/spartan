package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	grokOptions := map[string]interface{}{
		"regex": `^(?P<logdate>%{MONTHDAY}[-]%{MONTH}[-]%{YEAR} %{TIME}) client %{IP:clientip}#%{POSINT:clientport} \(%{GREEDYDATA:query}\): query: %{GREEDYDATA:target} IN %{GREEDYDATA:querytype} \(%{IP:dns}\)$`,
	}
	dateOptions := map[string]interface{}{
		"field":    "logdate",
		"patterns": "02-Jan-2006 15:04:05.999999999",
		"timezone": "America/Chicago",
	}
	mutateOptions := map[string]interface{}{
		"action": "remove_field",
		"fields": []string{"logdate", "message"},
	}

	// Inputs
	file, err := inputs.New("file", map[string]interface{}{
		"path": os.Args[1],
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Filters
	mutateFilter, err := filters.New("mutate", mutateOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dateFilter, err := filters.New("date", dateOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	grok, err := filters.New("grok", grokOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	end, _ := filters.New("end", nil)

	filter := filters.NewFilterController(grok, 10)
	grok.SetNext(dateFilter)
	dateFilter.SetNext(mutateFilter)
	mutateFilter.SetNext(end)

	// Outputs
	stdout, _ := outputs.NewStdoutOutput(nil)
	stdout.SetNext(&outputs.End{})
	output := outputs.NewOutputController(stdout, 10)

	// Communication channels
	inputChan := make(chan *event.Event)
	outputChan := make(chan *event.Event)

	// Start everything
	fmt.Println("Starting outputs")
	output.Start(outputChan)

	fmt.Println("Starting filters")
	filter.Start(inputChan, outputChan)

	fmt.Println("Starting inputs")
	file.Start(inputChan)

	// Wait for Ctrl+C
	fmt.Println("Waiting for signal")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan

	//Shutdown
	fmt.Println("Shutting down inputs")
	file.Close()

	fmt.Println("Shutting down filters")
	filter.Close()

	fmt.Println("Shutting down outputs")
	output.Close()
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
