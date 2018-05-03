package main

import "log"

type AutoAggregationSettings struct {
	IntensiveMinPrefix uint
	IntensiveMinNets   int
	IntensiveMaxFake   uint64
	MaxFake            uint64
	LoFakePercent      float64 // 0.2, netCount==2
	HiFakePercent      float64 // 0.51, netCount>=HiFakePercentNets
	HiFakePercentNets  int     // 8
	LogMaxPrefix       int

	// Calculated field
	k float64
}

var AutoAggregationDefaultSettings = AutoAggregationSettings{
	LogMaxPrefix:       23,
	IntensiveMinPrefix: 24,
	IntensiveMinNets:   2,
	IntensiveMaxFake:   128,
	MaxFake:            1024,
	LoFakePercent:      0.2,
	HiFakePercent:      0.51,
	HiFakePercentNets:  8,
}

func (s AutoAggregationSettings) Validate() {
	if s.IntensiveMinPrefix > 31 {
		panic(`invalid "intensive-aggregation-min-prefix"`)
	}
	if s.IntensiveMinNets < 2 {
		panic(`invalid "intensive-aggregation-min-nets"`)
	}
	if s.LoFakePercent <= 0 {
		panic(`invalid "aggregation-lo-fake-percent"`)
	}
	if s.HiFakePercent <= s.LoFakePercent {
		panic(`invalid "aggregation-hi-fake-percent"`)
	}
	if s.HiFakePercentNets <= 2 {
		panic(`invalid "aggregation-hi-fake-percent-nets"`)
	}
	if s.LogMaxPrefix > 30 || s.LogMaxPrefix < -1 {
		panic(`invalid "log-aggregation-max-prefix"`)
	}
}

func (s *AutoAggregationSettings) Enabled() bool {
	return s.IntensiveMinPrefix < 31 || s.MaxFake > 0
}

func (s *AutoAggregationSettings) Calc() {
	s.k = (s.HiFakePercent - s.LoFakePercent) / float64((s.HiFakePercentNets-2)*(s.HiFakePercentNets-2))
}

func AutoAggregateDecision(net NETv4, netCount int, ipCount uint64, settings AutoAggregationSettings) bool {
	if netCount < 2 {
		return false
	}

	fakeCount := net.Count() - ipCount

	//
	// Intensive
	//
	if net.Prefix >= Prefix(settings.IntensiveMinPrefix) {
		return netCount > settings.IntensiveMinNets || fakeCount <= settings.IntensiveMaxFake
	}

	fakePercent := float64(fakeCount) / float64(net.Count())

	//
	// Casual
	//

	if fakeCount > settings.MaxFake {
		return false
	}
	if netCount >= settings.HiFakePercentNets {
		return fakePercent <= settings.HiFakePercent
	}
	return fakePercent <= settings.LoFakePercent+settings.k*float64((netCount-2)*(netCount-2))
}

func AutoAggregate(nets NETv4s, settings AutoAggregationSettings) {
	// decided
	// Unable to summarize into "/32" because there is only 1 ip.
	// If normalize does not group into "/31" then this "/31" contains only 1 net and there is no reason to summarize.
	var decided [31]IPv4

netLoop:
	for i := 0; i < len(nets); {
		for m := Prefix(0); m < Prefix(len(decided)); m++ {
			if nets[i].Last() <= decided[m] {
				continue
			}
			if nets[i].Prefix <= m {
				continue
			}

			net := nets[i]
			net.Prefix = m
			decided[m] = net.Last()

			ipCount := nets[i].Count()
			j := i + 1
			for ; j < len(nets) && net.Contains(nets[j]); j++ {
				ipCount += nets[j].Count()
			}

			if j-1 == i { // if there are no pairs for current "m" then there are no pairs for "m+1"
				break
			}

			{ // maximize mask keeping nets[i:j] in
				effectiveM := nets[i].SummaryMask(nets[j-1])
				for m < effectiveM {
					m++
					decided[m] = net.Last()
				}
				net.Prefix = effectiveM
			}
			net.Normalize()

			if AutoAggregateDecision(net, j-i, ipCount, settings) {
				if int(net.Prefix) <= settings.LogMaxPrefix {
					log.Printf("Aggregate %v nets into %v, original IPs %.2f%% (%v), original nets: %v.\n", j-i, net.String(), float32(ipCount)*100/float32(net.Count()), ipCount, nets[i:j].String())
				}
				nets[i] = net
				for k := i + 1; k < j; k++ {
					nets[k].Prefix = RemoveMask
				}
				i = j
				continue netLoop
			}
		}
		i++
	}
}
