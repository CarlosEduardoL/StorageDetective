set shell := ["bash", "-uc"]
set dotenv-load

# paths
bin_dir := "bin"
bin_name := "stex"
bin_path := bin_dir / bin_name

# default recipe: list available recipes
default:
    @just --list

# build the binary into bin/
build:
    @mkdir -p {{bin_dir}}
    go build -o {{bin_path}} ./cmd/stex
    @echo "built {{bin_path}}"

# remove the bin/ directory
clean:
    rm -rf {{bin_dir}}

# run the app, forwarding any args (e.g. `just run -i ~/projects`)
run *args: build
    {{bin_path}} {{args}}

# run all tests
test:
    go test ./...

# run all tests in a container
test-container:
    podman run --rm -v {{justfile_directory()}}:/stex:Z -w /stex golang:1.26 go test ./...

# run a specific test in a container, e.g. `just testone-container TestExplorer_RebuildPopulatesItems`
testone-container name:
    podman run --rm -v {{justfile_directory()}}:/stex:Z -w /stex golang:1.26 go test -run '{{name}}' ./...

# run a specific test, e.g. `just testone TestExplorer_RebuildPopulatesItems`
testone name:
    go test -run '{{name}}' ./...

# go vet
vet:
    go vet ./...

# go fmt
fmt:
    go fmt ./...

# tidy modules
tidy:
    go mod tidy

# full check: fmt + vet + test
check: fmt vet test
