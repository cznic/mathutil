// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz


package main

import (
	"bytes"
	"github.com/cznic/mathutil"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"rand"
)

// see http://www.boallen.com/random-numbers.html
func main() {
	if r, err := mathutil.NewFC32(math.MinInt32, math.MaxInt32, true); err != nil {
		log.Fatal(err)
	} else {
		img := image.NewGray(512, 512)
		for y := 0; y < 512; y++ {
			for x := 0; x < 512; x++ {
				if r.Next()&1 != 0 {
					img.Set(x, y, image.White)
				}
			}
		}
		buf := bytes.NewBuffer(nil)
		if err := png.Encode(buf, img); err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile("rnd.png", buf.Bytes(), 0666); err != nil {
			log.Fatal(err)
		}
	}

	r := rand.New(rand.NewSource(0))
	img := image.NewGray(512, 512)
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			if r.Int()&1 != 0 {
				img.Set(x, y, image.White)
			}
		}
	}
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("rand.png", buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
}
