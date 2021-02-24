PROVIDER_VERSION := 1.0

all: build

clean:
	rm -rf bin/
	rm -rf out/

gofmt:
	gofmt -w .

build: gofmt scraper provider

.PHONY: scraper
scraper:
	go build -o ./bin/scraper github.com/benoittoulme/terraform-provider-itsi/scraper

test: gofmt
	go test -v -cover github.com/benoittoulme/terraform-provider-itsi/...

.PHONY: provider
provider:
	GOARCH=amd64 GOOS=darwin  go build -v -o ./out/darwin_amd64/terraform-provider-itsi_$(PROVIDER_VERSION)   github.com/benoittoulme/terraform-provider-itsi/provider
	GOARCH=amd64 GOOS=linux   go build -v -o ./out/darwin_linux/terraform-provider-itsi_$(PROVIDER_VERSION)   github.com/benoittoulme/terraform-provider-itsi/provider
	GOARCH=amd64 GOOS=windows go build -v -o ./out/darwin_windows/terraform-provider-itsi_$(PROVIDER_VERSION) github.com/benoittoulme/terraform-provider-itsi/provider

install: build
	# handle MacOS install for now, local install for testing purposes with the provider example:
	mkdir -p  ~/Library/Application\ Support/io.terraform/plugins/terraform.com/itsi/itsi/$(PROVIDER_VERSION)/darwin_amd64/
	cp ./out/darwin_amd64/terraform-provider-itsi_$(PROVIDER_VERSION)  ~/Library/Application\ Support/io.terraform/plugins/terraform.com/itsi/itsi/$(PROVIDER_VERSION)/darwin_amd64/terraform-provider-itsi
	# for linux:
	# mkdir -p ~/.terraform.d/plugins
	# cp ./bin/terraform-provider-itsi ~/.terraform.d/plugins
