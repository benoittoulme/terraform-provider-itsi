all: build

clean:
	rm -rf bin/

gofmt:
	gofmt -w .

build: gofmt scraper itsi-provider

.PHONY: scraper
scraper:
	go build -o ./bin/scraper github.com/benoittoulme/terraform-provider-itsi/scraper

test: gofmt
	go test -v -cover github.com/benoittoulme/terraform-provider-itsi/...

itsi-provider:
	go build -o ./bin/terraform-provider-itsi github.com/benoittoulme/terraform-provider-itsi/provider
	chmod u+x ./bin/terraform-provider-itsi

install: build
	mkdir -p ~/.terraform.d/plugins
	cp ./bin/terraform-provider-itsi ~/.terraform.d/plugins