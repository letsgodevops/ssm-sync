VERSION=0.1.0

build-docker:
	docker-compose build --pull

build:
	GORELEASER_CURRENT_TAG=$(VERSION) goreleaser release --skip-validate --rm-dist --skip-publish

release:
	GORELEASER_CURRENT_TAG=$(VERSION) goreleaser release --skip-validate --rm-dist

clean:
	rm -rf ./dist
