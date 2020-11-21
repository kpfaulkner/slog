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

	lp := processor.NewLogProcessor([]string{*term1})

	termMap, err := lp.ReadData(*filename, processor.RoundSecond)
  if err != nil {
  	fmt.Printf("error while reading file %s\n", err.Error())
  }

  // have a map of terms, times and counts...
  for k,v := range termMap {
  	fmt.Printf("term %s\n", k)
  	for kk,vv := range v {
  		fmt.Printf("  time %s  count %d\n", kk, vv)
	  }
  }
}

