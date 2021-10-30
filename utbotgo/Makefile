include .env

include $(MAKE_DIR)/*.mk

.PHONY: all build install uninstall clean deep-clean

all: install build_project

build: $(PROGRAM_NAME)

install: build
	install $(PROGRAM_NAME) $(INSTALL_PATH)

uninstall:
	rm -f $(INSTALL_PATH)/$(PROGRAM_NAME)

clean: uninstall
	rm -rf $(PROGRAM_NAME) $(MARK_DIR)

deep_clean: clean
	rm -f go.mod go.sum
