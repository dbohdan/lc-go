# lc(1)

## Synopsis

```
lc [-1abcdfp] [DIRECTORY]...
```

## Description

`lc` lists the contents of the specified DIRECTORY, or the current directory if none is provided.
It categorizes files by type (e.g. regular files, directories) and displays them in columns for readability.

## Options

- `-1`: List one file per line, instead of in columns.
- `-a`: List all files, including hidden entries '.' and '..'.
- `-b`: List block-special files only.
- `-c`: List character-special files only.
- `-d`: List directories only.
- `-f`: List regular files only.
- `-l`: List symbolic links only.
- `-p`: List pipe (FIFO) files only.

## Usage

Options can be combined for customized output. For example, to display all files, one per line:

```
lc -a1
```

If no directory is specified, `lc` defaults to the current working directory.

## Examples

The following command lists all regular files in the current directory:

```
lc -f
```

To copy all regular files to a directory called `mydir` you can use:

```
cp `lc -f` mydir
```

## See Also

- ls(1)
