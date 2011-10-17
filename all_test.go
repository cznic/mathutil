// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

package mathutil

import (
	"big"
	"math"
	"sort"
	"testing"
)

func benchmark1eN(b *testing.B, r *FC32) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r.Next()
	}
}

func Benchmark1e3(b *testing.B) {
	b.StopTimer()
	r, _ := NewFC32(0, 1e3, false)
	benchmark1eN(b, r)
}

func Benchmark1e6(b *testing.B) {
	b.StopTimer()
	r, _ := NewFC32(0, 1e6, false)
	benchmark1eN(b, r)
}

func Benchmark1e9(b *testing.B) {
	b.StopTimer()
	r, _ := NewFC32(0, 1e9, false)
	benchmark1eN(b, r)
}

func Test0(t *testing.T) {
	const N = 10000
	for n := 1; n < N; n++ {
		lo, hi := 0, n-1
		period := int64(hi) - int64(lo) + 1
		r, err := NewFC32(lo, hi, false)
		if err != nil {
			t.Fatal(err)
		}
		if r.Cycle()-period > period {
			t.Fatalf("Cycle exceeds 2 * period")
		}
	}
	for n := 1; n < N; n++ {
		lo, hi := 0, n-1
		period := int64(hi) - int64(lo) + 1
		r, err := NewFC32(lo, hi, true)
		if err != nil {
			t.Fatal(err)
		}
		if r.Cycle()-2*period > period {
			t.Fatalf("Cycle exceeds 3 * period")
		}
	}
}

func Test1(t *testing.T) {
	const (
		N = 360
		S = 3
	)
	for hq := 0; hq <= 1; hq++ {
		for n := 1; n < N; n++ {
			for seed := 0; seed < S; seed++ {
				lo, hi := -n, 2*n
				period := int64(hi - lo + 1)
				r, err := NewFC32(lo, hi, hq == 1)
				if err != nil {
					t.Fatal(err)
				}
				r.Seed(int64(seed))
				m := map[int]bool{}
				v := make([]int, period, period)
				p := make([]int64, period, period)
				for i := lo; i <= hi; i++ {
					x := r.Next()
					p[i-lo] = r.Pos()
					if x < lo || x > hi {
						t.Fatal("t1.0")
					}
					if m[x] {
						t.Fatal("t1.1")
					}
					m[x] = true
					v[i-lo] = x
				}
				for i := lo; i <= hi; i++ {
					x := r.Next()
					if x < lo || x > hi {
						t.Fatal("t1.2")
					}
					if !m[x] {
						t.Fatal("t1.3")
					}
					if x != v[i-lo] {
						t.Fatal("t1.4")
					}
					if r.Pos() != p[i-lo] {
						t.Fatal("t1.5")
					}
					m[x] = false
				}
				for i := lo; i <= hi; i++ {
					r.Seek(p[i-lo] + 1)
					x := r.Prev()
					if x < lo || x > hi {
						t.Fatal("t1.6")
					}
					if x != v[i-lo] {
						t.Fatal("t1.7")
					}
				}
			}
		}
	}
}

func Test2(t *testing.T) {
	const (
		N = 370
		S = 3
	)
	for hq := 0; hq <= 1; hq++ {
		for n := 1; n < N; n++ {
			for seed := 0; seed < S; seed++ {
				lo, hi := -n, 2*n
				period := int64(hi - lo + 1)
				r, err := NewFC32(lo, hi, hq == 1)
				if err != nil {
					t.Fatal(err)
				}
				r.Seed(int64(seed))
				m := map[int]bool{}
				v := make([]int, period, period)
				p := make([]int64, period, period)
				for i := lo; i <= hi; i++ {
					x := r.Prev()
					p[i-lo] = r.Pos()
					if x < lo || x > hi {
						t.Fatal("t2.0")
					}
					if m[x] {
						t.Fatal("t2.1")
					}
					m[x] = true
					v[i-lo] = x
				}
				for i := lo; i <= hi; i++ {
					x := r.Prev()
					if x < lo || x > hi {
						t.Fatal("t2.2")
					}
					if !m[x] {
						t.Fatal("t2.3")
					}
					if x != v[i-lo] {
						t.Fatal("t2.4")
					}
					if r.Pos() != p[i-lo] {
						t.Fatal("t2.5")
					}
					m[x] = false
				}
				for i := lo; i <= hi; i++ {
					s := p[i-lo] - 1
					if s < 0 {
						s = r.Cycle() - 1
					}
					r.Seek(s)
					x := r.Next()
					if x < lo || x > hi {
						t.Fatal("t2.6")
					}
					if x != v[i-lo] {
						t.Fatal("t2.7")
					}
				}
			}
		}
	}
}

