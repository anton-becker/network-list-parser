package main

import (
	"bufio"
	"github.com/apaxa-go/helper/mathh"
	"io"
	"log"
	"strings"
)

func isD(b uint8) bool { return b >= '0' && b <= '9' }

// Does not check anything after, ParseIPv4 do it.
func parseOctet(s string) (octet uint8, i int) {
	var v uint16
	for i < len(s) && i < 3 && isD(s[i]) {
		v = v*10 + uint16(s[i]-'0')
		i++
	}
	if v > mathh.MaxUint8 {
		i = 0
	}
	octet = uint8(v)
	return
}

// Also check what following is not digit.
func parseMask4(s string) (mask Prefix, i int) {
	i = len(s)
	if i == 0 {
		return 0, 0
	}
	if !isD(s[0]) {
		return 0, 0
	}
	mask = Prefix(s[0] - '0')
	if i == 1 || !isD(s[1]) {
		return mask, 1
	}
	mask = mask*10 + Prefix(s[1]-'0')
	if mask > 32 {
		return 0, 0
	}
	if i == 2 || !isD(s[2]) {
		return mask, 2
	}
	return 0, 0
}

func ParseIPv4(s string) (ip IPv4, i int) {
	for j := 0; j < 4; j++ {
		octet, skip := parseOctet(s[i:])
		if skip == 0 {
			return 0, 0
		}
		ip = (ip << 8) | IPv4(octet)
		i += skip

		if j < 3 {
			if i >= len(s) || s[i] != '.' {
				return 0, 0
			}
			i++
		} else {
			if i < len(s) && s[i] >= '0' && s[i] <= '9' {
				return 0, 0
			}
		}
	}
	return
}

func ParseNETv4(s string) (net NETv4, i int) {
	net.IP, i = ParseIPv4(s)
	if i == 0 {
		return
	}
	if i >= len(s) { // Use default mask
		net.Prefix = 32
		return
	}
	switch s[i] {
	case ' ', ';', '|': // Use default mask
		net.Prefix = 32
		return
	case '/': // Parse mask
		i++
		var deltaI int
		net.Prefix, deltaI = parseMask4(s[i:])
		if deltaI == 0 {
			return NETv4{}, 0
		}
		i += deltaI
		return
	default:
		return NETv4{}, 0
	}
}

func ParseIPsRow(s string) (nets NETv4s, ok bool) {
	nets = make(NETv4s, 0, 100)

	//
	// At least one net should be defined
	//
	s = strings.TrimSpace(s)
	net, i := ParseNETv4(s)
	if i == 0 {
		return nil, false
	}
	nets = append(nets, net)
	s = strings.TrimSpace(s[i:])

	//
	// Other nets
	//
	for len(s) > 0 {
		switch s[0] {
		case '|':
			s = strings.TrimSpace(s[1:])
			net, i := ParseNETv4(s)
			if i == 0 {
				return nil, false
			}
			nets = append(nets, net)
			s = strings.TrimSpace(s[i:])
		case ';':
			return nets, true
		default:
			return nil, false
		}
	}
	return nets, true
}

// Apply some hacks for broken rows
func fixIPsRow(s string) string {
	if len(s) > 0 && s[0] == ';' {
		return s[1:]
	}
	return s
}

func ParseIPs(src io.Reader) (nets NETv4s) {
	const KnownException = "Updated:"

	nets = make(NETv4s, 0, 10000)
	//s := bufio.NewScanner(src)
	s := bufio.NewReader(src)

	var err error
	var line string
	for row := 0; err == nil; row++ {
		line, err = s.ReadString('\n')
		if err != nil && len(line) == 0 {
			break
		}
		if strings.HasPrefix(line, KnownException) {
			continue
		}
		rowNets, ok := ParseIPsRow(fixIPsRow(line))
		if !ok {
			log.Printf("Error in row %v\n", row+1)
			continue
		}
		nets = append(nets, rowNets...)
	}

	if err != nil && err != io.EOF {
		panic(err.Error())
	}

	return
}
