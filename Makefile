
build:
	CGO_ENABLED=0 go build -o _out/kpxc-fdos-fuse ./cmd/kpxc_fdos_fuse

run: build
	./_out/dops

clean:
	go clean
	rm ./_out/*

package:
	go clean
	rm -rf ./_out/*

	GOARCH=386   GOOS=linux   CGO_ENABLED=0 go build -o _out/kpxc-fdos-fuse_linux-386-static                     ./cmd/kpxc_fdos_fuse  # Linux - 32 bit
	GOARCH=amd64 GOOS=linux   CGO_ENABLED=0 go build -o _out/kpxc-fdos-fuse_linux-amd64-static                   ./cmd/kpxc_fdos_fuse  # Linux - 64 bit
	GOARCH=arm64 GOOS=linux   CGO_ENABLED=0 go build -o _out/kpxc-fdos-fuse_linux-arm64-static                   ./cmd/kpxc_fdos_fuse  # Linux - ARM
	GOARCH=386   GOOS=linux                 go build -o _out/kpxc-fdos-fuse_linux-386                            ./cmd/kpxc_fdos_fuse  # Linux - 32 bit
	GOARCH=amd64 GOOS=linux                 go build -o _out/kpxc-fdos-fuse_linux-amd64                          ./cmd/kpxc_fdos_fuse  # Linux - 64 bit
	GOARCH=arm64 GOOS=linux                 go build -o _out/kpxc-fdos-fuse_linux-arm64                          ./cmd/kpxc_fdos_fuse  # Linux - ARM
	GOARCH=amd64 GOOS=darwin                go build -o _out/kpxc-fdos-fuse_macos-amd64                          ./cmd/kpxc_fdos_fuse  # macOS - 32 bit
	GOARCH=amd64 GOOS=darwin                go build -o _out/kpxc-fdos-fuse_macos-amd64                          ./cmd/kpxc_fdos_fuse  # macOS - 64 bit
	GOARCH=amd64 GOOS=openbsd               go build -o _out/kpxc-fdos-fuse_openbsd-amd64                        ./cmd/kpxc_fdos_fuse  # OpenBSD - 64 bit
	GOARCH=arm64 GOOS=openbsd               go build -o _out/kpxc-fdos-fuse_openbsd-arm64                        ./cmd/kpxc_fdos_fuse  # OpenBSD - ARM
	GOARCH=amd64 GOOS=freebsd               go build -o _out/kpxc-fdos-fuse_freebsd-amd64                        ./cmd/kpxc_fdos_fuse  # FreeBSD - 64 bit
	GOARCH=arm64 GOOS=freebsd               go build -o _out/kpxc-fdos-fuse_freebsd-arm64                        ./cmd/kpxc_fdos_fuse  # FreeBSD - ARM
