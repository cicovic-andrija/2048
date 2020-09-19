# Makefile for 2048
#
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINFNAME=2048
BINLOCAL=bin
BINUSER=$(HOME)/bin

.PHONY: install
install: build
	cp $(BINLOCAL)/$(BINFNAME) $(BINUSER)/$(BINFNAME)

.PHONY: build
build: prebuild
	$(GOBUILD) -o $(BINLOCAL)/$(BINFNAME) -v ./cmd

.PHONY: prebuild
prebuild:
	mkdir -p $(BINLOCAL)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BINLOCAL)
