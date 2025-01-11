# lc-go(1)

## Synopsis

```
lc-go [-1Fabcdfp] [DIRECTORY]...
```

## Description

`lc-go` lists the contents of the specified DIRECTORY, or the current directory if none is provided.
It categorizes files by type (e.g., regular files, directories) and displays them in columns for readability.
By default, filenames are printed in color based on their type.
Color output can be disabled by setting the environment variable `NO_COLOR=1`.

## Options

- `-1`: List one file per line, instead of in columns.
- `-F`: Enable file type indicators after filenames, like `*` after executable files.
- `-a`: List all files, including hidden entries `.` and `..`.
- `-b`: List block-special files only.
- `-c`: List character-special files only.
- `-d`: List directories only.
- `-f`: List regular files only.
- `-l`: List symbolic links only.
- `-p`: List pipe (FIFO) files only.

## Usage

Options can be combined for customized output.
For example, to display all files, one per line:

```
lc-go -a1
```

If no directory is specified, `lc` defaults to the current working directory.

## Examples

The following command lists all regular files in the current directory:

```
lc-go -f
```

## See Also

- lc(1)
- ls(1)
