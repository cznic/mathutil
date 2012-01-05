// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
)

/* 

$ # in "cznic/mathutil" dir
$ make -f example3.make
$ # Intended usage e.g.:
$ ./example3 -max 1024 > rand.dat # generate 1kB of "random" data

*/
func main() {
	r := rand.New(rand.NewSource(1))
	var mflag uint64
	flag.Uint64Var(&mflag, "max", 0, "limit output to max bytes")
	flag.Parse()
	stdout := bufio.NewWriter(os.Stdout)
	if mflag != 0 {
		for i := uint64(0); i < mflag; i++ {
			if err := stdout.WriteByte(byte(r.Int())); err != nil {
				log.Fatal(err)
			}
		}
		stdout.Flush()
	} else {
		for stdout.WriteByte(byte(r.Int())) == nil {
		}
	}
}
