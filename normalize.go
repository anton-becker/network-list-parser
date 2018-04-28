package main

import (
	"sort"
)

func NormalizeNETs(nets []NETv4) {
	for i := range nets {
		nets[i].Normalize()
	}
}

func PackNETs(nets []NETv4) []NETv4 {
	removeCount := 0
	for i := range nets {
		if nets[i].Prefix == RemoveMask {
			removeCount++
		}
	}

	r := make([]NETv4, 0, len(nets)-removeCount)
	for i := range nets {
		if nets[i].Prefix != RemoveMask {
			r = append(r, nets[i])
		}
	}

	return r
}

func EliminateNETs(nets NETv4s) {
	for i := 0; i < len(nets); {
		j := i + 1
		for j < len(nets) && nets[i].Contains(nets[j]) {
			nets[j].Prefix = RemoveMask
			j++
		}
		i = j
	}
}

func GroupUpNETs(nets NETv4s) {
	for i := nets.First(); i < len(nets); {
		next := nets.Next(i)
		if !nets[i].FirstInGroup() || next >= len(nets) {
			i = next
			continue
		}
		if nets[i].GroupPair() == nets[next] {
			nets[i].Prefix--
			nets[next].Prefix = RemoveMask
			if nets[i].LastInGroup() {
				prev := nets.Prev(i)
				if prev >= 0 {
					i = prev
				} else {
					i = nets.Next(next)
				}
			}
		} else {
			i = next
		}
	}
}

func NormalizeIPs(nets []NETv4) []NETv4 {
	NormalizeNETs(nets)

	less := func(i, j int) bool {
		iv := int64(nets[i].IP) << 8
		iv |= int64(nets[i].Prefix)
		jv := int64(nets[j].IP) << 8
		jv |= int64(nets[j].Prefix)
		return iv < jv
	}

	sort.Slice(nets, less)

	EliminateNETs(nets)
	GroupUpNETs(nets)
	nets = PackNETs(nets)

	return nets
}
