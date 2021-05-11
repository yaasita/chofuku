package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yaasita/chofuku/chofuku"
)

const version = "1.1.0"

func main() {
	size_only := flag.Bool("size-only", false, "check size only")
	head100k_only := flag.Bool("100k-only", false, "check only the first 100 kbytes")
	show_version := flag.Bool("version", false, "show version")
	flag.Parse()
	target_dir := flag.Arg(0)
	if target_dir == "" {
		fmt.Println(os.Args[0], "[options] /path/to/directory")
		fmt.Println("options:")
		flag.PrintDefaults()
		return
	}
	if *show_version {
		fmt.Println(version)
		return
	}

	c, err := chofuku.New(target_dir)
	check(err)
	defer finish(c)
	if *size_only {
		return
	}

	err = c.UpdateHead100k()
	check(err)
	if *head100k_only {
		return
	}

	err = c.UpdateFullHash()
	check(err)
}
func finish(c chofuku.Chofuku) {
	duplicates, err := c.GetDuplicates()
	check(err)
	b, err := json.MarshalIndent(duplicates, "", "  ")
	check(err)
	fmt.Println(string(b))
	c.Close()
}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
