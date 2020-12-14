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
	# handle MacOS install for now, local install for testing purposes with the provider example:
	cp ./bin/terraform-provider-itsi  ~/Library/Application\ Support/io.terraform/plugins/terraform.com/itsi/itsi/1.0/darwin_amd64/
	# for linux:
	# mkdir -p ~/.terraform.d/plugins
	# cp ./bin/terraform-provider-itsi ~/.terraform.d/plugins
