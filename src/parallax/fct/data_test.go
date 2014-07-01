package fct

import (
	"testing"
)

func TestLoadData(t *testing.T) {
	g, err := LoadGraph("N104.DAT", 1)
	if err != nil {
		t.Fatal("Error loading FCT data:", err)
	}
	if n := g.SourceOrder(); n != 10 {
		t.Error("Wrong number of sources (10):", n)
	}
	if n := g.SinkOrder(); n != 10 {
		t.Error("Wrong number of sinks (10):", n)
	}
	if n := g.Size(); n != 100 {
		t.Error("Wrong number of arks (100):", n)
	}
}
