// Copyright 2019 The gVisor Authors.
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

package header

import "gvisor.dev/gvisor/pkg/tcpip"

const (
	// NDPTargetLinkLayerAddressOptionType is the type of the Target
	// Link-Layer Address option, as per RFC 4861 section 4.6.1.
	NDPTargetLinkLayerAddressOptionType = 2

	// ndpTargetEthernetLinkLayerAddressSize is the size of a Target
	// Link Layer Option for an Ethernet address.
	ndpTargetEthernetLinkLayerAddressSize = 8

	// lengthByteUnits is the multiplier factor for the Length field of an
	// NDP option. That is, the length field for NDP options is in units of
	// 8 octets, as per RFC 4861 section 4.6.
	lengthByteUnits = 8
)

// NDPOptions is a buffer of NDP options as defined by RFC 4861 section 4.6.
type NDPOptions []byte

// Serialize serializes the provided list of NDP options into o.
func (b NDPOptions) Serialize(s NDPOptionsSerializer) {
	if s.Length() == 0 {
		return
	}

	done := 0
	for _, o := range s {
		l := o.Length()

		if l == 0 {
			continue
		}

		bytes := int(l) * lengthByteUnits

		buf := b[done:][:bytes]
		buf[0] = o.Type()
		buf[1] = l
		o.Serialize(buf[2:])

		done += bytes
	}
}

// NDPOption is the set of functions to be implemented by all NDP option types.
type NDPOption interface {
	// Type returns the type of this NDPOption.
	Type() uint8

	// Length returns the length of this NDPOption, including the bytes for
	// the Type and Length fields, in units of lengthByteUnits bytes..
	Length() uint8

	// Serialize serializes this NDPOption into the provided byte buffer.
	Serialize([]byte)
}

// NDPOptionsSerializer is a serializer for NDP options.
type NDPOptionsSerializer []NDPOption

// Length returns the total number of bytes required to serialize.
func (b NDPOptionsSerializer) Length() uint16 {
	l := uint8(0)

	for _, o := range b {
		l += o.Length()
	}

	// o.Length() returns bytes in lengthByteUnits units so we multiply by
	// lengthByteUnits to get the total number of bytes.
	return uint16(l) * lengthByteUnits
}

// NDPTargetLinkLayerAddressOption is the NDP Target Link Layer Option
// as defined by RFC 4861 section 4.6.1.
type NDPTargetLinkLayerAddressOption tcpip.LinkAddress

// Type implements NdpOption.Type.
func (o NDPTargetLinkLayerAddressOption) Type() uint8 {
	return NDPTargetLinkLayerAddressOptionType
}

// Length implements NdpOption.Length.
func (o NDPTargetLinkLayerAddressOption) Length() uint8 {
	if len(o) == 0 {
		return 0
	}

	// Length includes the 2 Type and Length bytes.
	l := len(o) + 2

	// Add extra bytes if needed to make sure the option is
	// lengthByteUnits-byte aligned. We do this by adding lengthByteUnits-1
	// to l and then stripping off the last few LSBits from l. This will
	// make sure that l is rounded up to the nearest unit of
	// lengthByteUnits. This works since lengthByteUnits is a power of 2
	// (= 8).
	mask := lengthByteUnits - 1
	l += mask
	l &^= mask

	if l > 255 {
		// Should never happen, so just return 0 so this option does not
		// get serialized.
		//
		// Returning 0 here will make sure that this option does not get
		// serialized when NDPOptions.Serialize is called with the
		// NDPOptionsSerializer that holds this option, effectively
		// skipping this option during serialization. Also note that
		// a value of zero for the Length field in an NDP option is
		// invalid so this is another sign to the caller that this NDP
		// option is malformed, as per RFC 4861 section 4.6.
		return 0
	}

	// Return length in units of lengthByteUnits bytes.
	return uint8(l / lengthByteUnits)
}

// Serialize implements NdpOption.Serialize.
func (o NDPTargetLinkLayerAddressOption) Serialize(b []byte) {
	copy(b[:len(o)], o)

	// Zero out remaining (padding) bytes, if any exists.
	for i := len(o); i < len(b); i++ {
		b[i] = 0
	}
}
