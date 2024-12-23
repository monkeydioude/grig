.PHONY: dev
dev:
	make -j 2 templ-watch gow-watch

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: gow-watch
gow-watch:
	gow run cmd/grig-server/*.go