func benchmarkBig1eN(b *testing.B, r *FCBig) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r.Next()
	}
}

func BenchmarkBig1e3(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e3)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkBig1e6(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e6)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkBig1e9(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e9)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkBig1e12(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e12)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkBig1e15(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e15)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

func BenchmarkBig1e18(b *testing.B) {
	b.StopTimer()
	hi := big.NewInt(0).SetInt64(1e18)
	r, _ := NewFCBig(big0, hi, false)
	benchmarkBig1eN(b, r)
}

var (
	big0 = big.NewInt(0)
	big1 = big.NewInt(1)
)

func TestBig0(t *testing.T) {
	const N = 7400
	lo := big.NewInt(0)
	hi := big.NewInt(0)
	period := big.NewInt(0)
	c := big.NewInt(0)
	for n := int64(1); n < N; n++ {
		hi.SetInt64(n - 1)
		period.Set(hi)
		period.Sub(period, lo)
		period.Add(period, big1)
		r, err := NewFCBig(lo, hi, false)
		if err != nil {
			t.Fatal(err)
		}
		if r.cycle.Cmp(period) < 0 {
			t.Fatalf("Period exceeds cycle")
		}
		c.Set(r.Cycle())
		c.Sub(c, period)
		if c.Cmp(period) > 0 {
			t.Fatalf("Cycle exceeds 2 * period")
		}
	}
	for n := int64(1); n < N; n++ {
		hi.SetInt64(n - 1)
		period.Set(hi)
		period.Sub(period, lo)
		period.Add(period, big1)
		r, err := NewFCBig(lo, hi, true)
		if err != nil {
			t.Fatal(err)
		}
		if r.cycle.Cmp(period) < 0 {
			t.Fatalf("Period exceeds cycle")
		}
		c.Set(r.Cycle())
		c.Sub(c, period)
		c.Sub(c, period)
		if c.Cmp(period) > 0 {
			t.Fatalf("Cycle exceeds 3 * period")
		}
	}
}

func TestBig1(t *testing.T) {
	const (
		N = 120
		S = 3
	)
	lo := big.NewInt(0)
	hi := big.NewInt(0)
	seek := big.NewInt(0)
	for hq := 0; hq <= 1; hq++ {
		for n := int64(1); n < N; n++ {
			for seed := 0; seed < S; seed++ {
				lo64 := -n
				hi64 := 2 * n
				lo.SetInt64(lo64)
				hi.SetInt64(hi64)
				period := hi64 - lo64 + 1
				r, err := NewFCBig(lo, hi, hq == 1)
				if err != nil {
					t.Fatal(err)
				}
				r.Seed(int64(seed))
				m := map[int64]bool{}
				v := make([]int64, period, period)
				p := make([]int64, period, period)
				for i := lo64; i <= hi64; i++ {
					x := r.Next().Int64()
					p[i-lo64] = r.Pos().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb1.0")
					}
					if m[x] {
						t.Fatal("tb1.1")
					}
					m[x] = true
					v[i-lo64] = x
				}
				for i := lo64; i <= hi64; i++ {
					x := r.Next().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb1.2")
					}
					if !m[x] {
						t.Fatal("tb1.3")
					}
					if x != v[i-lo64] {
						t.Fatal("tb1.4")
					}
					if r.Pos().Int64() != p[i-lo64] {
						t.Fatal("tb1.5")
					}
					m[x] = false
				}
				for i := lo64; i <= hi64; i++ {
					r.Seek(seek.SetInt64(p[i-lo64] + 1))
					x := r.Prev().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb1.6")
					}
					if x != v[i-lo64] {
						t.Fatal("tb1.7")
					}
				}
			}
		}
	}
}

