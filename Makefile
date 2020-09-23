hello:
	echo "Start build"

# Build
build:
	go build cmd/main.go

# Run
run:
	go run cmd/main.go

all: hello build