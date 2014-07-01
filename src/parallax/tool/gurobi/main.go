package main

import (
	"flag"
	"fmt"
	"parallax/core"
	"parallax/fct"
)

var optFile = flag.String("instance", "./data/N104.DAT", "FCTP data file name")
var verbose = flag.Int("verbose", 1, "Print a lot of messages, level 0, 1, 2, 3")

func main() {
	fmt.Println("Parallax Engine: Gurobi Tool")

	flag.Parse()

	g, err := fct.LoadGraph(*optFile, *verbose)
	if err != nil {
		fmt.Println("Error loading file:", *optFile, err)
		return
	}
	s := core.NewGurobiSolver()
	r, err := s.ComputeFlow(g)
	if err != nil {
		fmt.Println("Error computing flow:", *optFile, err)
		return
	}
	for _, f := range r {
		fmt.Println(f)
	}
}
