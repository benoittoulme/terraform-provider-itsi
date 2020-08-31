all: build

clean:
	rm -rf bin/

gofmt:
	gofmt -w .

build: gofmt scraper

.PHONY: scraper
scraper:
	go build -o ./bin/scraper github.com/benoittoulme/terraform-provider-itsi/scraper

test: gofmt
	go test -v -cover github.com/benoittoulme/terraform-provider-itsi/...
