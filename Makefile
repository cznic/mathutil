# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include ../../../../Make.inc

TARG=github.com/cznic/mathutil

GOFILES=\
		big.go\
		permute.go\
		primes.go\
		rnd.go\

include ../../../../Make.pkg
