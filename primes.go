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

func isPrimeUint16(n uint16) (isPrime bool) {
	_, isPrime = findSmallPrime(n)
	return
}

// IsPrimeUint16 returns true if n is prime. Typical run time is few tens of nsecs.
func IsPrimeUint16(n uint16) bool {
	switch {
	case n < 2 || n&1 == 0:
		return n == 2
	case n%3 == 0:
		return n == 3
	case n%3 == 0:
		return n == 3
	case n%5 == 0:
		return n == 5
	case n%7 == 0:
		return n == 7
	case n%11 == 0:
		return n == 11
	case n%13 == 0:
		return n == 13 // Benchmarked optimum
	default:
		return isPrimeUint16(n)
	}
	panic("unreachable")
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

func probablyPrime(n, a uint32) bool {
	e, s := n-1, 0
	for ; e&1 == 0; e, s = e>>1, s+1 {
	}
	x := uint64(ModPowUint32(a, e, n))
	if x == 1 || uint32(x) == n-1 {
		return true
	}

	for ; s > 1; s-- {
		x = (x * x) % uint64(n)
		if x == 1 {
			return false
		}

		if uint32(x) == n-1 {
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
	case n%2 == 0:
		return n == 2
	case n%3 == 0:
		return n == 3
	case n%5 == 0:
		return n == 5
	case n%7 == 0:
		return n == 7
	case n%11 == 0:
		return n == 11
	case n%13 == 0:
		return n == 13
	case n%17 == 0:
		return n == 17
	case n%19 == 0:
		return n == 19
	case n%23 == 0:
		return n == 23
	case n%29 == 0:
		return n == 29
	case n%31 == 0:
		return n == 31
	case n%37 == 0:
		return n == 37
	case n%41 == 0:
		return n == 41
	case n%43 == 0:
		return n == 43
	case n%47 == 0:
		return n == 47
	case n%53 == 0:
		return n == 53 // Benchmarked optimum
	case n < 65536:
		// use table data
		return isPrimeUint16(uint16(n))
	case n < 1373653:
		// it is enough to test a = 2 and 3
		return probablyPrime(n, 2) && probablyPrime(n, 3)
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
