# "Using the COHERENT system" &mdash; page 15

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
The following gives the commonly used options to the command **lc**:

- **-d** List directories only, omitting files
- **-f** List files only, omitting directories
- **-1** List files in single-column format
