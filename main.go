package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/coffeemakingtoaster/dockerfile-parser/pkg/wrapper"
)

func main() {
	startTime := time.Now()
	recursive := slices.Contains(os.Args, "-r")
	wrapper.ParsePath(os.Args[len(os.Args)-1], recursive)
	diff := time.Now().Sub(startTime)
	fmt.Printf("Parsing finished in %v\n", diff)
}
