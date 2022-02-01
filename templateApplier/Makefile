name = templateApplier

all: build

build: $(name)

go.mod:
	go mod init $(name)

$(name): go.mod *.go
	go build .
