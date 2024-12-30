.PHONY: dev-linux
dev-linux:
	docker compose up -d grig-server

.PHONY: dev
dev:
	make -j 3 templ-watch gow-watch tailwind

.PHONY: templ-watch
templ-watch:
	templ generate --watch

.PHONY: gow-watch
gow-watch:
	gow run cmd/grig-server/*.go

.PHONY: tailwind
tailwind:
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

.PHONY: install-tailwind
install-tailwind:
	@if [ "$(shell uname)" = "Darwin" ]; then \
		brew install tailwindcss; \
	elif [ "$(shell uname)" = "Linux" ]; then \
		curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64; \
		chmod +x tailwindcss-linux-x64; \
		sudo mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss; \
	elif [ "$(OS)" = "Windows_NT" ]; then \
		echo "Downloading TailwindCSS for Windows"; \
		echo "Please download manually from:"; \
		echo "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-windows-x64.exe"; \
		echo "Add the downloaded file to your PATH as 'tailwindcss'."; \
	else \
		echo "Unsupported operating system."; \
		exit 1; \
	fi

.PHONY: install
install: install-tailwind
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/mitranim/gow@latest
