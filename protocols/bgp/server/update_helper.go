package server

import (
	"fmt"
	"strings"

	"github.com/bio-routing/bio-rd/protocols/bgp/packet"
	"github.com/bio-routing/bio-rd/route"
)

func pathAttribues(p *route.Path, fsm *FSM) (*packet.PathAttribute, error) {
	asPathPA, err := packet.ParseASPathStr(strings.TrimRight(fmt.Sprintf("%d %s", fsm.localASN, p.BGPPath.ASPath), " "))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse AS path: %v", err)
	}

	origin := &packet.PathAttribute{
		TypeCode: packet.OriginAttr,
		Value:    p.BGPPath.Origin,
		Next:     asPathPA,
	}

	nextHop := &packet.PathAttribute{
		TypeCode: packet.NextHopAttr,
		Value:    p.BGPPath.NextHop,
	}
	asPathPA.Next = nextHop

	localPref := &packet.PathAttribute{
		TypeCode: packet.LocalPrefAttr,
		Value:    p.BGPPath.LocalPref,
	}
	nextHop.Next = localPref

	if p.BGPPath != nil {
		err := addOptionalPathAttribues(p, localPref)

		if err != nil {
			return nil, err
		}
	}

	return origin, nil
}

func addOptionalPathAttribues(p *route.Path, parent *packet.PathAttribute) error {
	current := parent

	if len(p.BGPPath.LargeCommunities) > 0 {
		largeCommunities, err := packet.LargeCommunityAttributeForString(p.BGPPath.LargeCommunities)
		if err != nil {
			return fmt.Errorf("Could not create large community attribute: %v", err)
		}

		current.Next = largeCommunities
		current = largeCommunities
	}

	return nil
}