func TestBig2(t *testing.T) {
	const (
		N = 120
		S = 3
	)
	lo := big.NewInt(0)
	hi := big.NewInt(0)
	seek := big.NewInt(0)
	for hq := 0; hq <= 1; hq++ {
		for n := int64(1); n < N; n++ {
			for seed := 0; seed < S; seed++ {
				lo64, hi64 := -n, 2*n
				lo.SetInt64(lo64)
				hi.SetInt64(hi64)
				period := hi64 - lo64 + 1
				r, err := NewFCBig(lo, hi, hq == 1)
				if err != nil {
					t.Fatal(err)
				}
				r.Seed(int64(seed))
				m := map[int64]bool{}
				v := make([]int64, period, period)
				p := make([]int64, period, period)
				for i := lo64; i <= hi64; i++ {
					x := r.Prev().Int64()
					p[i-lo64] = r.Pos().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb2.0")
					}
					if m[x] {
						t.Fatal("tb2.1")
					}
					m[x] = true
					v[i-lo64] = x
				}
				for i := lo64; i <= hi64; i++ {
					x := r.Prev().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb2.2")
					}
					if !m[x] {
						t.Fatal("tb2.3")
					}
					if x != v[i-lo64] {
						t.Fatal("tb2.4")
					}
					if r.Pos().Int64() != p[i-lo64] {
						t.Fatal("tb2.5")
					}
					m[x] = false
				}
				for i := lo64; i <= hi64; i++ {
					s := p[i-lo64] - 1
					if s < 0 {
						s = r.Cycle().Int64() - 1
					}
					r.Seek(seek.SetInt64(s))
					x := r.Next().Int64()
					if x < lo64 || x > hi64 {
						t.Fatal("tb2.6")
					}
					if x != v[i-lo64] {
						t.Fatal("tb2.7")
					}
				}
			}
		}
	}
}

func TestPermutations(t *testing.T) {
	data := sort.IntSlice{3, 2, 1}
	check := [][]int{
		{1, 2, 3},
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}
	i := 0
	for PermutationFirst(data); ; i++ {
		if i >= len(check) {
			t.Fatalf("too much permutations generated: %d > %d", i+1, len(check))
		}

		for j, v := range check[i] {
			got := data[j]
			if got != v {
				t.Fatalf("permutation %d:\ndata: %v\ncheck: %v\nexpected data[%d] == %d, got %d", i, data, check[i], j, v, got)
			}
		}

		if !PermutationNext(data) {
			if i != len(check)-1 {
				t.Fatal("permutations generated", i, "expected", len(check))
			}
			break
		}
	}
}

func TestIsPrime(t *testing.T) {
	const p4M = 283146 // # of primes < 4e6
	n := 0
	for i := uint32(0); i <= 4e6; i++ {
		if IsPrime(i) {
			n++
		}
	}
	t.Log(n)
	if n != p4M {
		t.Fatal(n)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(uint32(i))
	}
}

func BenchmarkNextPrime(b *testing.B) {
	var p uint32
	for i := 0; i < b.N; i++ {
		p, _ = NextPrime(p)
	}
}

func TestNextPrime(t *testing.T) {
	const p4M = 283146 // # of primes < 4e6
	n := 0
	var p uint32
	for {
		p, _ = NextPrime(p)
		if p >= 4e6 {
			break
		}
		n++
	}
	t.Log(n)
	if n != p4M {
		t.Fatal(n)
	}
}

func TestNextPrime2(t *testing.T) {
	type data struct {
		x  uint32
		y  uint32
		ok bool
	}
	tests := []data{
		{0, 2, true},
		{1, 2, true},
		{2, 3, true},
		{3, 5, true},
		{math.MaxUint32, 0, false},
		{math.MaxUint32 - 1, 0, false},
		{math.MaxUint32 - 2, 0, false},
		{math.MaxUint32 - 3, 0, false},
		{math.MaxUint32 - 4, 0, false},
		{math.MaxUint32 - 5, math.MaxUint32 - 4, true},
	}

	for _, test := range tests {
		y, ok := NextPrime(test.x)
		if ok != test.ok || ok && y != test.y {
			t.Fatalf("x %d, got y %d ok %t, expected y %d ok %t", test.x, y, ok, test.y, test.ok)
		}
	}
}
