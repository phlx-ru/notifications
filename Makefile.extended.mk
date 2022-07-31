.PHONY: go-env
# Add Go binaries to PATH making it accessible after `go install ...`
go-env:
	@export PATH="$PATH:$(go env GOPATH)/bin"

.PHONY: run
# Run the service
run:
	@go run ./cmd/notifications -conf=./configs

.PHONY: vendor
# Make ./vendor folder with dependencies
vendor:
	@go mod tidy && go mod vendor && go mod verify
