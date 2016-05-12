// Crane manifest validator
package main

import (
	"flag"
	"log"
	"strings"

	"github.com/RedCoolBeans/crane/util/gpg"
	"github.com/RedCoolBeans/crane/util/logging"
	"github.com/RedCoolBeans/crane/util/manifest"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	file := flag.String("file", "MANIFEST.yaml", "Path to manifest")
	debug := flag.Bool("debug", false, "Enable debug output")

	strict := flag.Bool("strict", false, "Enable signature checking")
	pubkey := flag.String("pubkey", "pubkey.asc", "Path to GPG public key")
	signature := flag.String("sig", "MANIFEST.yaml.sig", "Path to Manifest signature")

	flag.Parse()
	m := manifest.ReadFile(*file)

	if *debug {
		spew.Dump(m)
	}

	if *strict {
		if ok, ids := gpg.Verify(*pubkey, *signature, *debug); ok {
			logging.PrInfoBegin("Signature for MANIFEST.yaml verified\n")
			logging.PrInfoEnd("Signed by: %s\n", strings.Join(ids, "\n\t"))
		} else {
			logging.PrError("INVALID signature for MANIFEST.yaml! Aborting.")
		}
	}

	if err := manifest.Validate(m); err != nil {
		log.Fatalln(err)
	}
}
