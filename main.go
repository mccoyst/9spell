// © 2012 Steve McCoy. Available under the MIT license.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"index/suffixarray"
)

func main() {
	if len(os.Args) <= 1 {
		os.Stderr.WriteString("I need the names of files to check.\n")
		os.Exit(1)
	}

	for _, file := range os.Args[1:] {
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
	n := 1
	for {
		line, err := in.ReadSlice('\n')
		if err == io.EOF && len(line) == 0 {
			break
		} else if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Problem reading %q: %v\n", file, err)
			break
		}

		//BUG(mccoyst): This finds typos within word boundaries. E.g. "bufio" matches
		// "bufio" and "io".
		index := suffixarray.New(line)
		for typo := range typos {
			if len(index.Lookup([]byte(typo), 1)) == 1 {
				fmt.Printf("%s:%d: %s\n", file, n, typo)
			}
		}
		n++
	}
}

func findTypos(file string) map[string]bool {
	spell := exec.Command("9", "spell", file)
	o, err := spell.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem piping 9 spell: %v\n", err)
		return nil
	}

	typos := map[string]bool{}
	out := bufio.NewReader(o)
	err = spell.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem starting 9 spell: %v\n", err)
		return nil
	}
	for {
		typo, err := out.ReadString('\n')
		if err == io.EOF {
			typos[strings.TrimSpace(typo)] = true
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Problem reading 9 spell output: %v\n", err)
			break
		}
		typos[strings.TrimSpace(typo)] = true
	}

	err = spell.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem running 9 spell: %v\n", err)
		return nil
	}

	return typos
}
