username = Software-Analysis-Team
program = UTBotGo

go.mod:
	go mod init github.com/${username}/${program}

${program}: go.mod *.go
	go build .

build: ${program}
	sudo ln -fi ${program} /usr/bin/${program}

run: ${program}
	./$^ $(filter-out $@, $(MAKECMDGOALS))

clean:
	rm -f ${program}

deep-clean: clean
	rm -f go.mod

%:
	@:
