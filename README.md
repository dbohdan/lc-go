# lc-go

**lc-go** is an enhanced Go port of [gdm85/lc](https://github.com/gdm85/lc).
`lc` is a command similar to `ls` that was originally included with the [Coherent](https://en.wikipedia.org/wiki/Coherent_(operating_system)) Unix clone.
`lc` lists files grouped by category: directory, file, symlink, etc.
Its output consists of columns after [Bill Joy's `ls` in 1BSD](https://tldp.org/LDP/LG/issue48/fischer.html).
Unlike `ls`, `lc` always lists dot files.

## Sample output

This is the output of running `lc-go` in a clone of this repository:

```none
$ go build
$ ln -s lc-go lc
$ ./lc -F
Directories:
    .git/        .github/     doc/

Files:
    .gitignore   LICENSE.md   Makefile     README.md    go.mod
    go.sum       lc-go*       main.go      main_test.go

Symlinks:
    lc@
```

## Installation

First, install Go 1.22 or later.
Once it is installed, run the following:

```shell
go install dbohdan.com/lc-go@latest

# or

git clone https://github.com/dbohdan/lc-go
cd lc-go
make install
```

You may wish to symlink `~/go/bin/lc-go` to `~/go/bin/lc` for a shorter command.
`make install` does this automatically.
To avoid name conflict with the original C version of `lc`, you can rename it to, for example, `lco`.

## Compatibility

lc-go supports FreeBSD, Linux, macOS, NetBSD, and OpenBSD.
It is automatically tested on those systems.
lc-go builds but doesn't work on Windows.

## Documentation

### Manual page

The file [`doc/lc-go.1`](doc/lc-go.1) contains an update of the original manual page,
It has been converted to use modern `-man` macros.
You will also find it in Markdown in [`doc/lc-go.1.md`](doc/lc-go.1.md).
On a modern Linux or BSD system, it may be viewed with `man ./doc/lc-go.1`.
The original, unaltered manual page source is in [`doc/lc.orig.troff`](doc/lc.orig.troff).

### Introduction

A scanned and OCRed page introducing `lc` in an early chapter of the [Coherent user manual](https://archive.org/details/CoherentMan/page/n48/mode/1up) may be viewed in [`doc/intro.md`](doc/intro.md).
The page is also preserved as an [image](doc/intro.png).

## Differences from the C version

lc-go:

- Detects terminal width.
- Adds a mode for filenames longer than the available terminal width.
  It prints only one column when they are present.
  (This is different from the flag `-1`, which disables columns.)
- Adds optional color.
  Color is enabled by default.
  It can be disabled using the environment variable [`NO_COLOR`](https://no-color.org/).
- Adds a new flag `-F`.
  `-F` enables file type indicators after the filename:
  - `/` for directories;
  - `*` for executable files;
  - `@` for symlinks;
  - `=` for sockets;
  - `|` for pipes (FIFOs).
- Does not dereference symlinks when detecting their type (a bugfix).

## License

Released under the original BSD-3-Clause [license](LICENSE.md).
