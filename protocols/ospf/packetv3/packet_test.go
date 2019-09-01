package packetv3_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	ospf "github.com/bio-routing/bio-rd/protocols/ospf/packetv3"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

func TestDecode(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir += "/fixtures/"

	files := []string{
		"OSPFv3_multipoint_adjacencies.cap",
		"OSPFv3_broadcast_adjacency.cap",
		"OSPFv3_NBMA_adjacencies.cap",
	}

	for _, path := range files {
		t.Run(path, func(t *testing.T) {
			testDecodeFile(t, dir+path)
		})
	}
}

func testDecodeFile(t *testing.T, path string) {
	fmt.Printf("Testing on file: %s\n", path)
	f, err := os.Open(path)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	r, err := pcapgo.NewReader(f)
	if err != nil {
		t.Error(err)
		return
	}

	var packetCount int
	for {
		data, _, err := r.ReadPacketData()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			return
		}

		t.Run(fmt.Sprintf("Packet %d", packetCount+1), func(t *testing.T) {
			payload, err := getPayload(data)
			if err != nil {
				t.Error(err)
				return
			}

			if err := handlePacket(payload); err != nil {
				t.Error(err)
			}
		})
		packetCount++
	}
}

func getPayload(raw []byte) ([]byte, error) {
	packet := gopacket.NewPacket(raw, layers.LayerTypeEthernet, gopacket.Default)
	if err := packet.ErrorLayer(); err != nil {
		// fallback to handling of FrameRelay (cut-off header)
		packet = gopacket.NewPacket(raw[4:], layers.LayerTypeIPv6, gopacket.Default)
		if err := packet.ErrorLayer(); err != nil {
			return []byte{}, err.Error()
		}
	}

	return packet.NetworkLayer().LayerPayload(), nil
}

func handlePacket(payload []byte) error {
	buf := bytes.NewBuffer(payload)
	_, _, err := ospf.DeserializeOSPFv3Message(buf)
	return err
}
