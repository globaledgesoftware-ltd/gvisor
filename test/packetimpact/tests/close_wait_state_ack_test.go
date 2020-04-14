// Copyright 2020 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package close_wait_state_Ack_test

import (
	"testing"
	"time"

	"fmt"
	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/seqnum"
	tb "gvisor.dev/gvisor/test/packetimpact/testbench"
)

func TestTCPCloseWaitStateAck(t *testing.T) {
	for _, tt := range []struct {
		description string
		isSeqNum    bool
	}{
		{"SEQNUM", true},
		{"ACKNUM", false},
	} {
		t.Run(fmt.Sprintf("%s", tt.description), func(t *testing.T) {
			println("\nCASE \n")
			if tt.isSeqNum == true {
				println("\nCASE OOW SEQNUM\n")
			} else {
				println("\nCASE UNACCEPT ACKNUM\n")
			}
			dut := tb.NewDUT(t)
			defer dut.TearDown()
			listenFd, remotePort := dut.CreateListener(unix.SOCK_STREAM, unix.IPPROTO_TCP, 1)
			defer dut.Close(listenFd)
			conn := tb.NewTCPIPv4(t, tb.TCP{DstPort: &remotePort}, tb.TCP{SrcPort: &remotePort})
			defer conn.Close()
			conn.Handshake()
			acceptFd, _ := dut.Accept(listenFd)

			// Send FIN-ACK to DUT.
			conn.Send(tb.TCP{Flags: tb.Uint8(header.TCPFlagFin | header.TCPFlagAck)})
			// Expecting ACK from DUT
			if _, err := conn.Expect(tb.TCP{Flags: tb.Uint8(header.TCPFlagAck)}, 3*time.Second); err != nil {
				t.Fatal("expected an ACK packet within 3 seconds but got nonei")
			}
			fmt.Println("In CLOSED_WAIT STATE")
			if tt.isSeqNum == true {
				conn.Send(tb.TCP{
					SeqNum: tb.Uint32(uint32(conn.LocalSeqNum.Add(seqnum.Size(*conn.SynAck.WindowSize) + 2))),
					Flags:  tb.Uint8(header.TCPFlagAck),
				}, []tb.Layer{&tb.Payload{Bytes: []byte("abcdef12345")}}...)
			} else {
				conn.Send(tb.TCP{
					AckNum: tb.Uint32(uint32(conn.LocalSeqNum.Add(seqnum.Size(*conn.SynAck.WindowSize) + 2))),
					Flags:  tb.Uint8(header.TCPFlagAck),
				}, []tb.Layer{&tb.Payload{Bytes: []byte("abcdef12345")}}...)
			}

			_, err := conn.Expect(tb.TCP{Flags: tb.Uint8(header.TCPFlagAck)}, 3*time.Second)
			if err != nil {
				t.Fatal("expected an ACK packet within 3 seconds but got none")
			}
			// Verifying that DUT is in the CLOSE-WAIT state Causing DUT to issue a CLOSE call
			dut.Close(acceptFd)
			if _, err := conn.Expect(tb.TCP{Flags: tb.Uint8(header.TCPFlagFin | header.TCPFlagAck)}, time.Second); err != nil {
				t.Fatal("expected an FIN-ACK packet within a second but got none")
			}
			conn.Send(tb.TCP{Flags: tb.Uint8(header.TCPFlagAck)})

			// Sending a TCP data packet to DUT and Expecting RST response from DUT
			conn.Send(tb.TCP{Flags: tb.Uint8(header.TCPFlagAck)}, []tb.Layer{&tb.Payload{Bytes: []byte("sample data")}}...)
			if _, err := conn.Expect(tb.TCP{Flags: tb.Uint8(header.TCPFlagRst)}, time.Second); err != nil {
				t.Fatal("expected an RSTpacket within a second but got none")
			}
			conn.Send(tb.TCP{Flags: tb.Uint8(header.TCPFlagRst | header.TCPFlagAck)})
		})
	}
}
