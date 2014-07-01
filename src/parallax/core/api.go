package core

import (
	"fmt"
)

// Game Protocol - API

type Match struct {
	InstanceName  string
	NumberOfEdges int
}

func (m *Match) String() string {
	return fmt.Sprint(m.InstanceName, " ", m.NumberOfEdges)
}

type Flow struct {
	Streams []*Stream
}

func (f *Flow) String() string {
	return fmt.Sprint("Number of Edges ", len(f.Streams))
}

type Stream struct {
	Source, Sink int
	Amount       float64
	Owner        string
	Price        float64
	NumberOfBids int
}

func (s *Stream) String() string {
	return fmt.Sprintf("(%d)-[%.2f]->(%d) [%s, %.2f, %d]", s.Source, s.Amount, s.Sink, s.Owner, s.Price, s.NumberOfBids)
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
