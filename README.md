# [runzip](https://webinstall.dev/runzip)

A command line utility for unpacking .rar files.

```text
USAGE
	runzip <archive.rar> [./dst/]

EXAMPLES
	runzip ./archive.rar                 # ./inner-dir/
	runzip ./archive.rar ./existing-dir/ # ./existing-dir/inner-dir/
	runzip ./archive.rar ./new-dir/      # ./new-dir/
```

For archives with a single file or folder, this will extract that to the given directory.

For archives with multiple files or folders, it will create a directory of the same name as the archive.

# Install

## macOS, Linux, BSD, \*nix

```sh
curl -sS https://webi.sh/runzip | sh
source ~/.config/envman/PATH.env
```

## Windows

```sh
curl.exe -sS https://webi.ms/runzip | powershell
```

## Go

```sh
go install github.com/therootcompany/runzip
```

## Manual

Download the distributable version for your OS and Architecture from
<https://github.com/therootcompany/runzip/releases>, unzip, and place in your
`PATH`.

```sh
# note: tar works for unzip on Windows
tar xvf runzip-*.*
mv runzip ~/bin/
```

note: you can use [pathman](https://webinstall.dev/pathman) if you need a
cross-platform, cross-shell PATH manager.

# Build

## Go

```sh
./scripts/go-build
```

Or

```sh
curl -sS https://webi.sh/go | sh
source ~/.config/envman/PATH.env
```

```sh
go build -o ./runzip ./runzip.go
```

## TinyGo

```sh
./scripts/tinygo-build
```

Or

```sh
curl -sS https://webi.sh/tinygo | sh
source ~/.config/envman/PATH.env
```

```sh
tinygo build -o ./runzip ./runzip.go
```
