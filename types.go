package main

import (
	"github.com/apaxa-go/helper/mathh"
	"github.com/apaxa-go/helper/strconvh"
	"strings"
)

type IPv4 uint32

func (ip IPv4) String() string {
	o1 := uint8(ip >> 24)
	o2 := uint8((ip >> 16) & mathh.MaxUint8)
	o3 := uint8((ip >> 8) & mathh.MaxUint8)
	o4 := uint8(ip & mathh.MaxUint8)
	return strconvh.FormatUint8(o1) + "." + strconvh.FormatUint8(o2) + "." + strconvh.FormatUint8(o3) + "." + strconvh.FormatUint8(o4)
}

func (ip IPv4) Netmask(last IPv4) Prefix {
	for m := Prefix(32); m > 0; m-- {
		if ip == last {
			return m
		}
		ip = ip >> 1
		last = last >> 1
	}
	return 0
}

type Prefix uint8 // 0..32 for IPv4
const RemoveMask Prefix = 255

func (m Prefix) String() string {
	return strconvh.FormatUint8(uint8(m))
}
func (m Prefix) Mask() IPv4 {
	return IPv4(mathh.MaxUint32) << (32 - m)
}

type NETv4 struct {
	IP     IPv4
	Prefix Prefix
}

func (n NETv4) String() string {
	return n.IP.String() + "/" + n.Prefix.String()
}

func (n *NETv4) Normalize() {
	n.IP &= n.Prefix.Mask()
}

func (n NETv4) Contains(n1 NETv4) bool {
	return n.Prefix <= n1.Prefix && n.IP&n.Prefix.Mask() == n1.IP&n.Prefix.Mask()
}

func (n NETv4) Count() uint64 {
	return 1 << (32 - n.Prefix)
}

func (n NETv4) First() IPv4 {
	return n.IP & n.Prefix.Mask()
}

func (n NETv4) Last() IPv4 {
	return n.IP | ^n.Prefix.Mask()
}

// Last meaning bit is 0
func (n NETv4) FirstInGroup() bool {
	return n.Prefix != 0 && n.IP&(1<<(32-n.Prefix)) == 0
}

// Last meaning bit is 1
func (n NETv4) LastInGroup() bool {
	return n.Prefix != 0 && n.IP&(1<<(32-n.Prefix)) > 0
}

func (n NETv4) GroupPair() NETv4 {
	if n.Prefix == 0 {
		return NETv4{}
	}
	var pair NETv4
	pair.IP = n.IP ^ (1 << (32 - n.Prefix))
	pair.Prefix = n.Prefix
	return pair
}

func (n NETv4) SummaryMask(last NETv4) Prefix {
	return n.IP.Netmask(last.Last())
}

type NETv4s []NETv4

func (n NETv4s) First() int {
	i := 0
	for i < len(n) && n[i].Prefix == RemoveMask {
		i++
	}
	return i
}

func (n NETv4s) Last() int {
	i := len(n) - 1
	for i >= 0 && n[i].Prefix == RemoveMask {
		i--
	}
	return i
}

func (n NETv4s) Next(i int) int {
	return i + 1 + n[i+1:].First()
}

func (n NETv4s) Prev(i int) int {
	return n[:i].Last()
}

func (n NETv4s) Count() uint64 {
	count := uint64(0)
	for i := range n {
		count += n[i].Count()
	}
	return count
}

func (n NETv4s) string(sep string) string {
	s := make([]string, len(n))
	for i := range s {
		s[i] = n[i].String()
	}
	return strings.Join(s, sep)
}

func (n NETv4s) String() string {
	return n.string(", ")
}

func (n NETv4s) StringNL() string {
	return n.string("\n")
}
