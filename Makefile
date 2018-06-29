
tools:
	go get -u -v github.com/golang/dep/cmd/dep
	dep ensure

embed:
	go-bindata -o include/include.go --pkg include include/*yml include/keys/* include/templates/*

clean:
	rm -rf dist
	rm -rf vendor

test:
	go test -v ./...


