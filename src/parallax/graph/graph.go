package graph

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FCTP data model

type Vertex struct {
	Id    int
	Size  float64
	Edges []*Edge
}

func (v *Vertex) String() string {
	return fmt.Sprintf("(%d, [%.2f,%d])", v.Id, v.Size, len(v.Edges))
}

type Edge struct {
	I, J         *Vertex
	VCost, FCost float64
}

func (e *Edge) String() string {
	return fmt.Sprintf("(%d)-[%.2f,%.2f]->(%d)", e.I.Id, e.VCost, e.FCost, e.J.Id)
}

type Graph struct {
	Sources, Sinks map[int]*Vertex
	Edges          []*Edge
}

func (g *Graph) String() string {
	return fmt.Sprintf("Sources %d, Sinks %d, Edges %d", len(g.Sources), len(g.Sinks), len(g.Edges))
}

func NewGraph() *Graph {
	return &Graph{
		Sources: make(map[int]*Vertex),
		Sinks:   make(map[int]*Vertex),
		Edges:   make([]*Edge, 0),
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
	vi := g.v(g.Sources, i)
	vj := g.v(g.Sinks, j)
	e := &Edge{vi, vj, v, f}
	vi.Edges = append(vi.Edges, e)
	vj.Edges = append(vj.Edges, e)
	g.Edges = append(g.Edges, e)
	return e
}

func (g *Graph) SourceSize(i int, s float64) *Vertex {
	vi := g.v(g.Sources, i)
	vi.Size = s
	return vi
}

func (g *Graph) SinkSize(j int, s float64) *Vertex {
	vj := g.v(g.Sinks, j)
	vj.Size = s
	return vj
}

func LoadGraph(path string, verbose int) (*Graph, error) {
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
		if verbose > 2 {
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
			if verbose > 1 {
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
			if verbose > 1 {
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
			if verbose > 1 {
				fmt.Println("Demand:", v)
			}
		} else {
			fmt.Println("Ignoring line, unknown type:", line)
		}
	}
	return g, nil
}

type GraphLoader struct {
	dataPath string
	cache    map[string]*Graph
	verbose  int
}

func NewLoader(dataPath string, verbose int) *GraphLoader {
	return &GraphLoader{dataPath, make(map[string]*Graph), verbose}
}

func (d *GraphLoader) Instance(name string) *Graph {
	if g, ok := d.cache[name]; ok {
		return g
	}
	fmt.Print("Loading ", name, "... ")
	g, err := LoadGraph(d.dataPath+"/"+name+".DAT", d.verbose)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	d.cache[name] = g
	fmt.Println(g)
	return g
}

func (d *GraphLoader) LoadAll() {
	folder, err := os.Open(d.dataPath)
	if err != nil {
		fmt.Println("Error opening folder:", d.dataPath, err)
		return
	}
	defer folder.Close()

	files, err := folder.Readdirnames(0)
	if err != nil {
		fmt.Println("Error listing data files from:", d.dataPath, err)
		return
	}

	for _, file := range files {
		if !strings.HasSuffix(file, ".DAT") {
			fmt.Println("Ignoring", file)
			continue
		}
		fmt.Print("Loading ", file, "... ")
		g, err := LoadGraph(d.dataPath+"/"+file, d.verbose)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		name := strings.TrimSuffix(file, ".DAT")
		d.cache[name] = g
		fmt.Println(g)
	}

	fmt.Println("Total:", len(d.cache))
}
