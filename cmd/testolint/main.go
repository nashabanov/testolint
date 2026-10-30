package main

import (
	"fmt"
	"os"

	"github.com/nashabanov/testolint/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	a, err := analyzer.New(analyzer.Config{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	singlechecker.Main(a)
}
