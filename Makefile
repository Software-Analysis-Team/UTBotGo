run:
	go run *.go

test:
	for test_file in test-files/*_test.go; do echo "$${test_file%_*}.go:"; go test $${test_file%_*}*; done

clean:
	rm test-files/*_test.go
