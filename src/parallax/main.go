package main

import (
	"flag"
	"fmt"
	"parallax/core"
	"parallax/engine"
	"parallax/fct"
	"runtime"
)

var optName = flag.String("name", "Parallax", "Player Name")
var optServer = flag.String("server", "localhost:8080", "Game server")
var optData = flag.String("data", "./data", "Directory with FCTP data files")
var optPreload = flag.Bool("load", true, "Load all data files (instances)")
var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var verbose = flag.Int("verbose", 1, "Print a lot of messages, level 0, 1, 2, 3")

func main() {
	fmt.Println("Parallax Engine: Game Theory Player")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	graphs := fct.NewFileLoader(*optData, *verbose)
	if *optPreload {
		graphs.LoadAll()
	}

	n := engine.NewFirstEdges(graphs, 1)
	h := core.NewHandler(*optName, n, *verbose)
	h.Connect(*optServer)
}
