package main

import (
	"flag"
	"fmt"
	"parallax/core"
	"parallax/engine"
	"parallax/fct"
	"runtime"
)

var optEngine = flag.String("name", engine.BID_RANDOM_EDGES, "Engine Name (RandomEdges, FirstEdges, ...)")
var optFactor = flag.Float64("factor", 2., "Price multiplication factor (Variable cost)")
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
	n := engine.New(*optEngine, graphs, *optFactor)
	if n == nil {
		fmt.Println("Error loading engine:", *optEngine)
		return
	}
	k := int(20 * g.Size() / 100)
	if k < 1 {
		k = 1
	}
	m := &core.Match{gname, k}
	bids := n.ComputeBid(m)

	fmt.Println(m)
	fmt.Println(bids)

	// Computing Flow
	w := core.NewGurobiFlowEngine(g)
	r, err := w.ComputeFlow(map[string]*core.BidPack{"parallax": bids})
	if err != nil {
		fmt.Println("Error computing flow:", gname, err)
		return
	}

	if len(r.Streams) == 0 {
		fmt.Println("No bids selected!")
	}
	for _, s := range r.Streams {
		fmt.Println(s)
	}
}
