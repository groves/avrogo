include $(GOROOT)/src/Make.inc

TARG=avrogo
GOFILES=\
	primitives.go \
	schema.go \
	container.go

include $(GOROOT)/src/Make.pkg

