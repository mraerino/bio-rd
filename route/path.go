package route

import (
	"fmt"

	bnet "github.com/bio-routing/bio-rd/net"
)

// Path represents a network path
type Path struct {
	Type        uint8
	StaticPath  *StaticPath
	BGPPath     *BGPPath
	NetlinkPath *NetlinkPath
}

// Select returns negative if p < q, 0 if paths are equal, positive if p > q
func (p *Path) Select(q *Path) int8 {
	switch {
	case p == nil && q == nil:
		return 0
	case p == nil:
		return -1
	case q == nil:
		return 1
	default:
	}

	if p.Type > q.Type {
		return 1
	}

	if p.Type < q.Type {
		return -1
	}

	switch p.Type {
	case BGPPathType:
		return p.BGPPath.Select(q.BGPPath)
	case StaticPathType:
		return p.StaticPath.Select(q.StaticPath)
	case NetlinkPathType:
		return p.NetlinkPath.Select(q.NetlinkPath)
	}

	panic("Unknown path type")
}

// ECMP checks if path p and q are equal enough to be considered for ECMP usage
func (p *Path) ECMP(q *Path) bool {
	switch p.Type {
	case BGPPathType:
		return p.BGPPath.ECMP(q.BGPPath)
	case StaticPathType:
		return p.StaticPath.ECMP(q.StaticPath)
	case NetlinkPathType:
		return p.NetlinkPath.ECMP(q.NetlinkPath)
	}

	panic("Unknown path type")
}

// Equal checks if paths p and q are equal
func (p *Path) Equal(q *Path) bool {
	if p == nil || q == nil {
		return false
	}

	if p.Type != q.Type {
		return false
	}

	switch p.Type {
	case BGPPathType:
		return p.BGPPath.Equal(q.BGPPath)
	case StaticPathType:
		return p.StaticPath.Equal(q.StaticPath)
	}

	return p.Select(q) == 0
}

// PathsDiff gets the list of elements contained by a but not b
func PathsDiff(a, b []*Path) []*Path {
	ret := make([]*Path, 0)

	for _, pa := range a {
		if !pathsContains(pa, b) {
			ret = append(ret, pa)
		}
	}

	return ret
}

func pathsContains(needle *Path, haystack []*Path) bool {
	for _, p := range haystack {
		if p == needle {
			return true
		}
	}

	return false
}

// Print all known information about a route in logfile friendly format
func (p *Path) String() string {
	switch p.Type {
	case StaticPathType:
		return "not implemented yet"
	case BGPPathType:
		return p.BGPPath.String()
	case NetlinkPathType:
		return p.NetlinkPath.String()
	default:
		return "Unknown paty type. Probably not implemented yet"
	}
}

// Print all known information about a route in human readable form
func (p *Path) Print() string {
	protocol := ""
	switch p.Type {
	case StaticPathType:
		protocol = "static"
	case BGPPathType:
		protocol = "BGP"
	case NetlinkPathType:
		protocol = "Netlink"
	}

	ret := fmt.Sprintf("\tProtocol: %s\n", protocol)
	switch p.Type {
	case StaticPathType:
		ret += "Not implemented yet"
	case BGPPathType:
		ret += p.BGPPath.Print()
	case NetlinkPathType:
		ret += p.NetlinkPath.Print()
	}

	return ret
}

// Copy a route
func (p *Path) Copy() *Path {
	if p == nil {
		return nil
	}

	cp := *p
	cp.BGPPath = cp.BGPPath.Copy()
	cp.StaticPath = cp.StaticPath.Copy()

	return &cp
}

// NextHop returns the next hop IP Address
func (p *Path) NextHop() bnet.IP {
	switch p.Type {
	case BGPPathType:
		return p.BGPPath.NextHop
	case StaticPathType:
		return p.StaticPath.NextHop
	case NetlinkPathType:
		return p.NetlinkPath.NextHop
	}

	panic("Unknown path type")
}
