package parser

import (
	"io"
	"os"

	"github.com/lfkeitel/spartan/utils"
)

// A ParsedFile is the representation of a processed filters
// file.
type ParsedFile struct {
	// Inputs are definitions of the inputs defined in a file.
	// Input module names and options are not checked for correctness.
	Inputs []*InputDef

	// Filters are definitions of filters in the order written in
	// the file. Index 0 is the root filter in a pipeline.
	Filters []*PipelineDef

	// Outputs are definitions of outputs in the order written in
	// the file. Index 0 is the root filter in a pipeline.
	Outputs []*PipelineDef
}

// InputDef defines the module name and options map of an input.
type InputDef struct {
	Module  string
	Options *utils.InterfaceMap
}

// PipelineDef defines a pipeline object (Filter/Output). It contains
// the module name, options map, and connections to the rest of the pipeline.
type PipelineDef struct {
	Module  string
	Options *utils.InterfaceMap

	// Connections is a slice of index numbers corresponding to the index of
	// a pipeline definition in the parent ParsedFile struct.
	//
	// Length 0 -> end of pipeline
	// Length 1 -> normal next connection
	// Length 3 -> if statement connections
	//	Connections[0] -> true (inside if body)
	//	Connections[1] -> else (inside else body if present), nil if no else
	//	Connections[2] -> next object (after closing brace of if/else body)
	Connections []int
}

func ParseFile(path string) (*ParsedFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Parse(file)
}

func Parse(r io.Reader) (*ParsedFile, error) {
	return nil, nil
}
