build-osx: clean format
	go build
	tar cvzf lockit-automerge-osx.tgz lockit-automerge

build-linux: clean format
	env GOOS=linux go build
	tar cvzf lockit-automerge-linux.tgz lockit-automerge

format:
	go fmt

clean:
	go clean

.PHONY: format clean build-osx build-linux