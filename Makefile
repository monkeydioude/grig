.PHONY: dev
dev:
	make -j 2 templ-watch gow-watch

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: gow-watch
gow-watch:
	gow run cmd/grig-server/*.go

.PHONY: install
install:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/mitranim/gow@latest