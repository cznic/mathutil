// Copyright (c) 2011 jnml. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mersenne collects utilities related to Mersenne numbers[1] and/or some
of their properties.

Exponent

In this documentatoin the term 'exponent' refers to 'n' of a Mersenne number Mn
equal to 2^n-1. This package supports only uint32 sized exponents. New()
currently supports exponents only up to math.MaxInt32 (31 bits, up to 256 MB
required to represent such Mn in memory as a big.Int).

Links

Referenced from above:
 [1] http://en.wikipedia.org/wiki/Mersenne_number
*/
package mersenne

import (
	"github.com/cznic/mathutil"
	"math"
	"math/big"
)

var (
	_0 = big.NewInt(0)
	_1 = big.NewInt(1)
	_2 = big.NewInt(2)
)

// New returns Mn == 2^n-1 for n <= math.MaxInt32 or nil otherwise.
func New(n uint32) (m *big.Int) {
	if n > math.MaxInt32 {
		return
	}

	m = big.NewInt(0)
	return m.Sub(m.SetBit(m, int(n), 1), _1)
}

// HasFactorUint32 returns true if d | Mn. Typical run time for a 32 bit factor
// and a 32 bit exponent is < 1 µs.
func HasFactorUint32(d, n uint32) bool {
	return d == 1 || d&1 != 0 && mathutil.ModPowUint32(2, n, d) == 1
}

// HasFactorUint64 returns true if d | Mn. Typical run time for a 64 bit factor
// and a 32 bit exponent is < 30 µs.
func HasFactorUint64(d uint64, n uint32) bool {
	return d == 1 || d&1 != 0 && mathutil.ModPowUint64(2, uint64(n), d) == 1
}

// HasFactorBigInt returns true if d | Mn, d > 0. Typical run time for a 128
// bit factor and a 32 bit exponent is < 75 µs.
func HasFactorBigInt(d *big.Int, n uint32) bool {
	return d.Cmp(_1) == 0 || d.Sign() > 0 && d.Bit(0) == 1 &&
		mathutil.ModPowBigInt(_2, big.NewInt(int64(n)), d).Cmp(_1) == 0
}

/*
FromFactorBigInt returns n such that d | Mn if n <= max and d is odd. In other
cases zero is returned.

Every odd d ∊ N divides infinitely many Mersenne numbers. The returned n is the
exponent of smallest such Mn.

NOTE: The computation of n from a given d performs in roughly O(n). It is
thus highly recomended to use the 'max' argument to limit the "searched"
exponent upper bound as appropriate. Otherwise the computation can take a long
time as a large factor can be a divisor of a Mn with exponent above the uint32
limits.

The FromFactorBigInt function is a modification of the original Will
Edgington's "reverse method", discussed here:
http://tech.groups.yahoo.com/group/primenumbers/message/15061
*/
func FromFactorBigInt(d *big.Int, max uint32) (n uint32) {
	if d.Bit(0) == 0 {
		return
	}

	var m big.Int
	for n < max {
		m.Add(&m, d)
		i := 0
		for ; m.Bit(i) == 1; i++ {
			if n == math.MaxUint32 {
				return 0
			}

			n++
		}
		m.Rsh(&m, uint(i))
		if m.Sign() == 0 {
			if n > max {
				n = 0
			}
			return
		}
	}
	return 0
}
