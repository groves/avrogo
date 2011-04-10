include $(GOROOT)/src/Make.inc

TARG=avrogo
GOFILES=\
	primitives.go \
	schema.go

include $(GOROOT)/src/Make.pkg

