###############################################################################
#
#	Welcome to the most awesome Makefile you will ever see 😏
#
###############################################################################
#--------------------------------------
# Variables
#--------------------------------------
APP := polygon
PROJ := github.com/jeffizhungry/polygon
PROJPATH := $(shell pwd)
SOURCES := $(shell find $(PROJPATH) -name '*.go')

#--------------------------------------
# Targets
#--------------------------------------
$(APP): $(SOURCES)
	go build -o $(APP) .

#--------------------------------------
# Basic Rules
#--------------------------------------
all: $(APP)

clean:
	rm -f $(APP)

distclean: clean

#--------------------------------------
# Custom Rules
#--------------------------------------
.PHONY: run
run: $(APP)
	./$(APP)

.PHONY: updatedeps
updatedeps:
	godep restore
	godep save ./...
	go install ./...

.PHONY: startlocal
startlocal:
	@echo '-- Starting local services'

.PHONY: stoplocal
stoplocal:
	@echo '-- Stopping local services'

.PHONY: test1 
test1:
	curl -XPOST -d'{"s":"Hello, World"}' localhost:8008/toLower

.PHONY: test2
test2:
	curl -XPOST -d'{"s":"Hello, World"}' localhost:8008/toUpper

.PHONY: test3
test3:
	curl -XPOST -d'{"s":"Hello, World"}' localhost:8008/length


#--------------------------------------
# Testing Rules
#--------------------------------------

# Shortcut for local testing
.PHONY: testall
testall: clean all unittest vet errcheck ineffassign misspell

.PHONY: unittest
unittest:
	go test -v -timeout 30s $(PROJ)

.PHONY: vet
vet:
	go vet $(PROJ)

.PHONY: errcheck
errcheck:
	errcheck $(PROJ)

.PHONY: misspell
misspell:
	misspell $(PROJPATH)

.PHONY: ineffassign
ineffassign:
	ineffassign $(PROJPATH)

.PHONY: unittestwrace
unittestwrace:
	go test -v -race -timeout 30s $(PROJ)
