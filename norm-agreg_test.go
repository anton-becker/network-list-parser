package main

import (
	"os"
	"testing"
)

const source = "dump.csv"

func TestNormalizeAndAutoAggregate(t *testing.T) {
	src, err := os.Open(source)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer src.Close()

	origNets := ParseIPs(src)
	nets := NormalizeIPs(origNets)
	settings := AutoAggregationDefaultSettings
	settings.LogMaxPrefix = 32
	AutoAggregate(nets, settings)
	nets = PackNETs(nets)

orig:
	for i := range origNets {
		for j := range nets {
			if nets[j].Contains(origNets[i]) {
				continue orig
			}
		}
		t.Errorf("lost network: %v", origNets[i].String())
	}
}
