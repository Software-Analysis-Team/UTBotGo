username = Software-Analysis-Team
program = UTBotGo

go.mod:
	go mod init github.com/${username}/${program}

${program}: go.mod
	go build .

build: ${program}

run: ${program}
	./$^

clean:
	rm -f ${program}

deep-clean: clean
	rm -f go.mod
