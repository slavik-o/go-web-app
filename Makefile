.PHONY: run fmt css

run:
	@templ generate && go run .

fmt:
	@go fmt ./...

css:
	@cd styles && pnpx tailwindcss -i app.css -o ../assets/styles.css --watch

.DEFAULT_GOAL := run
