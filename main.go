package main

import (
	"flag"
	"fmt"
	"github.com/kpfaulkner/slog/pkg/processor"

)


func main() {
	fmt.Printf("so it begins....\n")
	filename := flag.String("f", "", "filename to process")
	term1 := flag.String("term1", "", "first term to search for")

	flag.Parse()

	if filename == nil || term1 == nil {
		fmt.Printf("Really must have at least a filename and a term\n")
		return
	}

	lp := processor.NewLogProcessor()
	err := lp.ReadData(*filename, processor.RoundMinute)
  if err != nil {
  	fmt.Printf("error while reading file %s\n", err.Error())
  }

}

