# Global constants
USERNAME = Software-Analysis-Team
PROGRAM = UTBotGo
INSTALL_PATH = /usr/bin

# Phony targets
.PHONY: build install run uninstall clean deep-clean

# Create Go-Module
go.mod:
	go mod init github.com/${USERNAME}/${PROGRAM}

# Create main executable program `UTBotGo`
${PROGRAM}: go.mod *.go
	go build .

# Alias for ${PROGRAM}
build: ${PROGRAM}
	@:

# Install `UTBotGo`
install: ${PROGRAM}
	install ${PROGRAM} ${INSTALL_PATH}

# Run `UTBotGo` with command line parameters of `make` command without first
# parameter (which should be `run`). So, running of this target should be like
# `make run <parameters of UTBotGo program>`
run: ${PROGRAM}
	./${PROGRAM} $(word 2, $(MAKECMDGOALS))

# Uninstall `UTBotGo`
uninstall:
	rm -f ${INSTALL_PATH}/${PROGRAM}

# Remove `UTBotGo`
clean: uninstall
	rm -f ${PROGRAM}

# Remove `UTBotGo` and `go.mod`
deep-clean: clean
	rm -f go.mod

# Don't do any actions for other targets (used in target `run` for other
# command lines parameters)
%:
	@:
