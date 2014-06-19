package main

import (
	"flag"
	"fmt"
	"parallax/engine"
	"parallax/graph"
	"runtime"
)

var optName = flag.String("name", "Parallax", "Player Name")
var optServer = flag.String("server", "localhost:8080", "Game server")
var optData = flag.String("data", "./data", "Directory with FCTP data files")
var optPreload = flag.Bool("load", true, "Load all data files (instances)")
var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var verbose = flag.Int("verbose", 1, "Print a lot of messages, level 0, 1, 2, 3")

// Dummy Player

type FirstEdges struct {
	graphs *graph.GraphLoader
}

func (n *FirstEdges) ComputeBid(m *engine.Match) *engine.BidPack {
	g := n.graphs.Instance(m.InstanceName)
	if g == nil {
		fmt.Println("Instance not found:", m.InstanceName)
		return engine.EmptyBidPack()
	}
	pack := engine.NewBidPack(m.NumberOfEdges)
	for i := 0; i < m.NumberOfEdges; i++ {
		e := g.Edges[i]
		pack.Bid(e.I.Id, e.J.Id, e.VCost)
	}
	return pack
}

func (n *FirstEdges) Update(f *engine.Flow) {
}

func main() {
	fmt.Println("Parallax Engine: Game Theory Player")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	graphs := graph.NewLoader(*optData, *verbose)
	if *optPreload {
		graphs.LoadAll()
	}

	n := &FirstEdges{graphs}
	h := engine.NewHandler(*optName, n, *verbose)
	h.Connect(*optServer)
}
