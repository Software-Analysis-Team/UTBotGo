go.mod:
	go mod init $(PROJECT_PATH)

# `go.sum` will be created if `go.mod` has `require` section. So, if `go.mod`
# hasn't, this rule will be always applied. Moreover, rule for
# `$(PROGRAM_NAME)` will be always applied because this rule won't generate
# file `go.sum`. Solution of this problem is `touch`ing file `go.sum` if it
# wasn't be generated.
go.sum: go.mod
	go mod tidy
	@if [ ! -f $@ ]; then touch $@; fi

$(PROGRAM_NAME): go.mod go.sum *.go
	go build .
