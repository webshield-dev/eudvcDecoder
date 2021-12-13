REPO_NAME   := eudvcdecoder

install-code-check:
	@echo "** install 3rd party static code checkers **"

	go install honnef.co/go/tools/cmd/staticcheck@latest

	go install github.com/securego/gosec/v2/cmd/gosec@latest

code-check:
	@echo "** static code checks **"

	# scan code staticcheck.io using standard config
	staticcheck ./...

	#scan code with vet finds more issues than compiler
	go vet ./...

	#scan code for security issues, only interested in medium of higher
	gosec -quiet -fmt=json -confidence=medium -tests=true ./...


drone-test:
	@make install-code-check
	@make code-check
	@echo "** Drone Testing (short) **"
	go get -v -d -t ./...
	go clean -testcache
	go test -short -mod=readonly -covermode=atomic  ./...

test:
	@make code-check
	go get -v -d -t ./...
	go clean -testcache
	go test -mod=readonly -covermode=atomic  ./...


macbin:
	@echo "** make mac executable **"
	go build -o ./bin/decoder.mac
