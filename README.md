# lc-go

**lc-go** is a Go port of [gdm85/lc](https://github.com/gdm85/lc).
lc(1) is a command similar to ls(1) included with the [Coherent](https://en.wikipedia.org/wiki/Coherent_(operating_system)) Unix clone.

## Installation

```shell
go install dbohdan.com/lc-go@latest
```

You may wish to rename `~/go/bin/lc-go` to `~/go/bin/lc`.

## Compatibility

lc-go supports FreeBSD, Linux, macOS, NetBSD, and OpenBSD.
It is automatically tested on those systems.
It builds but doesn't work on Windows.

## Documentation

### Manual page

The original manual page, lightly edited and converted to use modern `-man` macros, is in the file [`doc/lc.1`](doc/lc.1).
You will find it in Markdown in [`doc/lc.1.md`](doc/lc.1.md).
On a modern Linux or BSD system, it may be viewed with `man ./doc/lc.1`.
The original, unaltered manual page source is in [`doc/lc.orig.troff`](doc/lc.orig.troff).

### Introduction

A scanned and OCRed page introducing lc(1) in an early chapter of the [user manual](https://archive.org/details/CoherentMan/page/n48/mode/1up) may be viewed in [`doc/intro.md`](doc/intro.md).
The page is also preserved as an [image](doc/intro.png).

## Differences from the C version

- lc-go detects terminal width.
- lc-go does not dereference symlinks.
- lc-go handles very long filenames by printing only one column when they are present.

## License

Released under the original [license](LICENSE.md).
