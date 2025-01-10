# lc-go

**lc-go** is a Go port of [gdm85/lc](https://github.com/gdm85/lc).
lc(1) is a command similar to ls(1) included with the [Coherent](https://en.wikipedia.org/wiki/Coherent_(operating_system)) Unix clone.

## Installation

```shell
go install dbohdan.com/lc-go@latest
```

You may wish to rename `~/go/bin/lc-go` to `~/go/bin/lc`.

## Manual

The original manual page, lightly edited and converted to use modern `-man` macros, is in the file [`doc/lc.1`](doc/lc.1).
On a modern Linux or BSD system, it may be viewed with `man ./doc/lc.1`.
The original, unaltered manual page source is in [`doc/lc.orig.troff`](doc/lc.orig.troff).

A scanned and OCRed page introducing lc(1) in an early chapter of the [1994 user manual](https://archive.org/details/CoherentMan/page/n48/mode/1up) may be viewed below.
The page is also preserved as an [image](lc-intro.png).

<details>
<summary><h3>From the Coherent user manual</i></h3></summary>

The command **lc** also lists file names, but it prints the files and directories separately, in columns across the screen.
For example, typing:

```
lc
```

gives something of the form:

```
Directories:
    backup newdirectory
Files:
    another doc1 doc2 file01 file02
    stuff
```

If you want to list files in a directory other than your own, name that directory as an argument to the command.
For example, **/bin** is a directory in the COHERENT system that contains commands.
Type:

```
lc /bin
```

and **lc** will print the contents of **/bin**.

Both **ls** and **lc** can take options.
An option is indicated by a hyphen '-'.
The option must appear before any other argument.
For example, to list only the files in the directory for user **carol**, leaving out any directories, use the **f** option with **lc**:

```
lc -f /usr/carol
```

Or, if you type the command:

```
lc -f
```

the COHERENT system prints all of the files in the current directory.
The following gives the commonly used options to the command `lc`:

- **-d** List directories only, omitting files
- **-f** List files only, omitting directories
- **-1** List files in single-column format
</details>

## Differences from the C version

- lc-go detects terminal width.
- lc-go does not derefence symlinks like gdm85/lc (commit `f6e696ef0e`).

## License

Released under the original [license](LICENSE.md).
