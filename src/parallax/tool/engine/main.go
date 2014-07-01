package main

import (
	"flag"
	"fmt"
	"parallax/core"
	"parallax/engine"
	"parallax/fct"
	"runtime"
)

const (
	BID_RANDOM_EDGES string = "RandomEdges"
	BID_FIRST_EDGES         = "FirstEdges"
)

var optEngine = flag.String("name", "RandomEdges", "Engine Name (RandomEdges, FirstEdges, ...)")
var optFactor = flag.Int("factor", 20, "Price multiplication factor (Variable cost)")
var optFile = flag.String("instance", "./data/N104.DAT", "FCTP data file name")
var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var verbose = flag.Int("verbose", 1, "Print a lot of messages, level 0, 1, 2, 3")

func main() {
	fmt.Println("Parallax Engine: Engine Tool")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	gname := *optFile

	// Loading Graph
	g, err := fct.LoadGraph(gname, *verbose)
	if err != nil {
		fmt.Println("Error loading file:", gname, err)
		return
	}
	graphs := fct.NewStaticLoader(map[string]*fct.Graph{gname: g})

	// Computing Bids
	n := NewEngine(*optEngine, graphs)
	if n == nil {
		fmt.Println("Error loading engine:", *optEngine)
		return
	}
	bids := n.ComputeBid(&core.Match{gname, 5})

	// Computing Flow
	w := core.NewGurobiFlowEngine(g)
	r, err := w.ComputeFlow(map[string]*core.BidPack{"parallax": bids})
	if err != nil {
		fmt.Println("Error computing flow:", gname, err)
		return
	}
	for _, s := range r.Streams {
		fmt.Println(s)
	}
}

func NewEngine(name string, graphs fct.GraphLoader) core.BidEngine {
	switch name {
	case BID_RANDOM_EDGES:
		return engine.NewRandomEdges(graphs, *optFactor)
	case BID_FIRST_EDGES:
		return engine.NewFirstEdges(graphs, *optFactor)
	default:
		return nil
	}
}
