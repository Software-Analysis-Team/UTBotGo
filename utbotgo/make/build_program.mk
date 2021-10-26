go.mod:
	go mod init $(PROJECT_PATH)

go.sum: go.mod
	go mod tidy
	@if [ ! -f $@ ]; then touch $@; fi

$(PROGRAM_NAME): go.mod go.sum *.go
	go build .
