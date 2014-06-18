package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var optServer = flag.String("server", "localhost:8080", "Game server")
var optData = flag.String("data", "./data", "Directory with FCTP data files")
var optPreload = flag.Bool("load", true, "Load all data files (instances)")
var verbose = flag.Bool("verbose", false, "Print a lot of messages")

type Vertex struct {
	id    int
	size  float64
	edges []*Edge
}

func (v *Vertex) String() string {
	return fmt.Sprintf("(%d, [%.2f,%d])", v.id, v.size, len(v.edges))
}

type Edge struct {
	i, j         *Vertex
	vCost, fCost float64
}

func (e *Edge) String() string {
	return fmt.Sprintf("(%d)-[%.2f,%.2f]->(%d)", e.i.id, e.vCost, e.fCost, e.j.id)
}

type Graph struct {
	sources, sinks map[int]*Vertex
	edges          []*Edge
}

func (g *Graph) String() string {
	return fmt.Sprintf("Sources %d, Sinks %d, Edges %d", len(g.sources), len(g.sinks), len(g.edges))
}

func NewGraph() *Graph {
	return &Graph{
		sources: make(map[int]*Vertex),
		sinks:   make(map[int]*Vertex),
		edges:   make([]*Edge, 0),
	}
}

func (g *Graph) NewEdge(i, j int, v, f float64) *Edge {
	vi, found := g.sources[i]
	if !found {
		vi = &Vertex{i, .0, make([]*Edge, 0)}
		g.sources[i] = vi
	}
	vj, found := g.sinks[j]
	if !found {
		vj = &Vertex{j, .0, make([]*Edge, 0)}
		g.sinks[j] = vj
	}
	e := &Edge{vi, vj, v, f}
	vi.edges = append(vi.edges, e)
	vj.edges = append(vj.edges, e)
	g.edges = append(g.edges, e)
	return e
}

func (g *Graph) SourceSize(i int, s float64) *Vertex {
	vi, found := g.sources[i]
	if !found {
		vi = &Vertex{i, .0, make([]*Edge, 0)}
		g.sources[i] = vi
	}
	vi.size = s
	return vi
}

func (g *Graph) SinkSize(j int, s float64) *Vertex {
	vj, found := g.sinks[j]
	if !found {
		vj = &Vertex{j, .0, make([]*Edge, 0)}
		g.sinks[j] = vj
	}
	vj.size = s
	return vj
}

func LoadGraph(path string) (*Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	g := NewGraph()

	type GraphParser int
	const (
		EDGES GraphParser = iota
		SUPPLY
		DEMAND
	)

	parser := EDGES
	scan := bufio.NewScanner(file)
	scan.Scan() //skep line 1
	scan.Scan() //skep line 2
	scan.Scan() //skep line 3
	for scan.Scan() {
		line := scan.Text()
		//fmt.Println(">>", line)
		if line == "S" {
			parser = SUPPLY
			continue
		}
		if line == "D" {
			parser = DEMAND
			continue
		}
		if line == "END" {
			break
		}
		n := strings.Fields(line)
		if parser == EDGES {
			if len(n) < 4 {
				fmt.Println("Ignoring line, no edge:", line)
				continue
			}
			i, err := strconv.ParseInt(n[0], 10, 0)
			if err != nil {
				fmt.Println("Error parsing edge source:", err)
				continue
			}
			j, err := strconv.ParseInt(n[1], 10, 0)
			if err != nil {
				fmt.Println("Error parsing edge sink:", err)
				continue
			}
			v, err := strconv.ParseFloat(n[2], 64)
			if err != nil {
				fmt.Println("Error parsing edge variable cost:", err)
				continue
			}
			f, err := strconv.ParseFloat(n[3], 64)
			if err != nil {
				fmt.Println("Error parsing edge fixed cost:", err)
				continue
			}
			e := g.NewEdge(int(i), int(j), v, f)
			if *verbose {
				fmt.Println("New Edge:", e)
			}
		} else if parser == SUPPLY {
			if len(n) < 2 {
				fmt.Println("Ignoring line, no supply:", line)
				continue
			}
			i, err := strconv.ParseInt(n[0], 10, 0)
			if err != nil {
				fmt.Println("Error parsing source:", err)
				continue
			}
			s, err := strconv.ParseFloat(n[1], 64)
			if err != nil {
				fmt.Println("Error parsing supply value:", err)
				continue
			}
			v := g.SourceSize(int(i), s)
			if *verbose {
				fmt.Println("Supply:", v)
			}
		} else if parser == DEMAND {
			if len(n) < 2 {
				fmt.Println("Ignoring line, no demand:", line)
				continue
			}
			j, err := strconv.ParseInt(n[0], 10, 0)
			if err != nil {
				fmt.Println("Error parsing sink:", err)
				continue
			}
			s, err := strconv.ParseFloat(n[1], 64)
			if err != nil {
				fmt.Println("Error parsing demand value:", err)
				continue
			}
			v := g.SinkSize(int(j), s)
			if *verbose {
				fmt.Println("Demand:", v)
			}
		} else {
			fmt.Println("Ignoring line, unknown type:", line)
		}
	}
	return g, nil
}

