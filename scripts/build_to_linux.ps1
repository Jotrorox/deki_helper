# Set the target Operating System and Architecture
$env:GOOS = "linux"
$env:GOARCH = "amd64"

# Build the Go project
go build .