include .env

SHELL = /usr/bin/env bash

.PHONY: all build install uninstall clean

all: build

build: .docker_image
	docker run \
        --rm \
        -v $${PWD}:/utbotgo \
        -w /utbotgo/utils \
        utbotgo:${UTBOTGO_VERSION} \
        make

.docker_image: Dockerfile .env
	docker build -t utbotgo:$(UTBOTGO_VERSION) .
	touch .docker_image

install: $(UTBOTGO_BIN_DIR)/utbotgo

$(UTBOTGO_BIN_DIR)/utbotgo: utbotgo.sh .env build
	sed -e "s/##utbotgo_dir/$${PWD//\//\\\/}/" < $< > $@
	chmod ugo+x $@

uninstall:
	rm -f $(UTBOTGO_BIN_DIR)/utbotgo

clean: uninstall
	cd utils; $(MAKE) clean
	docker rmi utbotgo:$(UTBOTGO_VERSION)
	rm -f .docker_image
