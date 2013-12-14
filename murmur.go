// Package murmur implements the 32-bit Murmur hash, version 3.
package murmur

import (
	"hash"
)

type murmur32 struct {
	hash, seed, k, full uint32
	bytes               int
}

const (
	C1 = 0xcc9e2d51
	C2 = 0x1b873593
	R1 = 15
	R2 = 13
	M  = 5
	N  = 0xe6546b64
)

func New32(seed uint32) hash.Hash32 {
	return &murmur32{hash: seed, seed: seed}
}

func (m *murmur32) Write(data []byte) (int, error) {
	hash, k, bytes := m.hash, m.k, m.bytes

	for _, b := range data {
		k = k<<8 | uint32(b)
		bytes++
		if bytes&3 == 0 {
			k = k * C1
			k = (k << R1) | (k >> (32 - R1))
			k = k * C2
			hash = hash ^ k
			hash = (hash << R2) | (hash >> (32 - R2))
			hash = hash*M + N
		}
	}

	full := hash
	for i, c, v := 0, m.bytes&3, k; i < c; i, v = i+1, v>>8 {
		rem := uint32(v&0xff) * C1
		rem = (rem << R1) | (rem >> (32 - R1))
		rem = rem * C2
		full = full ^ rem
	}
	full = full ^ uint32(bytes)
	full = full ^ (full >> 16)
	full = full * 0x85ebca6b
	full = full ^ (full >> 13)
	full = full * 0xc2b2ae35
	full = full ^ (full >> 16)

	m.hash, m.full, m.k, m.bytes = hash, full, k, bytes
	return len(data), nil
}

func (m *murmur32) Reset() {
	m.hash = m.seed
}

func (m *murmur32) Sum32() uint32 {
	return m.full
}

func (m *murmur32) Size() int {
	return 4
}

func (m *murmur32) BlockSize() int {
	return 1
}

func (m *murmur32) Sum(in []byte) []byte {
	v := m.full
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
