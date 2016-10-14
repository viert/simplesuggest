all: simplesuggest

simplesuggest: dependencies src/simplesuggest.go src/trie/trie.go src/trie/node.go src/config/config.go src/web/server.go
	GOPATH=$(CURDIR) /usr/local/go/bin/go build src/simplesuggest.go

test: test_dependencies src/trie/trie_test.go src/trie/trie.go src/trie/node.go
	$(MAKE) -C $(CURDIR)/src/trie

dependencies: src/github.com/op/go-logging src/github.com/viert/properties src/github.com/gorilla/mux

test_dependencies: src/github.com/nu7hatch/gouuid

src/github.com/op/go-logging:
	GOPATH=$(CURDIR) /usr/local/go/bin/go get github.com/op/go-logging

src/github.com/viert/properties:
	GOPATH=$(CURDIR) /usr/local/go/bin/go get github.com/viert/properties

src/github.com/gorilla/mux:
	GOPATH=$(CURDIR) /usr/local/go/bin/go get github.com/gorilla/mux

src/github.com/nu7hatch/gouuid:
	GOPATH=$(CURDIR) /usr/local/go/bin/go get github.com/nu7hatch/gouuid

clean:
	rm -f simplesuggest
	find $(CURDIR) -name '*.a' -delete
