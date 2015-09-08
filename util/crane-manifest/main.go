// Crane manifest validator
package main

import (
	"flag"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/RedCoolBeans/crane/util/manifest"
)

func main() {
	file := flag.String("file", "MANIFEST.yaml", "Path to manifest")
	debug := flag.Bool("debug", false, "Enable debug output")

	flag.Parse()
	m := manifest.ReadFile(*file)

	if *debug {
		spew.Dump(m)
	}

	if err := manifest.Validate(m); err != nil {
		log.Fatalln(err)
	}
}
