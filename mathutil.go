// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

// Package mathutil provides utilities supplementing the standard 'math' and 'rand' packages.
package mathutil

import (
	"fmt"
	"math"
)

// FC32 is a full cycle PRNG covering the 32 bit signed integer range.
// In contrast to full cycle generators shown at e.g. http://en.wikipedia.org/wiki/Full_cycle,
// this code doesn't produce values at constant delta (mod cycle length).
// The 32 bit limit is per this implementation, the algorithm used has no intrinsic limit on the cycle size.
// Properties include:
//	- Adjustable limits on creation (hi, lo).
//	- Positionable/randomly accessible (Pos, Seek).
//	- Repeatable (deterministic).
//	- Can run forward or backward (Next, Prev).
//	- For a billion numbers cycle the Next/Prev PRN can be produced in cca 100-150ns.
//	  That's like 5-10 times slower compared to PRNs generated using the (non FC) rand package.
type FC32 struct {
	cycle   int64     // On average: 3 * delta / 2, (HQ: 2 * delta)
	delta   int64     // hi - lo
	factors [][]int64 // This trades some space for hopefully a bit of speed (multiple adding vs multiplying).
	lo      int
	mods    []int   // pos % set
	pos     int64   // Within cycle.
	primes  []int64 // Ordered. ∏ primes == cycle.
	set     []int64 // Reordered primes (magnitude order bases) according to seed.
}

// NewFC32 returns a newly created FC32 adjusted for the closed interval [lo, hi] or an Error if any.
// If hq == true then trade some generation time for improved (pseudo)randomness.
func NewFC32(lo, hi int, hq bool) (r *FC32, err error) {
	if lo > hi {
		return nil, fmt.Errorf("invalid range %d > %d", lo, hi)
	}

	delta := int64(hi) - int64(lo)
	if delta > math.MaxUint32 {
		return nil, fmt.Errorf("range out of int32 limits %d, %d", lo, hi)
	}

	// Find the primorial covering whole delta
	n, set, p := int64(1), []int64{}, uint32(2)
	if hq {
		p++
	}
	for {
		set = append(set, int64(p))
		n *= int64(p)
		if n > delta {
			break
		}
		p, _ = NextPrime(p)
	}

	// Adjust the set so n ∊ [delta, 2 * delta] (HQ: [delta, 3 * delta])
	// while keeping the cardinality of the set (correlates with the statistic "randomness quality")
	// at max, i.e. discard atmost one member.
	i := -1 // no candidate prime
	if n > 2*(delta+1) {
		for j, p := range set {
			q := n / p
			if q < delta+1 {
				break
			}

			i = j // mark the highest candidate prime set index
		}
	}
	if i >= 0 { // shrink the inner cycle
		n = n / set[i]
		set = delete(set, i)
	}
	r = &FC32{
		cycle:   n,
		delta:   delta,
		factors: make([][]int64, len(set)),
		lo:      lo,
		mods:    make([]int, len(set)),
		primes:  set,
	}
	r.Seed(1) // the default seed should be always non zero
	return
}

// Cycle reports the length of the inner FCPRNG cycle.
// Cycle is atmost the double (HQ: triple) of the generator period (hi - lo + 1).
func (r *FC32) Cycle() int64 {
	return r.cycle
}

// Next returns the first PRN after Pos.
func (r *FC32) Next() int {
	return r.step(1)
}

// Pos reports the current position within the inner cycle.
func (r *FC32) Pos() int64 {
	return r.pos
}

// Prev return the first PRN before Pos.
func (r *FC32) Prev() int {
	return r.step(-1)
}

// Seed uses the provided seed value to initialize the generator to a deterministic state.
// A zero seed produces a "canonical" generator with worse randomness than for most non zero seeds.
// Still, the FC property holds for any seed value.
func (r *FC32) Seed(seed int64) {
	u := uint64(seed)
	r.set = mix(r.primes, &u)
	n := int64(1)
	for i, p := range r.set {
		k := make([]int64, p)
		v := int64(0)
		for j := range k {
			k[j] = v
			v += n
		}
		n *= p
		r.factors[i] = mix(k, &u)
	}
}

// Seek sets Pos to |pos| % Cycle.
func (r *FC32) Seek(pos int64) {
	if pos < 0 {
		pos = -pos
	}
	pos %= r.cycle
	r.pos = pos
	for i, p := range r.set {
		r.mods[i] = int(pos % p)
	}
}

func (r *FC32) step(dir int) int {
	for { // avg loops per step: 3/2 (HQ: 2)
		y := int64(0)
		pos := r.pos
		pos += int64(dir)
		switch {
		case pos < 0:
			pos = r.cycle - 1
		case pos >= r.cycle:
			pos = 0
		}
		r.pos = pos
		for i, mod := range r.mods {
			mod += dir
			p := int(r.set[i])
			switch {
			case mod < 0:
				mod = p - 1
			case mod >= p:
				mod = 0
			}
			r.mods[i] = mod
			y += r.factors[i][mod]
		}
		if y <= r.delta {
			return int(y) + r.lo
		}
	}
	panic("unreachable")
}

func delete(set []int64, i int) (y []int64) {
	for j, v := range set {
		if j != i {
			y = append(y, v)
		}
	}
	return
}

func mix(set []int64, seed *uint64) (y []int64) {
	for len(set) != 0 {
		*seed = rol(*seed)
		i := int(*seed % uint64(len(set)))
		y = append(y, set[i])
		set = delete(set, i)
	}
	return
}

func rol(u uint64) (y uint64) {
	y = u << 1
	if int64(u) < 0 {
		y |= 1
	}
	return
}
