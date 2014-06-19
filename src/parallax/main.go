package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var optName = flag.String("name", "Parallax", "Player Name")
var optServer = flag.String("server", "localhost:8080", "Game server")
var optData = flag.String("data", "./data", "Directory with FCTP data files")
var optPreload = flag.Bool("load", true, "Load all data files (instances)")
var optThreads = flag.Int("threads", runtime.NumCPU(), "Number of system threads")
var verbose = flag.Int("verbose", 1, "Print a lot of messages, level 0, 1, 2, 3")

// FCTP data model

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

func (g *Graph) v(m map[int]*Vertex, id int) *Vertex {
	v, found := m[id]
	if !found {
		v = &Vertex{id, .0, make([]*Edge, 0)}
		m[id] = v
	}
	return v
}

func (g *Graph) NewEdge(i, j int, v, f float64) *Edge {
	vi := g.v(g.sources, i)
	vj := g.v(g.sinks, j)
	e := &Edge{vi, vj, v, f}
	vi.edges = append(vi.edges, e)
	vj.edges = append(vj.edges, e)
	g.edges = append(g.edges, e)
	return e
}

func (g *Graph) SourceSize(i int, s float64) *Vertex {
	vi := g.v(g.sources, i)
	vi.size = s
	return vi
}

func (g *Graph) SinkSize(j int, s float64) *Vertex {
	vj := g.v(g.sinks, j)
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
	scan.Scan() //skip line 1
	scan.Scan() //skip line 2
	scan.Scan() //skip line 3
	for scan.Scan() {
		line := scan.Text()
		if *verbose > 2 {
			fmt.Println(">>", line)
		}
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
			if *verbose > 1 {
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
			if *verbose > 1 {
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
			if *verbose > 1 {
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

// Game Protocol

type Match struct {
	instanceName string
	edges        int
}

func (m *Match) String() string {
	return fmt.Sprint(m.instanceName, " ", m.edges)
}

func ParseMatch(m string) (*Match, error) {
	n := strings.Fields(m)
	if len(n) != 3 {
		return nil, fmt.Errorf("Wrong number of fields (3): %d", len(n))
	}
	instanceName := n[1]
	edges, err := strconv.ParseInt(n[2], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("Error parsing number of edges: %s", err)
	}
	return &Match{instanceName, int(edges)}, nil
}

type Flow struct {
	streams []*Stream
}

type Stream struct {
	source, sink int
	amount       float64
	owner        string
	price        float64
	bids         int
}

func (f *Flow) String() string {
	return fmt.Sprint("Number of Edges ", len(f.streams))
}

func ParseStream(m string) (*Stream, error) {
	n := strings.Fields(m)
	if len(n) != 6 {
		return nil, fmt.Errorf("Wrong number of fields (6): %d", len(n))
	}
	source, err := strconv.ParseInt(n[0], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("Error parsing result source: %s", err)
	}
	sink, err := strconv.ParseInt(n[1], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("Error parsing result sink: %s", err)
	}
	owner := n[2]
	bids, err := strconv.ParseInt(n[3], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("Error parsing result number of bids: %s", err)
	}
	price, err := strconv.ParseFloat(n[4], 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing result bid: %s", err)
	}
	amount, err := strconv.ParseFloat(n[5], 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing result amount: %s", err)
	}
	return &Stream{
		int(source),
		int(sink),
		amount,
		owner,
		price,
		int(bids),
	}, nil
}

func ParseFlow(m string, buf *bufio.Reader) (*Flow, error) {
	n := strings.Fields(m)
	if len(n) != 2 {
		return nil, fmt.Errorf("Wrong number of fields (2): %d", len(n))
	}
	k, err := strconv.ParseInt(n[1], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("Error parsing number of results: %s", err)
	}
	streams := make([]*Stream, k)
	for i := 0; i < int(k); i++ {
		r, err := buf.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("Error reading stream result (%d): %s", i+1, err)
		}
		r = strings.TrimSpace(r)
		s, err := ParseStream(r)
		if err != nil {
			return nil, fmt.Errorf("Error parsing stream result (%d): %s", i+1, err)
		}
		streams[i] = s
	}
	return &Flow{streams}, nil
}

type Bid struct {
	source, sink int
	price        float64
}

func (b *Bid) String() string {
	return fmt.Sprintf("%d %d %.2f", b.source, b.sink, b.price)
}

func NewBid(source, sink int, price float64) *Bid {
	return &Bid{
		source: source,
		sink:   sink,
		price:  price,
	}
}

type BidPack struct {
	bids []*Bid
}

func (p *BidPack) String() string {
	out := "bid"
	for _, b := range p.bids {
		out += "\n" + b.String()
	}
	return out
}

func NewBidPack(capacity int) *BidPack {
	return &BidPack{make([]*Bid, 0, capacity)}
}

func EmptyBidPack() *BidPack {
	return NewBidPack(0)
}

func (p *BidPack) Bid(source, sink int, price float64) *Bid {
	bid := NewBid(source, sink, price)
	p.bids = append(p.bids, bid)
	return bid
}

type Proffit struct {
	name  string
	value float64
}

type ProffitSlice []*Proffit

func (p *Proffit) String() string {
	return fmt.Sprintf("%s %.2f", p.name, p.value)
}

func (m ProffitSlice) String() string {
	out := ""
	for i, p := range m {
		if i != 0 {
			out += "\n"
		}
		out += p.String()
	}
	return out
}

func ParseProffit(m string) (*Proffit, error) {
	n := strings.Fields(m)
	if len(n) != 2 {
		return nil, fmt.Errorf("Wrong number of fields (2): %d", len(n))
	}
	name := n[0]
	value, err := strconv.ParseFloat(n[1], 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing value of proffits: %s", err)
	}
	return &Proffit{name, value}, nil
}

func ParseProffits(m string, buf *bufio.Reader) (ProffitSlice, error) {
	n := strings.Fields(m)
	if len(n) != 2 {
		return nil, fmt.Errorf("Wrong number of fields (2): %d", len(n))
	}
	k, err := strconv.ParseInt(n[1], 10, 0)
	if err != nil {
		return nil, fmt.Errorf("Error parsing number of proffits: %s", err)
	}
	proffits := make(ProffitSlice, k)
	for i := 0; i < int(k); i++ {
		p, err := buf.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("Error reading proffit (%d): %s", i+1, err)
		}
		p = strings.TrimSpace(p)
		pi, err := ParseProffit(p)
		if err != nil {
			return nil, fmt.Errorf("Error parsing proffit (%d): %s", i+1, err)
		}
		proffits[i] = pi
	}
	return proffits, nil
}

type Engine interface {
	ComputeBid(m *Match) *BidPack
	Update(f *Flow)
}

type Handler struct {
	name   string
	engine Engine
}

func NewHandler(name string, engine Engine) *Handler {
	return &Handler{name, engine}
}

func (h *Handler) Run(conn io.ReadWriter) {
	master := bufio.NewReader(conn)
	for {
		m, err := master.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
		m = strings.TrimSpace(m)
		fmt.Println("Master>", m)
		if m == "name" {
			name := "name " + h.name
			fmt.Println("Parallax>", name)
			fmt.Fprint(conn, name)
		} else if strings.HasPrefix(m, "instance") {
			t, err := ParseMatch(m)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if *verbose > 0 {
				fmt.Println("Match:", t)
			}
			result := h.engine.ComputeBid(t)
			fmt.Println("Parallax>")
			_out := result.String()
			if *verbose < 2 && len(_out) > 50 {
				_out = _out[0:50] + "..."
			}
			fmt.Println(_out)
			fmt.Fprint(conn, result)
		} else if strings.HasPrefix(m, "result") {
			t, err := ParseFlow(m, master)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if *verbose > 0 {
				fmt.Println("Flow:", t)
			}
			h.engine.Update(t)
		} else if strings.HasPrefix(m, "end") {
			t, err := ParseProffits(m, master)
			if err != nil {
				fmt.Println(err)
				break
			}
			if *verbose > 0 {
				fmt.Println("Proffit:", t)
			}
			fmt.Println("Parallax> that's all for now!")
			break
		} else {
			fmt.Println("Parallax> dont know what to do!")
		}
	}
}

func (h *Handler) Connect(server string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("Error connecting", server)
		return
	}
	defer conn.Close()
	h.Run(conn)
}

// Dummy Player

type FirstEdges struct{}

func (*FirstEdges) ComputeBid(m *Match) *BidPack {
	g := Instance(m.instanceName)
	if g == nil {
		fmt.Println("Instance not found:", m.instanceName)
		return EmptyBidPack()
	}
	pack := NewBidPack(m.edges)
	for i := 0; i < m.edges; i++ {
		e := g.edges[i]
		pack.Bid(e.i.id, e.j.id, e.vCost)
	}
	return pack
}

func (*FirstEdges) Update(f *Flow) {
}

func main() {
	fmt.Println("Game Theory Player: Parallax Engine")

	flag.Parse()

	fmt.Println("Threads:", *optThreads)
	runtime.GOMAXPROCS(*optThreads)

	if *optPreload {
		LoadAllInstances()
	}

	engine := &FirstEdges{}
	handler := NewHandler(*optName, engine)
	handler.Connect(*optServer)
}
