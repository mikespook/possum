default: test

fmt: 
	go fmt github.com/mikespook/possum

coverage: fmt
	go test ./ -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

test: fmt 
	go vet github.com/mikespook/possum
	go test github.com/mikespook/possum
