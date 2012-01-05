// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

package mathutil

import (
	"math"
)

func findSmallPrime(n uint16) (index int, found bool) {
	x := n >> 8
	l, h := int(smallLimits[x]), int(smallLimits[x+1])
	for l <= h {
		index := (l + h) / 2
		p := smallPrimes[index]
		switch {
		case n > p:
			l = index + 1
		case n < p:
			h = index - 1
		default:
			return index, true
		}
	}
	return l, false
}

// IsPrimeUint16 returns true if n is prime. Typical run time is few hundred nsecs.
func IsPrimeUint16(n uint16) (isPrime bool) {
	_, isPrime = findSmallPrime(n)
	return
}

// NextPrimeUint16 returns first prime > n and true if successful or an undefined value and false if there
// is no next prime in the uint16 limits. Typical run time is few µsecs.
func NextPrimeUint16(n uint16) (p uint16, ok bool) {
	if n >= 65521 {
		return
	}
	i, _ := findSmallPrime(n + 1)
	return smallPrimes[i], true
}

func probablyPrime(n32, a32 uint32) bool {
	n, a := uint64(n32), uint64(a32)
	d, s := n-1, 0
	for ; d&1 == 0; d, s = d>>1, s+1 {
	}
	// http://en.wikipedia.org/wiki/Modular_exponentiation#Right-to-left_binary_method
	x := uint64(1)
	for ; d > 0; a, d = (a*a)%n, d>>1 {
		if d&1 != 0 {
			x = (x * a) % n
		}
	}
	if x == 1 || x == n-1 {
		return true
	}

	for ; s > 1; s-- {
		x = (x * x) % n
		if x == 1 {
			return false
		}

		if x == n-1 {
			return true
		}
	}
	return false
}

// IsPrime returns true if n is prime. Typical run time is few hundred nsecs.
//
// http://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test#Deterministic_variants_of_the_test
func IsPrime(n uint32) bool {
	switch {
	case n&1 == 0:
		return n == 2
	case n < 65536:
		// use table data
		return IsPrimeUint16(uint16(n))
	case n < 1373653:
		// it is enough to test a = 2 and 3
		return n == 3 || probablyPrime(n, 2) && probablyPrime(n, 3)
	case n < 9080191:
		// it is enough to test a = 31 and 73
		return probablyPrime(n, 31) && probablyPrime(n, 73)
	default: // case n < 4759123141
		// it is enough to test a = 2, 7, and 61
		return probablyPrime(n, 2) && probablyPrime(n, 7) && probablyPrime(n, 61)
	}
	panic("unreachable")
}

// NextPrime returns first prime > n and true if successful or an undefined value and false if there
// is no next prime in the uint32 limits. Typical run time is few µsecs.
func NextPrime(n uint32) (p uint32, ok bool) {
	switch {
	case n < 2:
		return 2, true
	case n == 2:
		return 3, true
	case n < 65521:
		p16, _ := NextPrimeUint16(uint16(n))
		return uint32(p16), true
	case n >= math.MaxUint32-4:
		return
	}

	n++
	var d0, d uint32
	switch mod := n % 6; mod {
	case 0:
		d0, d = 1, 4
	case 1:
		d = 4
	case 2, 3, 4:
		d0, d = 5-mod, 2
	case 5:
		d = 2
	}

	p = n + d0
	if p < n { // overflow
		return
	}

	for {
		if ok = IsPrime(p); ok {
			break
		}

		p0 := p
		p += d
		if p < p0 { // overflow
			break
		}

		d ^= 6
	}
	return
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

// FactorTerm is one term of an integer factorization.
type FactorTerm struct {
	Prime uint32 // The divisor
	Power uint32 // Term == Prime^Power
}

// FactorInt returns prime factorization of n > 1 or nil otherwise.
// Resulting factors are ordered by Prime. Typical run time is < 10 µsec.
func FactorInt(n uint32) (f []FactorTerm) {
	switch {
	case n < 2:
		return
	case IsPrime(n):
		return []FactorTerm{{n, 1}}
	}

	f, w := make([]FactorTerm, 9), 0
	for _, prime16 := range smallPrimes {
		prime := uint32(prime16)
		if prime*prime > n {
			break
		}

		power := uint32(0)
		for n%prime == 0 {
			n /= prime
			power++
		}
		if power != 0 {
			f[w] = FactorTerm{prime, power}
			w++
		}
		if n == 1 {
			break
		}
	}
	if n != 1 {
		f[w] = FactorTerm{n, 1}
		w++
	}
	return f[:w]
}
