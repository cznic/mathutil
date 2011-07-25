// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz


package mathutil

import (
	"math"
)

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

// IsPrimes returns true if n is prime. Typical run time is few hundred nsecs.
//
// http://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test#Deterministic_variants_of_the_test
func IsPrime(n uint32) bool {
	switch {
	case n&1 == 0:
		return n == 2
	case n < 3:
		return false
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
