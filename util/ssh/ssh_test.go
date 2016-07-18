package ssh

import (
	"os"
	"testing"

	"github.com/RedCoolBeans/crane/util/ssh"
)

// helper which can be defered to close and remove a file
func cleanup(fd *os.File, path string, t *testing.T) {
	fd.Close()
	if err := os.Remove(path); err != nil {
		t.Errorf(err.Error())
	}
}

func TestInit(t *testing.T) {
	var tests = []struct {
		sshkey          string
		sshkeycreate    bool
		repository      string
		cargo           string
		sshpubkey       string
		sshpubkeycreate bool
		sshrepo         string
		sshuser         string
		err             error
	}{
		{"/tmp/.crane-key", true, "ssh://git@github.com/RedCoolBeans/", "crane", "/tmp/.crane-key.pub", true, "git@github.com/RedCoolBeans/crane", "j", nil},
	}

	for i, tt := range tests {
		// If needed, create the keyfiles. If omited we'll test for errors being thrown.
		if tt.sshkeycreate {
			f, err := os.Create(tt.sshkey)
			if err != nil {
				t.Errorf(err.Error())
			}
			defer cleanup(f, tt.sshkey, t)
		}

		if tt.sshpubkeycreate {
			f, err := os.Create(tt.sshpubkey)
			if err != nil {
				t.Errorf(err.Error())
			}
			defer cleanup(f, tt.sshpubkey, t)
		}

		sshpubkey, sshrepo, sshuser, err := ssh.Init(tt.sshkey, tt.repository, tt.cargo)

		if err != nil {
			t.Errorf(err.Error())
		}

		if sshpubkey != tt.sshpubkey {
			t.Errorf("%d. %q => %q, wanted: %q (%q)", i, tt.sshkey, sshpubkey, tt.sshpubkey, err.Error())
		}

		if sshrepo != tt.sshrepo {
			t.Errorf("%d. %q, wanted sshrepo: %q", i, sshrepo, tt.sshrepo)
		}

		if sshuser != tt.sshuser {
			t.Errorf("%d. %q, wanted sshuser: %q", i, sshuser, tt.sshuser)
		}
	}
}

func TestPubKey(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{"", ".pub"},
		{"id_rsa", "id_rsa.pub"},
	}

	for i, tt := range tests {
		p := ssh.PubKey(tt.in)
		if p != tt.out {
			t.Errorf("%d. %q => %q, wanted: %q", i, tt.in, p, tt.out)
		}
	}
}

func TestRepository(t *testing.T) {
	var tests = []struct {
		repository string
		cargo      string
		out        string
	}{
		{"ssh://git@github.com/RedCoolBeans", "crane", "git@github.com/RedCoolBeans/crane"},
		{"ssh://git@github.com/RedCoolBeans/", "crane", "git@github.com/RedCoolBeans/crane"},
	}

	for i, tt := range tests {
		r := ssh.Repository(tt.repository, tt.cargo)
		if r != tt.out {
			t.Errorf("%d. %q/%q => %q, wanted: %q", i, tt.repository, tt.cargo, r, tt.out)
		}
	}
}

func TestFindUserName(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{"ssh://git@github.com/RedCoolBeans", "j"},
		{"git@github.com/RedCoolBeans", "j"},
		{"ssh://github.com/RedCoolBeans/", "git"},
		{"github.com/RedCoolBeans/", "git"},
	}

	for i, tt := range tests {
		u := ssh.FindUserName(tt.in)
		if u != tt.out {
			t.Errorf("%d. %q => %q, wanted: %q", i, tt.in, u, tt.out)
		}
	}
}
