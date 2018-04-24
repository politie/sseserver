
VERSION?=dev

clean:
	rm -rvf vendor
	rm -rvf bin/*

test-unit:
	go test -race -cover -v .

test: test-unit

lint:
	gometalinter --install
	gometalinter ./ --checkstyle > report.xml

dev-dependencies:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	dep ensure -v

release: clean dev-dependencies test
	./release.sh $(VERSION)