var INSTANCES map[string]*Graph = make(map[string]*Graph)

func Instance(name string) *Graph {
	if g, ok := INSTANCES[name]; ok {
		return g
	}
	fmt.Print("Loading ", name, "... ")
	g, err := LoadGraph(*optData + "/" + name + ".DAT")
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	INSTANCES[name] = g
	fmt.Println(g)
	return g
}

func LoadAllInstances() {
	folder, err := os.Open(*optData)
	if err != nil {
		fmt.Println("Error opening folder:", *optData, err)
		return
	}
	defer folder.Close()

	files, err := folder.Readdirnames(0)
	if err != nil {
		fmt.Println("Error listing data files from:", *optData, err)
		return
	}

	for _, file := range files {
		if !strings.HasSuffix(file, ".DAT") {
			fmt.Println("Ignoring", file)
			continue
		}
		fmt.Print("Loading ", file, "... ")
		g, err := LoadGraph(*optData + "/" + file)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		name := strings.TrimSuffix(file, ".DAT")
		INSTANCES[name] = g
		fmt.Println(g)
	}

	fmt.Println("Total:", len(INSTANCES))
}

type Bid struct {
	source, sink int
	price        float64
}

func NewBid(edge *Edge, price float64) *Bid {
	return &Bid{
		source: edge.i.id,
		sink:   edge.j.id,
		price:  price,
	}
}

type BidPack struct {
	bids []*Bid
}

func (t *BidPack) String() string {
	out := "bid\n"
	first := true
	for _, n := range t.bids {
		if first {
			first = false
		} else {
			out += "\n"
		}
		out += fmt.Sprint(n.source, n.sink, n.price)
	}
	return out
}

func NewBidPack(capacity int) *BidPack {
	return &BidPack{make([]*Bid, 0, capacity)}
}

func (p *BidPack) add(bid *Bid) {
	p.bids = append(p.bids, bid)
}

func Bidding(instanceName string, edges int) *BidPack {
	g := Instance(instanceName)
	if g == nil {
		fmt.Println("Instance not found:", instanceName)
		return nil
	}
	pack := NewBidPack(edges)
	for i := 0; i < edges; i++ {
		e := g.edges[i]
		pack.add(NewBid(e, e.vCost))
	}
	return pack
}

func Connect() {
	conn, err := net.Dial("tcp", *optServer)
	if err != nil {
		fmt.Println("Error connecting", *optServer)
		return
	}
	defer conn.Close()
	master := bufio.NewReader(conn)
	for {
		m, err := master.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		m = strings.TrimSpace(m)
		fmt.Println("Master> ", m)
		if m == "name" {
			fmt.Println("Parallax> name Parallax")
			fmt.Fprint(conn, "name Parallax")
		} else if strings.HasPrefix(m, "instance") {
			n := strings.Fields(m)
			instanceName := n[1]
			edges, err := strconv.ParseInt(n[2], 10, 0)
			if err != nil {
				fmt.Println("Error parsing number of edges:", err)
				continue
			}
			if *verbose {
				fmt.Println("Instance:", instanceName)
				fmt.Println("Number of Edges:", edges)
			}
			result := Bidding(instanceName, int(edges))
			fmt.Println("Parallax>")
			fmt.Println(result)
			fmt.Fprint(conn, result)
		} else {
			fmt.Println("unknown")
		}
	}
}

func main() {
	fmt.Println("Game Theory Player: Parallax")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	if *optPreload {
		LoadAllInstances()
	}

	Connect()
}
