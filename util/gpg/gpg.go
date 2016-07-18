package gpg

import (
	"os"
	"path"

	"github.com/RedCoolBeans/crane/util/logging"
	"golang.org/x/crypto/openpgp"
)

const MANIFEST = "MANIFEST.yaml"

func Pubkey(path string) *os.File {
	pubkey, err := os.Open(path)
	if err != nil {
		logging.PrError("Could not open public key %s: %v", path, err)
	}

	return pubkey
}

func Signature(path string) *os.File {
	signature, err := os.Open(path)
	if err != nil {
		logging.PrError("Could not open signature file %s: %v", path, err)
	}

	return signature
}

func Verify(pubkeyPath string, signaturePath string, clonedir string, verbose bool) (bool, []string) {
	ids := make([]string, 0)

	pubkey := Pubkey(pubkeyPath)
	signature := Signature(path.Join(clonedir, signaturePath))

	manifest, err := os.Open(path.Join(clonedir, MANIFEST))
	if err != nil {
		logging.PrError("Could not open %s: %v", MANIFEST, err)
	}

	keyring, err := openpgp.ReadArmoredKeyRing(pubkey)
	if err != nil {
		logging.PrError(err.Error())
	}

	// Unless we're in debug mode, we don't care about the specifics of why the
	// signature didn't check out. Yes/no is all that matters then.
	signer, err := openpgp.CheckArmoredDetachedSignature(keyring, manifest, signature)
	if err != nil {
		if verbose {
			logging.PrError(err.Error())
		} else {
			return false, ids
		}
	}

	for id := range signer.Identities {
		ids = append(ids, id)
	}

	return true, ids
}
