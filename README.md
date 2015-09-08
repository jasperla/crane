# Crane

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

Currently Crane is built as a dynamic executable. This means that in
order for it to run the dependant libraries must be installed
too. These libraries must also be installed at build time:

- libgit2
- libssh2

When these have been installed, install Crane with:

	go get github.com/RedCoolBeans/crane

## Usage

    crane -package=dockerlint -repo=ssh://git@github.com:RedCoolBeans \
        -destination=/ -sshkey=/home/jasper/.ssh/id_rsa

The public key name is derived from `sshkey`; if the key requires a
password it can be passed with `-sshkey` though this is not
recommended for unattended use.

When Crane has installed all software, it can remove itself by passing
the `-clean` argument to the last invocation. Or simply with:

    crane -clean

## MANIFEST.yml

The `MANIFEST.yaml` file will be used by Crane to resolve and deploy
dependencies (if any). It will also be used to track versions and
serve as a project/repository description. The following fields are
supported:

- `name`: name of the software (REQUIRED)
- `maintainer`: maintainer name (REQUIRED)
- `email`: maintainer email (REQUIRED)
- `homepage`: project homepage
- `version`: project version (REQUIRED)
- `revision`: cargo revision (starts at _0_)
- `architecture`: (array) supported architectures. NB: This field
  is currently ignored and may require repository layout changes. By
  default `x86_64` will be assumed.
- `dependencies`: (array) this contains a hash with the following
  fields:
  - `name`: dependency name (must correspond with the `name` of it's
    manifest) (REQUIRED)
  - `repo`: repository in syntax as passed to crane with `-repo`
    (REQUIRED)
  - `branch`, `prefix`: branch/prefix to use for this dependency.
    By default the same branch/prefix as the dependant package will
    be used (which in turn defaults to `master` and `` respectively)

Unless otherwise noted, all fields are strings. A basic utility called
`crane-manifest` can be build with:

	go get github.com/RedCoolBeans/crane/crane-manifest

Then run:

	crane-manifest -file MANIFEST.yaml

in order to find missing fields; it currently lacks strict type-checking.

## ToDo

### Short term goals:

- Add 'customer' flag which can be re-used in combination with the manifest
- Allow building on CargOS nativly
- Fix static linking to create a truly standalone binary
- WANTLIB-like script to find out what libs something has linked
  against, in order to list those as dependencies

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
