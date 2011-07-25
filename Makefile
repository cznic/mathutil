# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=github.com/cznic/mathutil

GOFILES=\
	big.go\
	permute.go\
	primes.go\
	rnd.go\

CLEANFILES += example example2 example3 rnd.png rand.png *~

include $(GOROOT)/src/Make.pkg

example: example.go $(GOFILES)
	make install clean && make -f example.make

example2: example2.go $(GOFILES)
	make install clean && make -f example2.make

example3: example3.go $(GOFILES)
	make install clean && make -f example2.make
