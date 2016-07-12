package ssh

// sshOptions is a per-repository struct to track ssh parameters
type SshOptions struct {
	Enabled   bool
	Sshkey    string
	Sshpass   string
	Sshpubkey string
	Sshrepo   string
	Sshuser   string
}
