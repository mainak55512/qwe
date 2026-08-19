.SILENT: run
run:
	go run .

.SILENT: build
build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./build/qwe .

.SILENT: release
release:
	@read -p "Enter git tag name (e.g., v1.0.0): " TAG; \
	if [ -z "$$TAG" ]; then \
		echo "Error: Tag name cannot be empty."; \
		exit 1; \
	fi; \
	echo "1. Creating and pushing tag $$TAG..."; \
	git tag $$TAG && \
	git push origin $$TAG && \
	echo "2. Building binaries..." && \
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./release/x64_linux/qwe . && \
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./release/x64_windows/qwe.exe . && \
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./release/x64_mac/qwe . && \
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ./release/arm64_mac/qwe . && \
	echo "3. Packaging archives..." && \
	zip -r ./release/x64_linux.zip ./release/x64_linux/ && \
	zip -r ./release/x64_windows.zip ./release/x64_windows/ && \
	zip -r ./release/x64_mac.zip ./release/x64_mac/ && \
	zip -r ./release/arm64_mac.zip ./release/arm64_mac/ && \
	rm -r ./release/x64_linux/ ./release/x64_windows/ ./release/x64_mac/ ./release/arm64_mac/ && \
	echo "4. Creating GitHub release and uploading zips..." && \
	gh release create $$TAG ./release/*.zip --generate-notes && \
	rm -rf ./release/ && \
	echo "Release $$TAG published successfully!"
