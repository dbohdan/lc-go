// lc lists files in categories (and columns).

package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/term"
)

const (
	GAP     = 1   // Minimum gap between columns.
	INDENT1 = 4   // Indent for multiple directories.
	INDENT2 = 4   // Indent for files in a category.
	NFNAME  = 400 // Maximum a filename can expand to.
)

var (
	twidth   = 79   // Default line width if we can't detect terminal.
	oneflag  bool   // One per line.
	aflag    bool   // Do all entries, including "." and "..".
	bflag    bool   // Do block special.
	cflag    bool   // Do character special.
	dflag    bool   // Do directories.
	fflag    bool   // Do regular files.
	lflag    bool   // Do symlinks.
	mflag    bool   // Do multiplexor files.
	pflag    bool   // Do pipes.
	allflag  = true // Do all types.
	ndir     int
	printed  = false // Set when we have printed something.
	maxwidth int     // Maximum width of a filename.
	lwidth   int
	sb       fs.FileInfo
	fname    string
)

type Entry struct {
	e_name string
}

var (
	files  []Entry
	links  []Entry
	dirs   []Entry
	blocks []Entry
	chars  []Entry
	pipes  []Entry
	mults  []Entry
)

func prindent(format string, a ...any) {
	if ndir > 1 {
		fmt.Printf("%*s", INDENT1, "")
	}
	fmt.Printf(format, a...)
	printed = true
}

func prindent_empty() {
	if ndir > 1 {
		fmt.Printf("%*s", INDENT1, "")
	}
	printed = true
}

func main() {
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		twidth = width
	}
	lwidth = twidth - INDENT2

	args := os.Args
	estat := 0

	argIdx := 1
	for argIdx < len(args) && strings.HasPrefix(args[argIdx], "-") {
		for _, c := range args[argIdx][1:] {
			switch c {

			case 'f':
				fflag = true
				allflag = false

			case 'd':
				dflag = true
				allflag = false

			case 'l':
				lflag = true
				allflag = false

			case 'm':
				mflag = true
				allflag = false

			case 'b':
				bflag = true
				allflag = false

			case 'c':
				cflag = true
				allflag = false

			case 'p':
				pflag = true
				allflag = false

			case '1':
				oneflag = true

			case 'a':
				aflag = true

			default:
				usage()
			}
		}
		argIdx++
	}

	if allflag {
		fflag = true
		dflag = true
		cflag = true
		bflag = true
		mflag = true
		lflag = true
		pflag = true
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	// Process directories.
	if argIdx >= len(args) {
		ndir = 1
		if lc(".") != 0 {
			estat = 1
		}
	} else {
		if len(args)-argIdx > 1 {
			lwidth -= INDENT1
		}
		ndir = len(args) - argIdx
		for i := argIdx; i < len(args); i++ {
			if lc(args[i]) != 0 {
				estat = 1
			}
			w.Flush()
		}
	}

	os.Exit(estat)
}

// lc does "lc" on a single name.
func lc(name string) int {
	var err error
	sb, err = os.Stat(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: not found\n", name)
		return 1
	}

	mode := sb.Mode()
	typ := ""

	switch {

	case mode.IsDir():
		return lcdir(name)

	case mode.IsRegular():
		typ = "file"

	case mode&os.ModeSymlink != 0:
		typ = "symlink"

	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0:
		typ = "character special file"

	case mode&os.ModeDevice != 0:
		typ = "block special file"

	case mode&os.ModeNamedPipe != 0:
		typ = "FIFO/pipe"

	case mode&os.ModeSocket != 0:
		typ = "socket"

	default:
		fmt.Printf("%s: unknown file type\n", name)
		return 1
	}

	fmt.Printf("%s: %s\n", name, typ)
	return 0
}

// lcdir processes one directory.
func lcdir(dname string) int {
	clearlist(&files)
	clearlist(&links)
	clearlist(&dirs)
	clearlist(&blocks)
	clearlist(&chars)
	clearlist(&pipes)
	clearlist(&mults)
	maxwidth = 0

	if ndir > 1 {
		if printed {
			fmt.Println()
		}
		fmt.Printf("%s:\n", dname)
	}
	printed = false

	dir, err := os.Open(dname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open directory `%s`\n", dname)
		return 1
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: directory read error: %s\n", dname, err)
		return 1
	}

	for _, entry := range entries {
		doentry(dname, entry)
	}

	prnames()
	return 0
}

// doentry processes a single directory entry.
func doentry(dirname string, dp fs.DirEntry) int {
	width := 0

	name := dp.Name()
	if !aflag && (name == "." || name == "..") {
		return 0
	}

	fullPath := filepath.Join(dirname, name)
	fname = fullPath

	width = len(name)
	if width > maxwidth {
		maxwidth = width
	}

	var err error
	sb, err = dp.Info()
	if err != nil {
		prindent("%s: cannot stat\n", name)
		return 1
	}

	var list *[]Entry
	mode := sb.Mode()
	switch {

	case mode.IsRegular():
		list = &files

	case mode&os.ModeSymlink != 0:
		list = &links

	case mode.IsDir():
		list = &dirs

	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice == 0:
		list = &blocks

	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0:
		list = &chars

	case mode&os.ModeSocket != 0:
		list = &mults

	case mode&os.ModeNamedPipe != 0:
		list = &pipes

	default:
		prindent("%s: unknown file type\n", name)
		return 1
	}

	e := Entry{e_name: name}
	addlist(list, e)
	return 0
}

// addlist adds an entry to the list in sorted order.
func addlist(lsp *[]Entry, e Entry) {
	*lsp = append(*lsp, e)
	slices.SortFunc(*lsp, func(a, b Entry) int {
		return strings.Compare(a.e_name, b.e_name)
	})
}

// clearlist clears a list.
func clearlist(lsp *[]Entry) {
	*lsp = []Entry{}
}

// prnames prints all the appropriate names.
func prnames() {
	if dflag {
		prtype(dirs, "Directories")
	}

	if fflag {
		prtype(files, "Files")
	}

	if lflag {
		prtype(links, "Symlinks")
	}

	if cflag {
		prtype(chars, "Character special files")
	}

	if bflag {
		prtype(blocks, "Block special files")
	}

	if pflag {
		prtype(pipes, "Pipes")
	}

	if mflag {
		prtype(mults, "Multiplexed files")
	}
}

// prtype prints one type of files.
func prtype(list []Entry, typeStr string) {
	if len(list) == 0 {
		return
	}

	if printed {
		fmt.Println()
	}

	if !oneflag {
		prindent("%s:\n", typeStr)
	}

	npl := 1
	if !oneflag {
		npl = lwidth / (maxwidth + GAP)
	}

	i := 0
	for j, e := range list {
		if !oneflag && i == 0 {
			fmt.Printf("%*s", INDENT2, "")
			prindent_empty()
		}

		name := e.e_name
		fmt.Print(name)

		if i+1 != npl && j != len(list)-1 {
			padding := maxwidth + GAP - len(name)
			fmt.Printf("%*s", padding, "")
			i++
		} else {
			i = 0
			fmt.Println()
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: lc [-afdcbpl] [-1] [name ...]")
	os.Exit(1)
}
