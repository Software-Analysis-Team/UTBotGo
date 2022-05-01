include .env

SHELL = /usr/bin/env bash

.PHONY: all build install uninstall clean

all: build

build: .docker_image
	docker run \
        --rm \
        -v $${PWD}:/utbotgo \
        -w /utbotgo/utils \
        $(UTBOTGO_DOCKER_IMAGE) \
        make

.docker_image: Dockerfile .env
	docker pull $(UTBOTGO_DOCKER_IMAGE)
	touch .docker_image

install: $(UTBOTGO_BIN_DIR)/utbotgo

$(UTBOTGO_BIN_DIR)/utbotgo: utbotgo.sh .env build
	sed -e "s/##utbotgo_dir/$${PWD//\//\\\/}/" < $< > $@
	chmod ugo+x $@

uninstall:
	rm -f $(UTBOTGO_BIN_DIR)/utbotgo

clean: uninstall
	cd utils; $(MAKE) clean
	docker rmi $(UTBOTGO_DOCKER_IMAGE)
	rm -f .docker_image
