all: build

build:
	go build -o terraform-provider-dalpha

clean:
	rm -f terraform-provider-dalpha