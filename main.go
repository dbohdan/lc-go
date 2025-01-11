// lc lists files in categories (and columns).

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

const (
	gap     = 1 // Minimum gap between columns.
	indent1 = 4 // Indent for multiple directories.
	indent2 = 4 // Indent for files in a category.
)

var (
	oneflag  bool   // One per line.
	Fflag    bool   // Classify files.
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
	twidth   = 80 // Default line width if we can't detect the terminal.

	fname string

	wout io.Writer
	werr io.Writer

	rsColor = color.New()                                          // rs=0
	bdColor = color.New(color.FgYellow, color.BgBlack, color.Bold) // bd=40;33;01
	cdColor = color.New(color.FgYellow, color.BgBlack, color.Bold) // cd=40;33;01
	diColor = color.New(color.FgBlue, color.Bold)                  // di=01;34
	exColor = color.New(color.FgGreen, color.Bold)                 // ex=01;32
	lnColor = color.New(color.FgCyan, color.Bold)                  // ln=01;36
	orColor = color.New(color.FgRed, color.BgBlack, color.Bold)    // or=40;31;01
	soColor = color.New(color.FgMagenta, color.Bold)               // so=01;35
	piColor = color.New(color.FgYellow, color.BgBlack)             // pi=40;33
	suColor = color.New(color.FgWhite, color.BgRed)                // su=37;41
	sgColor = color.New(color.FgBlack, color.BgYellow)             // sg=30;43
	stColor = color.New(color.FgWhite, color.BgBlue)               // st=37;44
	owColor = color.New(color.FgBlue, color.BgGreen)               // ow=34;42
	twColor = color.New(color.FgBlack, color.BgGreen)              // tw=30;42

	rsIndicator = ""
	diIndicator = "/"
	exIndicator = "*"
	lnIndicator = "@"
	soIndicator = "="
	piIndicator = "|"
)

type Entry struct {
	name string

	color     color.Color
	indicator string
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
		fmt.Fprintf(wout, "%*s", indent1, "")
	}
	fmt.Fprintf(wout, format, a...)
	printed = true
}

func prindent_empty() {
	if ndir > 1 {
		fmt.Fprintf(wout, "%*s", indent1, "")
	}
	printed = true
}

func main() {
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		twidth = width
	}
	lwidth = twidth - indent2

	wout = os.Stdout
	werr = os.Stderr

	args := os.Args
	var lastErr error

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

			case 'F':
				Fflag = true

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
		if err := lc("."); err != nil {
			lastErr = err
		}
	} else {
		if len(args)-argIdx > 1 {
			lwidth -= indent1
		}
		ndir = len(args) - argIdx
		for i := argIdx; i < len(args); i++ {
			if err := lc(args[i]); err != nil {
				lastErr = err
			}
			w.Flush()
		}
	}

	if lastErr != nil {
		os.Exit(1)
	}
}

// lc processes a single name.
func lc(name string) error {
	sb, err := os.Stat(name)
	if err != nil {
		fmt.Fprintf(wout, "%s: not found\n", name)
		return fmt.Errorf("stat failed: %w", err)
	}

	mode := sb.Mode()
	typeStr := ""

	switch {

	case mode.IsDir():
		return lcdir(name)

	case mode.IsRegular():
		typeStr = "file"

	case mode&os.ModeSymlink != 0:
		typeStr = "symlink"

	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0:
		typeStr = "character special file"

	case mode&os.ModeDevice != 0:
		typeStr = "block special file"

	case mode&os.ModeNamedPipe != 0:
		typeStr = "FIFO/pipe"

	case mode&os.ModeSocket != 0:
		typeStr = "socket"

	default:
		fmt.Printf("%s: unknown file type\n", name)
		return fmt.Errorf("unknown file type")
	}

	nameColor, indicator := style(mode)

	nameColor.Fprint(wout, name)
	fmt.Fprintf(wout, "%s: %s\n", indicator, typeStr)
	return nil
}

