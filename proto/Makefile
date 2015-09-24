
DESTDIR     ?= .
PROTO_FILES := $(wildcard *.proto)
GO_DIR      := $(DESTDIR)/go
GO_FILES    := $(patsubst %.proto,$(GO_DIR)/%.pb.go,$(PROTO_FILES))
PY_DIR      := $(DESTDIR)/py
PY_FILES    := $(patsubst %.proto,$(PY_DIR)/%_pb2.py,$(PROTO_FILES))
C_DIR       := $(DESTDIR)/c
C_SOURCE    := $(patsubst %.proto,$(C_DIR)/%.pb-c.c,$(PROTO_FILES))
C_HDR       := $(patsubst %.proto,$(C_DIR)/%.pb-c.h,$(PROTO_FILES))
C_FILES     := $(C_SOURCE) $(C_HDR)
CPP_DIR     := $(DESTDIR)/cpp
CPP_SOURCE  := $(patsubst %.proto,$(CPP_DIR)/%.pb.cc,$(PROTO_FILES))
CPP_HDR     := $(patsubst %.proto,$(CPP_DIR)/%.pb.h,$(PROTO_FILES))
CPP_FILES   := $(CPP_SOURCE) $(CPP_HDR)

default: example

all: go py c cpp

go: $(GO_FILES)

$(GO_DIR)/%.pb.go: %.proto
	@mkdir -p $(GO_DIR)
	protoc --go_out=$(GO_DIR) $^

example: go
	go run ./example.go

c: $(C_FILES)

$(C_DIR)/%.pb-c.c: %.proto
	@mkdir -p $(C_DIR)
	protoc-c --c_out=$(C_DIR) $^

cpp: $(CPP_FILES)

$(CPP_DIR)/%.pb.cc: %.proto
	@mkdir -p $(CPP_DIR)
	protoc --cpp_out=$(CPP_DIR)/ $^

py: $(PY_FILES)

$(PY_DIR)/%_pb2.py: %.proto
	@mkdir -p $(PY_DIR)
	protoc --python_out=$(PY_DIR) $^

clean:
	rm -rf *~ $(GO_FILES) $(C_FILES) $(PY_FILES) $(CPP_FILES)

