# Crane

![Crane](crane.png)

Crane is a tool to deploy pre-built software or configuration directly
from Git onto the local filesystem of a Docker image. This allows for
installing only what's needed without any "recommended" dependencies
or other unneeded cruft.

### What it does
Essentially it does a Git clone and copy the files around, removing
any traces of the clones and itself when finished. This way no further
software can be installed as there's no package manager around anymore.

### What it does not
Crane does not build or configure anything. This allows it to be very
flexible and handle any file format. Build once, deploy everywhere.

## Install

Dependencies:

- libgit2 (0.24)
- libssh2
- [Glide](https://glide.sh/) (for updating dependencies)

When these have been installed, install Crane with:

	go get -u github.com/RedCoolBeans/crane/crane

### Static version

In order to be truly standalone a static version of crane can be built.
This requires that that libgit2 is built with `-DBUILD_SHARED_LIBS=Off`.

While CargOS ships with Go < 1.5 (i.e. no support for `vendor/`), you have
to retrieve the dependencies first. Easiest is to run the above `go get`
command and skip cloning the repository in the block below.

For example on CargOS:

	mkdir -p $GOPATH/src/github.com/RedCoolBeans
	cd $GOPATH/src/github.com/RedCoolBeans
	git clone https://github.com/RedCoolBeans/crane.git
	pkgin in libgit2-static
	cd crane
	bmake

This creates a `crane/crane.static`.

## Usage

    crane -package=dockerlint -repo=ssh://git@github.com:RedCoolBeans \
        -destination=/ -sshkey=/home/jasper/.ssh/id_rsa

### Repositories

Crane's building blocks are _packages_. These are Git repositories that
contain the manifest files. In order for Crane to download such a _package_,
it requires two parameters to specify the location: `-package` and `-repo`.
The way they work is that they're concatenated like `repo + / + package` to
get the actual URI. For example in the above it would result in:
`ssh://git@github.com:RedCoolBeans/dockerlint`.

Crane supports both HTTPS and SSH repositories. The `-repo` URI is to be
specified as follows:

- HTTPS: `https://SERVER/SUB/DIR`
- SSH: `ssh://USERNAME@SERVER:SUB/DIR`

### SSH keys

The public key name is derived from `sshkey`; if the key requires a
password it can be passed with `-sshpass` though this is not
recommended for unattended use, or security.

### Strict mode

By default Crane operates in _strict mode_ which means the following:

- the MANIFEST.yaml is required to be signed. I.e. the `MANIFEST.yaml.sig`
  has to exist the package repository.
- all files that are to be installed have to have their checksum recorded
  in the `contents` section of the manifest. A missing or incorrect checksum
  will result in termination. In non-strict mode a mismatch is fatal but the
  absence of a checksum is not.

Strict mode can be disabled with `-strict=false`

### "Self-destruct"

When Crane has installed all software, it can remove itself by passing
the `-clean` argument to the last invocation. Or simply with:

    crane -clean

## MANIFEST.yaml

The `MANIFEST.yaml` file will be used by Crane to resolve and deploy
dependencies (if any). It will also be used to track versions and
serve as a project/repository description. The following fields are
supported:

- `name`: (string) name of the software (REQUIRED)
- `maintainer`: (string) maintainer name (REQUIRED)
- `email`: (string) maintainer email (REQUIRED)
- `homepage`: (string) project homepage
- `version`: (string) project version (REQUIRED)
- `revision`: (string) cargo revision (starts at _0_)
- `architecture`: (array) supported architectures. NB: This field
  is currently ignored and may require repository layout changes. By
  default `x86_64` will be assumed.
- `dependencies`: (array) this contains a hash with the following
  fields:
  - `name`: (string) dependency name (must correspond with the `name` of it's
    manifest) (REQUIRED)
  - `repo`: (string) repository in syntax as passed to crane with `-repo`
    (REQUIRED)
  - `branch`, `prefix`: branch/prefix to use for this dependency.
    By default the same branch/prefix as the dependant package will
    be used (which in turn defaults to `master` and `` respectively)
- `contents`: (array) contains a hash with names of files that are to
   be installed. If not specified all files are installed verbatim.
   Valid fields are:
     - `path`: (string) path to install the file to
	 - `sha256`: (string) SHA256 sum of the file
	 - `mode`: (int) filemode, defaults to `0644`

Unless otherwise noted, all fields are strings. A basic utility called
`crane-manifest` can be build with:

	go get github.com/RedCoolBeans/crane/crane-manifest

Then run:

	crane-manifest -file MANIFEST.yaml

in order to find missing fields; it currently lacks strict type-checking.

### Example

An example manifest that would be located in the `dockerlint.git` repository:

```
---
  name: 'dockerlint'
  version: '1.0'
  maintainer: 'Jasper Lievisse Adriaanse'
  email: 'jasper@redcoolbeans.com'
  homepage: 'https://github.com/redcoolbeans/dockerlint'
  dependencies:
    - name: 'nodejs'
      repo: 'ssh://git@git.redcoolbeans.com:software/nodejs'
      branch: '{{ .Customer }}'
  contents:
    - path: README
      sha256: 52eba98ea2584afc1a03d92344181b09aa7ac7b9715d2b03942a88160a769bf3
      mode: 0644
    - path: script.sh
      sha256: 4bc94a8b8ad00708da1a0986fd507146bf01fe2b97c95cd5ecf4c5c3147b8779
      mode: 0755
```

The first fields (`name` through `homepage`) are metadata fields which describe the
package at hand. The `version` would generally need to be set to the upstream package
version, and the `revision` is to track local changes to the package itself.

The `dependencies` block lists all repositories on which the current package depends.
Note that the `branch` field in the above example contains [`text/template`](https://golang.org/pkg/text/template/)
syntax, however that is currently not yet supported and merely a future task. The
`branch` defaults to `master` and can be set to any arbitrary branch of the repository
as needed by the package at hand.

`contents` function as a packaging list, describing which files in this repository
are to be installed. The `path` field is concatenated to the `-destination` flag
of crane, or the `destination` field in the manifest (the latter takes precedence).
Thus assuming a `-destination=/usr/local`, then the `script.sh` would be installed
as `/usr/local/script.sh`.

If omited, all files (excluding the `MANIFEST.yaml`) will be installed verbatim
with the default modes corresponding to their type (see listing above).

## Signed MANIFEST files

Checksums for clandestinely modified files are still considered valid if
an the checksums in the manifest file are updated correspondingly. To combat this, the
manifest file can be GPG signed, and by default `crane` will verify the signature
before attempting to use a manifest.

Crane requires two additional files:

- The "detached" signature (default: `MANIFEST.yaml.sig`)
- Public key of the signer (default: `pubkey.asc`)

To sign a manifest: `gpg --armor --output MANIFEST.yaml.sig --detach-sig MANIFEST.yaml`

## ToDo

### Short term goals:

- Add 'customer' flag which can be re-used in combination with the manifest

### Long term goals (roadmap)

- multi-platform support, including FreeBSD and Windows

## License

Apache v2, please see the LICENSE file.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
