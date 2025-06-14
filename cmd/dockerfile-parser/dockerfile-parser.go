package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/coffeemakingtoaster/dockerfile-parser/internal/pkg/wrapper"
)

func main() {
	startTime := time.Now()
	recursive := slices.Contains(os.Args, "-r")
	output := slices.Contains(os.Args, "-o")
	count := wrapper.ParsePath(os.Args[len(os.Args)-1], recursive, output)
	diff := time.Now().Sub(startTime)
	fmt.Printf("Parsing %d files finished in %v\n", count, diff)
}
