package manifest

type Manifest struct {
	Name         string
	Maintainer   string
	Email        string
	Homepage     string
	Architecture []string
	Dependencies struct {
		Name   string
		Repo   string
		Branch string
	}
	Contents []struct {
		Path   string
		Sha256 string
		Mode   string
	}
}
