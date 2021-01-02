

cross-build-sesamefs:
	mkdir -p build 
	cd sesamefs && go mod vendor
	# Linux AMD64
	cd sesamefs && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../build/sesamefs
	tar -czf build/sesamefs-linux-amd64.tar.gz -C build/ sesamefs
	rm -f build/sesamefs
	sha256sum build/sesamefs-linux-amd64.tar.gz > build/sesamefs-linux-amd64.asc
	# Linux ARM
	cd sesamefs && env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ../build/sesamefs
	tar -czf build/sesamefs-linux-arm.tar.gz -C build sesamefs
	rm -f build/sesamefs
	sha256sum build/sesamefs-linux-arm.tar.gz > build/sesamefs-linux-arm.asc
	# Linux ARM64
	cd sesamefs && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ../build/sesamefs
	tar -czf build/sesamefs-linux-arm64.tar.gz -C build sesamefs
	rm -f build/sesamefs
	sha256sum build/sesamefs-linux-arm64.tar.gz > build/sesamefs-linux-arm64.asc
	# FreeBSD AMD64
	cd sesamefs && env GOOS=freebsd GOARCH=amd64 go build -o ../build/sesamefs
	tar -czf build/sesamefs-freebsd-amd64.tar.gz -C build sesamefs
	rm -f build/sesamefs
	sha256sum build/sesamefs-freebsd-amd64.tar.gz > build/sesamefs-freebsd-amd64.asc
	# FreeBSD ARM
	cd sesamefs && env GOOS=freebsd GOARCH=arm go build -o ../build/sesamefs
	tar -czf build/sesamefs-freebsd-arm.tar.gz -C build sesamefs
	rm -f build/sesamefs
	sha256sum build/sesamefs-freebsd-arm.tar.gz > build/sesamefs-freebsd-arm.asc
	# FreeBSD ARM64
	cd sesamefs && env GOOS=freebsd GOARCH=arm64 go build -o ../build/sesamefs
	tar -czf build/sesamefs-freebsd-arm64.tar.gz -C build sesamefs
	rm -f build/sesamefs
	sha256sum build/sesamefs-freebsd-arm64.tar.gz > build/sesamefs-freebsd-arm64.asc