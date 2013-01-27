// © 2012 Steve McCoy. Available under the MIT license.

/*
The 9spell command runs the plan9port "9 spell" program on each
of the files supplied as arguments. It prints output like:
	file0:44+/teh/
	file0:63+/frgo/
	file1:0+/fner/
A program such as acme can read those addresses and navigate
to the misspelled word.

If a filename ends in ".tex", that file is piped through the plan9port
"9 detex" program before "9 spell".
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/mccoyst/pipeline"
)

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		os.Stderr.WriteString("I need the names of files to check.\n")
		os.Exit(1)
	}

	for _, file := range files {
		check(file)
	}
}

func check(file string) {
	typos := findTypos(file)
	if typos == nil {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %q: %v\n", file, err)
		return
	}
	defer f.Close()

	in := bufio.NewReader(f)
	n := 0
	for {
		line, err := in.ReadString('\n')
		if err == io.EOF && len(line) == 0 {
			break
		} else if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Problem reading %q: %v\n", file, err)
			break
		}

		for _, w := range strings.FieldsFunc(trim(line), isWordSep) {
			if typos[w] {
				fmt.Printf("%s:%d+/%s/\n", file, n, w)
			}
		}
		n++
	}
}

func findTypos(file string) map[string]bool {
	var cmds pipeline.P
	var err error

	if strings.HasSuffix(file, ".tex") {
		cmds, err = pipeline.New(
			exec.Command("9", "delatex", file),
			exec.Command("9", "spell"))
	} else {
		cmds, err = pipeline.New(exec.Command("9", "spell", file))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem creating commands: %v\n", err)
		return nil
	}

	o, err := cmds.Last().StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem piping 9 spell: %v\n", err)
		return nil
	}

	typos := map[string]bool{}
	out := bufio.NewReader(o)
	err = cmds.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem starting 9 spell: %v\n", err)
		return nil
	}
	for {
		typo, err := out.ReadString('\n')
		typos[trim(typo)] = true
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Problem reading 9 spell output: %v\n", err)
			break
		}
	}

	errs := cmds.Wait()
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "Problems running 9 spell: %v\n", errs)
		return nil
	}

	return typos
}

var trim = strings.TrimSpace

func isWordSep(r rune) bool {
	return r != '.' && !unicode.IsLetter(r)
}
