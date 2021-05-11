binary := \
	bin/linux-amd64-chofuku \
	bin/windows-amd64-chofuku.exe \
	bin/darwin-amd64-chofuku

build_options := --tags "json1"

all: $(binary)

bin/linux-amd64-chofuku: main.go chofuku/*
	GOOS=linux GOARCH=amd64 go build $(build_options) -o $@

bin/windows-amd64-chofuku.exe: main.go chofuku/*
	GOOS=windows GOARCH=amd64 go build $(build_options) -o $@

bin/darwin-amd64-chofuku: main.go chofuku/*
	GOOS=darwin GOARCH=amd64 go build $(build_options) -o $@

clean:
	rm $(binary)
