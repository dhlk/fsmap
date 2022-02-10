package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dhlk/fsmap"
)

func mapfs(prefix, algorithm, key string, create bool) (path string, err error) {
	var f *fsmap.Fsmap
	if f, err = fsmap.New(prefix, algorithm); err != nil {
		return
	}

	return f.Lookup([]byte(key), create)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [option...] dir key\n", os.Args[0])
		flag.PrintDefaults()
	}
	lookup := flag.Bool("l", false, "lookup (do not create)")
	algorithm := flag.String("a", "SHA-512", "naming algorithm")
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	path, err := mapfs(args[0], *algorithm, args[1], !*lookup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fsmap: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(path)
}
