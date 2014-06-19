package engine

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// Game Protocol

type Match struct {
	InstanceName  string
	NumberOfEdges int
}

func (m *Match) String() string {
	return fmt.Sprint(m.InstanceName, " ", m.NumberOfEdges)
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
	Streams []*Stream
}

type Stream struct {
	Source, Sink int
	Amount       float64
	Owner        string
	Price        float64
	NumberOfBids int
}

func (f *Flow) String() string {
	return fmt.Sprint("Number of Edges ", len(f.Streams))
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
	out := "bid\n"
	for _, b := range p.bids {
		out += b.String() + "\n"
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

func (pp ProffitSlice) String() string {
	out := ""
	for i, p := range pp {
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
	name    string
	engine  Engine
	verbose int
}

func NewHandler(name string, engine Engine, verbose int) *Handler {
	return &Handler{name, engine, verbose}
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
			n, err := ParseMatch(m)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if h.verbose > 0 {
				fmt.Println("Match:", n)
			}
			result := h.engine.ComputeBid(n)
			fmt.Println("Parallax>")
			_out := result.String()
			if h.verbose < 2 && len(_out) > 50 {
				_out = _out[0:50] + "..."
			}
			fmt.Println(_out)
			fmt.Fprint(conn, result)
		} else if strings.HasPrefix(m, "result") {
			n, err := ParseFlow(m, master)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if h.verbose > 0 {
				fmt.Println("Flow:", n)
			}
			h.engine.Update(n)
		} else if strings.HasPrefix(m, "end") {
			n, err := ParseProffits(m, master)
			if err != nil {
				fmt.Println(err)
				break
			}
			if h.verbose > 0 {
				fmt.Println("Proffit:", n)
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
