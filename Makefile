
.PHONY: build
build: build-frontend build-bin

.PHONY: build-bin
build-bin:
	mkdir -p bin
	CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -o bin/cirrus main.go

.PHONY: build-frontend
build-frontend:
	cd frontend && yarn build