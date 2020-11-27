package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/kpfaulkner/slog/pkg/graph"
	"github.com/kpfaulkner/slog/pkg/processor"
	"strings"
)


func getRounding(round string) (processor.RoundFactor, error) {
  round = strings.ToLower(strings.TrimSpace(round))
  switch round {
  case "h" :
  	return processor.RoundHour,nil
  case "m":
  	return processor.RoundMinute,nil
  case "s":
  	return processor.RoundSecond,nil
  }

  return processor.RoundError, errors.New("Invalid rounding")
}

func main() {
	fmt.Printf("so it begins....\n")
	filename := flag.String("f", "", "filename to process")
	terms := flag.String("terms", "", "double quoted, comma separated list of terms to search. If term needs a comma... well, dont :)")
  rounding := flag.String("round","m", "rounding by hour, minute, second. Parameter should be h,m,s")

	flag.Parse()

	if filename == nil || terms == nil || rounding == nil {
		fmt.Printf("Really must have at least a filename and a term\n")
		return
	}

	sp := strings.Split(*terms,",")
	lp := processor.NewLogProcessor(sp)

	roundFactor,err  := getRounding(*rounding)
	if err != nil {
		fmt.Printf("invalid rounding\n")
		return
	}

	termMap, err := lp.ReadData(*filename,roundFactor)
  if err != nil {
  	fmt.Printf("error while reading file %s\n", err.Error())
  	return
  }

  graphData,err  := lp.GenerateGraphData(termMap)
  if err != nil {
  	fmt.Printf("Unable to generate graph data : %s\n", err.Error())
  	return
  }


  /*
  // have a map of terms, times and counts...
  for k,v := range termMap {
  	fmt.Printf("term %s\n", k)
  	for kk,vv := range v {
  		fmt.Printf("  time %s  count %d\n", kk, vv)
	  }
  }
  */
  graph.DrawChart(graphData)

}

