binary = \
	bin/linux-amd64-chofuku \
	bin/windows-amd64-chofuku.exe \
	bin/darwin-amd64-chofuku

all: $(binary)

bin/linux-amd64-chofuku: main.go chofuku/*
	GOOS=linux GOARCH=amd64 go build -o $@

bin/windows-amd64-chofuku.exe: main.go chofuku/*
	GOOS=windows GOARCH=amd64 go build -o $@

bin/darwin-amd64-chofuku: main.go chofuku/*
	GOOS=darwin GOARCH=amd64 go build -o $@

clean:
	rm $(binary)
