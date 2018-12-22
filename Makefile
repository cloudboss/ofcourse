export GOARCH=amd64

_output/go/bin/go-bindata:
	mkdir -p _output
	GOPATH=`pwd`/_output/go go get -u github.com/go-bindata/go-bindata
	GOPATH=`pwd`/_output/go go get -u github.com/go-bindata/go-bindata/...

ofcourse/bindata.go: _output/go/bin/go-bindata templates/resource/resource.go \
	templates/resource/resource_test.go templates/Dockerfile templates/Makefile \
	templates/pipeline.yml templates/go.mod templates/cmd/check/main.go \
	templates/cmd/out/main.go templates/cmd/in/main.go
	./_output/go/bin/go-bindata \
		-o ofcourse/bindata.go \
		-pkg ofcourse \
		-prefix templates templates/...

GOSRC = main.go cmd/root.go cmd/init.go ofcourse/ofcourse.go ofcourse/bindata.go

_output/darwin/ofcourse: $(GOSRC)
	mkdir -p _output/darwin
	GOOS=darwin go build -o _output/darwin/ofcourse .

_output/darwin/ofcourse_darwin_amd64.zip: _output/darwin/ofcourse
	cd _output/darwin && zip ofcourse_darwin_amd64.zip ofcourse

_output/linux/ofcourse: $(GOSRC)
	mkdir -p _output/linux
	GOOS=linux go build -o _output/linux/ofcourse .

_output/linux/ofcourse_linux_amd64.zip: _output/linux/ofcourse
	cd _output/linux && zip ofcourse_linux_amd64.zip ofcourse

_output/windows/ofcourse.exe: $(GOSRC)
	mkdir -p _output/windows
	GOOS=windows go build -o _output/windows/ofcourse.exe .

_output/windows/ofcourse_windows_amd64.zip: _output/windows/ofcourse.exe
	cd _output/windows && zip ofcourse_windows_amd64.zip ofcourse.exe

ofcourse: _output/darwin/ofcourse _output/linux/ofcourse _output/windows/ofcourse.exe

dist: _output/darwin/ofcourse_darwin_amd64.zip _output/linux/ofcourse_linux_amd64.zip \
	_output/windows/ofcourse_windows_amd64.zip

test:
	go test -v ./... -run .

clean:
	GOPATH=`pwd`/_output/go go clean -modcache
	rm -rf _output

fmt:
	find . -name '*.go' -not -path './templates/*' -not -path './_output/*' | while read -r f; do \
		gofmt -w -s "$$f"; \
	done

.DEFAULT_GOAL := ofcourse
.PHONY: test
