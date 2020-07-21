GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

VERSION=$(shell git describe --exact-match --tags 2>/dev/null)
BUILD_DIR=build
PACKAGE_RPI=geizhalsinflux-$(VERSION)_linux_armhf

unexport GOPATH

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(GOBUILD)

package: build
	tar -cvzf $(PACKAGE_RPI).tar.gz -C $(BUILD_DIR) $(PACKAGE_RPI)

build:
	GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $(BUILD_DIR)/$(PACKAGE_RPI)/usr/bin/geizhalsinflux -i cmd/geizhalsinflux/main.go