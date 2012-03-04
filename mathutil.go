// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

// Package mathutil provides utilities supplementing the standard 'math' and 'rand' packages.
package mathutil

import (
	"math"
	"math/big"
)

// GCDByte returns the greatest common divisor of a and b.
//
// Based on: http://en.wikipedia.org/wiki/Euclidean_algorithm#Implementations
func GCDByte(a, b byte) byte {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// GCD16 returns the greatest common divisor of a and b.
func GCDUint16(a, b uint16) uint16 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// GCD returns the greatest common divisor of a and b.
func GCDUint32(a, b uint32) uint32 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// GCD64 returns the greatest common divisor of a and b.
func GCDUint64(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// ISqrt returns floor(sqrt(n)). Typical run time is few hundred nsecs.
func ISqrt(n uint32) (x uint32) {
	if n == 0 {
		return
	}

	if n >= math.MaxUint16*math.MaxUint16 {
		return 65535
	}

	var px, nx uint32
	for x = n; ; px, x = x, nx {
		nx = (x + n/x) / 2
		if nx == x || nx == px {
			break
		}
	}
	return
}

// Log2Byte returns log base 2 of n. It's the same as index of the highest
// bit set in n.  For n == 0 -1 is returned.
func Log2Byte(n byte) int {
	return log2[n]
}

// Log2Uint16 returns log base 2 of n. It's the same as index of the highest
// bit set in n.  For n == 0 -1 is returned.
func Log2Uint16(n uint16) int {
	if b := n >> 8; b != 0 {
		return log2[b] + 8
	}

	return log2[n]
}

// Log2Uint32 returns log base 2 of n. It's the same as index of the highest
// bit set in n.  For n == 0 -1 is returned.
func Log2Uint32(n uint32) int {
	if b := n >> 24; b != 0 {
		return log2[b] + 24
	}

	if b := n >> 16; b != 0 {
		return log2[b] + 16
	}

	if b := n >> 8; b != 0 {
		return log2[b] + 8
	}

	return log2[n]
}

// Log2Uint64 returns log base 2 of n. It's the same as index of the highest
// bit set in n.  For n == 0 -1 is returned.
func Log2Uint64(n uint64) int {
	if b := n >> 56; b != 0 {
		return log2[b] + 56
	}

	if b := n >> 48; b != 0 {
		return log2[b] + 48
	}

	if b := n >> 40; b != 0 {
		return log2[b] + 40
	}

	if b := n >> 32; b != 0 {
		return log2[b] + 32
	}

	if b := n >> 24; b != 0 {
		return log2[b] + 24
	}

	if b := n >> 16; b != 0 {
		return log2[b] + 16
	}

	if b := n >> 8; b != 0 {
		return log2[b] + 8
	}

	return log2[n]
}

// ModPowByte computes (b^e)%m. It panics for m == 0.
//
// See also: http://en.wikipedia.org/wiki/Modular_exponentiation#Right-to-left_binary_method
func ModPowByte(b, e, m byte) byte {
	r := uint16(1)
	for b, m := uint16(b), uint16(m); e > 0; b, e = b*b%m, e>>1 {
		if e&1 == 1 {
			r = r * b % m
		}
	}
	return byte(r)
}

// ModPowByte computes (b^e)%m. It panics for m == 0.
func ModPowUint16(b, e, m uint16) uint16 {
	r := uint32(1)
	for b, m := uint32(b), uint32(m); e > 0; b, e = b*b%m, e>>1 {
		if e&1 == 1 {
			r = r * b % m
		}
	}
	return uint16(r)
}

// ModPowUint32 computes (b^e)%m. It panics for m == 0.
func ModPowUint32(b, e, m uint32) uint32 {
	r := uint64(1)
	for b, m := uint64(b), uint64(m); e > 0; b, e = b*b%m, e>>1 {
		if e&1 == 1 {
			r = r * b % m
		}
	}
	return uint32(r)
}

// ModPowUint64 computes (b^e)%m. It panics for m == 0.
func ModPowUint64(b, e, m uint64) (r uint64) {
	r, _ = Uint64FromBigInt(modPowBigInt(Uint64ToBigInt(b), Uint64ToBigInt(e), Uint64ToBigInt(m)))
	return
}

func modPowBigInt(b, e, m *big.Int) (r *big.Int) {
	r = big.NewInt(1)
	for e.Sign() > 0 {
		if e.Bit(0) == 1 {
			r.Mod(r.Mul(r, b), m)
		}
		e.Rsh(e, 1)
		b.Mod(b.Mul(b, b), m)
	}
	return
}

// ModPowBigInt computes (b^e)%m. Returns nil for e < 0. It panics for m == 0.
func ModPowBigInt(b, e, m *big.Int) (r *big.Int) {
	if e.Sign() < 0 {
		return
	}

	return modPowBigInt(big.NewInt(0).Set(b), big.NewInt(0).Set(e), m)
}

var uint64ToBigIntDelta big.Int

func init() {
	uint64ToBigIntDelta.SetString("8000000000000000", 16)
}

// Uint64ToBigInt returns a big.Int set to n.
func Uint64ToBigInt(n uint64) *big.Int {
	if n <= math.MaxInt64 {
		return big.NewInt(int64(n))
	}

	y := big.NewInt(int64(n - uint64(math.MaxInt64) - 1))
	return y.Add(y, &uint64ToBigIntDelta) // Benchamrked to be faster than y.SetBit(y, 63, 1)
}

// Uint64FromBigInt returns (uint64 value of n, true) if 0 <= n <=
// math.MaxUint64.  Otherwise it returns  (undefined value, false).
func Uint64FromBigInt(n *big.Int) (uint64, bool) {
	switch bits := n.BitLen(); {
	case bits == 0:
		return 0, true
	case n.Sign() < 0 || bits > 64:
		return 0, false
	case bits <= UintptrBits():
		return uint64(n.Bits()[0]), true
	default:
		b := n.Bits()
		return uint64(b[1])<<uint(uintptrBits) | uint64(b[0]), true
	}
	panic("unreachable")
}

var uintptrBits int

func init() {
	x := uint64(math.MaxUint64)
	uintptrBits = BitLenUintptr(uintptr(x))
}

// UintptrBits returns the bit width of an uintptr at the executing machine.
func UintptrBits() int {
	return uintptrBits
}
