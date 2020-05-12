# makefile
pro=$(shell pwd)
GOBIN = $(shell pwd)/bin

def:
	@echo "default compile start"
	go build -o $(GOBIN)/repe ./main
	cp -rf $(pro)/main/tmpl $(GOBIN)/ && cp -f $(pro)/main/repe.yaml $(GOBIN)/repe.yaml

linux:
	@echo "linux compile start"
	xgo --targets linux/amd64 --dest $(GOBIN) -out repe -pkg main  $(GOPATH)/src/repe

