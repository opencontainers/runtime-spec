
DESTDIR     ?= .
PROTO_FILES := $(wildcard *.proto)
GO_FILES    := $(patsubst %.proto,%.pb.go,$(PROTO_FILES))
C_FILES     := $(patsubst %.proto,%.pb-c.c,$(PROTO_FILES))
C_HDR_FILES := $(patsubst %.proto,%.pb-c.h,$(PROTO_FILES))
PY_FILES    := $(patsubst %.proto,%_pb2.py,$(PROTO_FILES))

default: go

all: go py c

go: $(GO_FILES)

%.pb.go: %.proto
	protoc --go_out=$(DESTDIR) $^

c: $(C_FILES)

%.pb-c.c: %.proto
	protoc-c --c_out=$(DESTDIR) $^

py: $(PY_FILES)

%_pb2.py: %.proto
	protoc --python_out=$(DESTDIR) $^


clean:
	rm -rf *~ $(GO_FILES) $(C_FILES) $(C_HDR_FILES) $(PY_FILES)

