USERNAME = Software-Analysis-Team
PROGRAM = utbotgo
PROJECT_PATH = github.com/${USERNAME}/UTBotGo/${PROGRAM}
INSTALL_PATH = /usr/bin
MAKE_DIR = .make

.PHONY: is_installed is_not_installed build install run uninstall clean deep-clean

# +----------------------------+
# |  Targets for internal use  |
# +----------------------------+

is_installed = rm -f ${MAKE_DIR}/is_not_installed

is_not_installed = mkdir -p ${MAKE_DIR}; touch ${MAKE_DIR}/is_not_installed

go.mod:
	go mod init ${PROJECT_PATH}

go.sum: go.mod
	go mod tidy
	@if [ ! -f $@ ]; then touch $@; fi

${PROGRAM}: go.mod go.sum *.go
	go build .
	@${is_not_installed}

# Don't do any actions for other targets (used in target `run` for other
# command lines parameters)
%:

# +---------------------+
# |  Targets for users  |
# +---------------------+

all: build

build: ${PROGRAM}

install: build
	install ${PROGRAM} ${INSTALL_PATH}
	@${is_installed}

# Run program with command line parameters of `make` command without first
# parameter (which should be `run`). So, running of this target should be like
# `make run <parameters of program>`
run: build
	./${PROGRAM} $(word 2, $(MAKECMDGOALS))

uninstall: is_not_installed
	@if [ ! -f ${MAKE_DIR}/is_not_installed ]; then echo "uninstall ${PROGRAM} (it's not shell command)"; rm -f ${INSTALL_PATH}/${PROGRAM}; ${is_not_installed}; fi

clean: uninstall
	rm -f ${PROGRAM}

deep-clean: clean
	rm -rf go.mod go.sum ${MAKE_DIR}
