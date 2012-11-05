// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

// Package mathutil provides utilities supplementing the standard 'math' and
// 'math/rand' packages.
package mathutil

import (
	"math"
	"math/big"
)

var (
	_1 = big.NewInt(1)
	_2 = big.NewInt(2)
)

// GCDByte returns the greatest common divisor of a and b. Based on:
// http://en.wikipedia.org/wiki/Euclidean_algorithm#Implementations
func GCDByte(a, b byte) byte {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// GCDUint16 returns the greatest common divisor of a and b.
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

// ISqrt returns floor(sqrt(n)). Typical run time is few hundreds of ns.
func ISqrt(n uint32) (x uint32) {
	if n == 0 {
		return
	}

	if n >= math.MaxUint16*math.MaxUint16 {
		return math.MaxUint16
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

// SqrtUint64 returns floor(sqrt(n)). Typical run time is about 0.5 µs.
func SqrtUint64(n uint64) (x uint64) {
	if n == 0 {
		return
	}

	if n >= math.MaxUint32*math.MaxUint32 {
		return math.MaxUint32
	}

	var px, nx uint64
	for x = n; ; px, x = x, nx {
		nx = (x + n/x) / 2
		if nx == x || nx == px {
			break
		}
	}
	return
}

// SqrtBig returns floor(sqrt(n)). It panics on n < 0.
func SqrtBig(n *big.Int) (x *big.Int) {
	switch n.Sign() {
	case -1:
		panic(-1)
	case 0:
		return big.NewInt(0)
	}

	var px, nx big.Int
	x = big.NewInt(0)
	x.SetBit(x, n.BitLen()/2+1, 1)
	for {
		nx.Rsh(nx.Add(x, nx.Div(n, x)), 1)
		if nx.Cmp(x) == 0 || nx.Cmp(&px) == 0 {
			break
		}
		px.Set(x)
		x.Set(&nx)
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

// ModPowByte computes (b^e)%m. It panics for m == 0 || b == e == 0.
//
// See also: http://en.wikipedia.org/wiki/Modular_exponentiation#Right-to-left_binary_method
func ModPowByte(b, e, m byte) byte {
	if b == 0 && e == 0 {
		panic(0)
	}

	if m == 1 {
		return 0
	}

	r := uint16(1)
	for b, m := uint16(b), uint16(m); e > 0; b, e = b*b%m, e>>1 {
		if e&1 == 1 {
			r = r * b % m
		}
	}
	return byte(r)
}

// ModPowByte computes (b^e)%m. It panics for m == 0 || b == e == 0.
func ModPowUint16(b, e, m uint16) uint16 {
	if b == 0 && e == 0 {
		panic(0)
	}

	if m == 1 {
		return 0
	}

	r := uint32(1)
	for b, m := uint32(b), uint32(m); e > 0; b, e = b*b%m, e>>1 {
		if e&1 == 1 {
			r = r * b % m
		}
	}
	return uint16(r)
}

// ModPowUint32 computes (b^e)%m. It panics for m == 0 || b == e == 0.
func ModPowUint32(b, e, m uint32) uint32 {
	if b == 0 && e == 0 {
		panic(0)
	}

	if m == 1 {
		return 0
	}

	r := uint64(1)
	for b, m := uint64(b), uint64(m); e > 0; b, e = b*b%m, e>>1 {
		if e&1 == 1 {
			r = r * b % m
		}
	}
	return uint32(r)
}

// ModPowUint64 computes (b^e)%m. It panics for m == 0 || b == e == 0.
func ModPowUint64(b, e, m uint64) (r uint64) {
	if b == 0 && e == 0 {
		panic(0)
	}

	if m == 1 {
		return 0
	}

	r, _ = Uint64FromBigInt(modPowBigInt(Uint64ToBigInt(b), Uint64ToBigInt(e), Uint64ToBigInt(m)))
	return
}

func modPowBigInt(b, e, m *big.Int) (r *big.Int) {
	r = big.NewInt(1)
	for i, n := 0, e.BitLen(); i < n; i++ {
		if e.Bit(i) != 0 {
			r.Mod(r.Mul(r, b), m)
		}
		b.Mod(b.Mul(b, b), m)
	}
	return
}

// ModPowBigInt computes (b^e)%m. Returns nil for e < 0. It panics for m == 0 || b == e == 0.
func ModPowBigInt(b, e, m *big.Int) (r *big.Int) {
	if b.Sign() == 0 && e.Sign() == 0 {
		panic(0)
	}

	if m.Cmp(_1) == 0 {
		return big.NewInt(0)
	}

	if e.Sign() < 0 {
		return
	}

	return modPowBigInt(big.NewInt(0).Set(b), big.NewInt(0).Set(e), m)
}

var uint64ToBigIntDelta big.Int

func init() {
	uint64ToBigIntDelta.SetBit(&uint64ToBigIntDelta, 63, 1)
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

// AddUint128_64 returns the uint128 sum of uint64 a and b.
func AddUint128_64(a, b uint64) (hi uint64, lo uint64) {
	lo = a + b
	if lo < a {
		hi = 1
	}
	return
}

// MulUint128_64 returns the uint128 bit product of uint64 a and b.
func MulUint128_64(a, b uint64) (hi, lo uint64) {
	/*
		2^(2 W) ahi bhi + 2^W alo bhi + 2^W ahi blo + alo blo

		FEDCBA98 76543210 FEDCBA98 76543210
		                  ---- alo*blo ----
		         ---- alo*bhi ----
		         ---- ahi*blo ----
		---- ahi*bhi ----
	*/
	const w = 32
	const m = 1<<w - 1
	ahi, bhi, alo, blo := a>>w, b>>w, a&m, b&m
	lo = alo * blo
	mid1 := alo * bhi
	mid2 := ahi * blo
	c1, lo := AddUint128_64(lo, mid1<<w)
	c2, lo := AddUint128_64(lo, mid2<<w)
	_, hi = AddUint128_64(ahi*bhi, mid1>>w+mid2>>w+uint64(c1+c2))
	return
}

// PowerizeBigInt returns (e, p) such that e is the smallest number for which p
// == b^e is greater or equal n. For n < 0 or b < 2 (0, nil) is returned.
//
// NOTE: Run time for large values of n (above about 2^1e6 ~= 1e300000) can be
// significant and/or unacceptabe.  For any smaller values of n the function
// typically performs in sub second time.  For "small" values of n (cca bellow
// 2^1e3 ~= 1e300) the same can be easily below 10 µs.
//
// A special (and trivial) case of b == 2 is handled separately and performs
// much faster.
func PowerizeBigInt(b, n *big.Int) (e uint32, p *big.Int) {
	switch {
	case b.Cmp(_2) < 0 || n.Sign() < 0:
		return
	case n.Sign() == 0 || n.Cmp(_1) == 0:
		return 0, big.NewInt(1)
	case b.Cmp(_2) == 0: //TODO actually all powers of 2 should use the short path
		p = big.NewInt(0)
		e = uint32(n.BitLen() - 1)
		p.SetBit(p, int(e), 1)
		if p.Cmp(n) < 0 {
			p.Mul(p, _2)
			e++
		}
		return
	}

	bw := b.BitLen()
	nw := n.BitLen()
	p = big.NewInt(1)
	var bb, r big.Int
	for {
		switch p.Cmp(n) {
		case -1:
			x := uint32((nw - p.BitLen()) / bw)
			if x == 0 {
				x = 1
			}
			e += x
			switch x {
			case 1:
				p.Mul(p, b)
			default:
				r.Set(_1)
				bb.Set(b)
				e := x
				for {
					if e&1 != 0 {
						r.Mul(&r, &bb)
					}
					if e >>= 1; e == 0 {
						break
					}

					bb.Mul(&bb, &bb)
				}
				p.Mul(p, &r)
			}
		case 0, 1:
			return
		}
	}
	panic("unreachable")
}

// PowerizeUint32BigInt returns (e, p) such that e is the smallest number for
// which p == b^e is greater or equal n. For n < 0 or b < 2 (0, nil) is
// returned.
//
// More info: see PowerizeBigInt.
func PowerizeUint32BigInt(b uint32, n *big.Int) (e uint32, p *big.Int) {
	switch {
	case b < 2 || n.Sign() < 0:
		return
	case n.Sign() == 0 || n.Cmp(_1) == 0:
		return 0, big.NewInt(1)
	case b == 2: //TODO actually all powers of 2 should use the short path
		p = big.NewInt(0)
		e = uint32(n.BitLen() - 1)
		p.SetBit(p, int(e), 1)
		if p.Cmp(n) < 0 {
			p.Mul(p, _2)
			e++
		}
		return
	}

	var bb big.Int
	bb.SetInt64(int64(b))
	return PowerizeBigInt(&bb, n)
}

/*
ProbablyPrimeUint32 returns true if n is prime or n is a pseudoprime to base a.
It implements the Miller-Rabin primality test for one specific value of 'a' and
k == 1.

Wrt pseudocode shown at
http://en.wikipedia.org/wiki/Miller-Rabin_primality_test#Algorithm_and_running_time

 Input: n > 3, an odd integer to be tested for primality;
 Input: k, a parameter that determines the accuracy of the test
 Output: composite if n is composite, otherwise probably prime
 write n − 1 as 2^s·d with d odd by factoring powers of 2 from n − 1
 LOOP: repeat k times:
    pick a random integer a in the range [2, n − 2]
    x ← a^d mod n
    if x = 1 or x = n − 1 then do next LOOP
    for r = 1 .. s − 1
       x ← x^2 mod n
       if x = 1 then return composite
       if x = n − 1 then do next LOOP
    return composite
 return probably prime

... this function behaves like passing 1 for 'k' and additionaly a
fixed/non-random 'a'.  Otherwise it's the same algorithm.

See also: http://mathworld.wolfram.com/Rabin-MillerStrongPseudoprimeTest.html
*/
func ProbablyPrimeUint32(n, a uint32) bool {
	d, s := n-1, 0
	for ; d&1 == 0; d, s = d>>1, s+1 {
	}
	x := uint64(ModPowUint32(a, d, n))
	if x == 1 || uint32(x) == n-1 {
		return true
	}

	for ; s > 1; s-- {
		if x = x * x % uint64(n); x == 1 {
			return false
		}

		if uint32(x) == n-1 {
			return true
		}
	}
	return false
}

// ProbablyPrimeUint64_32 returns true if n is prime or n is a pseudoprime to
// base a. It implements the Miller-Rabin primality test for one specific value
// of 'a' and k == 1.  See also ProbablyPrimeUint32.
func ProbablyPrimeUint64_32(n uint64, a uint32) bool {
	d, s := n-1, 0
	for ; d&1 == 0; d, s = d>>1, s+1 {
	}
	x := ModPowUint64(uint64(a), d, n)
	if x == 1 || x == n-1 {
		return true
	}

	bx, bn := Uint64ToBigInt(x), Uint64ToBigInt(n)
	for ; s > 1; s-- {
		if x, _ = Uint64FromBigInt(bx.Mod(bx.Mul(bx, bx), bn)); x == 1 {
			return false
		}

		if x == n-1 {
			return true
		}
	}
	return false
}

// ProbablyPrimeBigInt_32 returns true if n is prime or n is a pseudoprime to
// base a. It implements the Miller-Rabin primality test for one specific value
// of 'a' and k == 1.  See also ProbablyPrimeUint32.
func ProbablyPrimeBigInt_32(n *big.Int, a uint32) bool {
	var d big.Int
	d.Set(n)
	d.Sub(&d, _1) // d <- n-1
	s := 0
	for ; d.Bit(s) == 0; s++ {
	}
	nMinus1 := big.NewInt(0).Set(&d)
	d.Rsh(&d, uint(s))

	x := ModPowBigInt(big.NewInt(int64(a)), &d, n)
	if x.Cmp(_1) == 0 || x.Cmp(nMinus1) == 0 {
		return true
	}

	for ; s > 1; s-- {
		if x = x.Mod(x.Mul(x, x), n); x.Cmp(_1) == 0 {
			return false
		}

		if x.Cmp(nMinus1) == 0 {
			return true
		}
	}
	return false
}

// ProbablyPrimeBigInt returns true if n is prime or n is a pseudoprime to base
// a. It implements the Miller-Rabin primality test for one specific value of
// 'a' and k == 1.  See also ProbablyPrimeUint32.
func ProbablyPrimeBigInt(n, a *big.Int) bool {
	var d big.Int
	d.Set(n)
	d.Sub(&d, _1) // d <- n-1
	s := 0
	for ; d.Bit(s) == 0; s++ {
	}
	nMinus1 := big.NewInt(0).Set(&d)
	d.Rsh(&d, uint(s))

	x := ModPowBigInt(a, &d, n)
	if x.Cmp(_1) == 0 || x.Cmp(nMinus1) == 0 {
		return true
	}

	for ; s > 1; s-- {
		if x = x.Mod(x.Mul(x, x), n); x.Cmp(_1) == 0 {
			return false
		}

		if x.Cmp(nMinus1) == 0 {
			return true
		}
	}
	return false
}