// lcdir processes one directory.
func lcdir(dname string) error {
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
			fmt.Fprintln(wout)
		}

		info, err := os.Stat(dname)
		if err != nil {
			return fmt.Errorf("cannot stat directory: %w", err)
		}

		if Fflag {
			dname = strings.TrimSuffix(dname, diIndicator)
		}

		dColor, indicator := style(info.Mode())
		dColor.Fprint(wout, dname)
		fmt.Fprintf(wout, "%s:\n", indicator)
	}
	printed = false

	dir, err := os.Open(dname)
	if err != nil {
		fmt.Fprintf(wout, "Cannot open directory `%s`\n", dname)
		return fmt.Errorf("cannot open directory: %w", err)
	}
	defer dir.Close()

	// If aflag is set, manually print an entry for "." and "..".
	if aflag {
		// Current directory.
		info, err := os.Stat(dname)
		if err != nil {
			return fmt.Errorf("cannot stat current directory: %w", err)
		}
		doentry(dname, "", fs.FileInfoToDirEntry(info))

		// Parent directory.
		abs, err := filepath.Abs(dname)
		if err != nil {
			return fmt.Errorf("cannot resolve absolute path: %w", err)
		}
		parentDir := filepath.Dir(abs)
		info, err = os.Stat(parentDir)
		if err != nil {
			return fmt.Errorf("cannot stat current directory: %w", err)
		}

		parentEntry := fs.FileInfoToDirEntry(info)
		doentry(dname, "..", parentEntry)
	}

	// Read and process regular directory entries.
	entries, err := dir.ReadDir(-1)
	if err != nil {
		fmt.Fprintf(werr, "%s: directory read error: %s\n", dname, err)
		return fmt.Errorf("directory read error: %w", err)
	}

	for _, entry := range entries {
		doentry(dname, "", entry)
	}

	prnames()
	return nil
}

func style(mode fs.FileMode) (color.Color, string) {
	var clr *color.Color = rsColor
	ind := rsIndicator

	switch {

	case mode.IsDir():
		clr = diColor
		ind = diIndicator

		// Check for a world-writable directory.
		if mode&0o002 != 0 {
			if mode&os.ModeSticky == 0 {
				clr = owColor
			} else {
				clr = twColor
			}
		} else if mode&os.ModeSticky != 0 {
			clr = stColor
		}

	case mode.IsRegular():
		if mode&0o001|mode&0o010|mode&0o100 != 0 {
			clr = exColor
			ind = exIndicator
		}
		if mode&os.ModeSetuid != 0 {
			clr = suColor
		} else if mode&os.ModeSetgid != 0 {
			clr = sgColor
		}

	case mode&os.ModeSymlink != 0:
		// Check if the symlink is broken.
		if _, err := os.Stat(fname); err != nil {
			clr = orColor
		} else {
			clr = lnColor
		}
		ind = lnIndicator

	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice == 0:
		clr = bdColor

	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0:
		clr = cdColor

	case mode&os.ModeSocket != 0:
		clr = soColor
		ind = soIndicator

	case mode&os.ModeNamedPipe != 0:
		clr = piColor
		ind = piIndicator
	}

	if !Fflag {
		ind = rsIndicator
	}

	return *clr, ind
}

// doentry processes a single directory entry.
func doentry(dirname, nameoverride string, dp fs.DirEntry) error {
	width := 0

	name := dp.Name()
	if nameoverride != "" {
		name = nameoverride
	}

	fullPath := filepath.Join(dirname, name)
	fname = fullPath

	sb, err := dp.Info()
	if err != nil {
		prindent("%s: cannot stat\n", name)
		return fmt.Errorf("cannot stat %s", name)
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
		return fmt.Errorf("unknown file type: %s", name)
	}

	entryColor, indicator := style(mode)

	width = len(name) + len(indicator)
	if width > maxwidth {
		maxwidth = width
	}

	e := Entry{
		name: name,

		color:     entryColor,
		indicator: indicator,
	}
	addlist(list, e)
	return nil
}

// addlist adds an entry to the list in sorted order.
func addlist(lsp *[]Entry, e Entry) {
	*lsp = append(*lsp, e)
	slices.SortFunc(*lsp, func(a, b Entry) int {
		return strings.Compare(a.name, b.name)
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
		fmt.Fprintln(wout)
	}

	prindent("%s:\n", typeStr)

	npl := 1
	if !oneflag {
		npl = lwidth / (maxwidth + gap)
	}

	col := 0
	for i, e := range list {
		if !oneflag && col == 0 {
			fmt.Fprintf(wout, "%*s", indent2, "")
			prindent_empty()
		}

		e.color.Fprint(wout, e.name)
		fmt.Fprint(wout, e.indicator)

		if col+1 != npl && i != len(list)-1 && maxwidth < lwidth {
			padding := maxwidth + gap - len(e.name) - len(e.indicator)
			fmt.Fprintf(wout, "%*s", padding, "")
			col++
		} else {
			col = 0
			fmt.Fprintln(wout)
		}
	}
}

func usage() {
	fmt.Fprintf(werr, "Usage: lc [-afdcbpl] [-1F] [name ...]\n")
	os.Exit(2)
}
