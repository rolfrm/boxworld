# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=run
GOFILES=\
	math32.go\
	Vector.go\
	BSP.go\
	ABSP.go\
	physics.go\
	graphics.go\
	utils.go\
	functional.go\
	main.go\

include $(GOROOT)/src/Make.cmd
