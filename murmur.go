// Package murmur implements the 32-bit Murmur hash, version 3.
package murmur

import (
	"reflect"
	"unsafe"
)

// Murmur32 is a 32-bit hash code generator using the Murmur hash 2.1
// algorithm.  It implements the hash.Hash32 interface.
type Murmur32 struct {
	hash, seed, k, full uint32
	bytes               int
}

const (
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
	r1 = 15
	r2 = 13
	m1 = 5
	n1 = 0xe6546b64
)

func New32(seed uint32) *Murmur32 {
	return &Murmur32{hash: seed, seed: seed}
}

func (m *Murmur32) Write(data []byte) (int, error) {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	if m.bytes&3 == 0 {
		return m.writeFast(header.Data, header.Len)
	} else {
		return m.writeSlow(data)
	}
}

func (m *Murmur32) writeFast(data uintptr, size int) (int, error) {
	hash := m.hash
	term := data + uintptr(size)

	for ; data+4 <= term; data += 4 {
		val := *(*uint32)(unsafe.Pointer(data))
		val = val * c1
		val = (val << r1) | (val >> (32 - r1))
		val = val * c2
		hash = hash ^ val
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = hash*m1 + n1
	}

	full := hash
	for ; data < term; data++ {
		b := uint32(*(*byte)(unsafe.Pointer(data)))
		m.k = uint32(b)<<24 | m.k>>8
		rem := b * c1
		rem = (rem << r1) | (rem >> (32 - r1))
		rem = rem * c2
		full = full ^ rem
	}
	m.bytes += size
	full = full ^ uint32(m.bytes)
	full = full ^ (full >> 16)
	full = full * 0x85ebca6b
	full = full ^ (full >> 13)
	full = full * 0xc2b2ae35
	full = full ^ (full >> 16)

	m.hash, m.full = hash, full
	return size, nil
}

func (m *Murmur32) writeSlow(data []byte) (int, error) {
	hash, k, bytes := m.hash, m.k, m.bytes

	for _, b := range data {
		k = uint32(b)<<24 | k>>8
		bytes++
		if bytes&3 == 0 {
			k = k * c1
			k = (k << r1) | (k >> (32 - r1))
			k = k * c2
			hash = hash ^ k
			hash = (hash << r2) | (hash >> (32 - r2))
			hash = hash*m1 + n1
		}
	}

	full := hash
	count := bytes & 3
	for i, v := count, k; i > 0; i-- {
		rem := ((v >> 24) & 0xff) * c1
		rem = (rem << r1) | (rem >> (32 - r1))
		rem = rem * c2
		full = full ^ rem
		v = v << 8
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

func (m *Murmur32) Reset() {
	m.hash = m.seed
}

func (m *Murmur32) Sum32() uint32 {
	return m.full
}

func (m *Murmur32) Size() int {
	return 4
}

func (m *Murmur32) BlockSize() int {
	return 1
}

func (m *Murmur32) Sum(in []byte) []byte {
	v := m.full
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
