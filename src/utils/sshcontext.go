package utils

type SSHContext struct {
	RemoteUser   string
	RemoteHost   string
	IdentityFile string
}

func NewSSHContext(remoteUser, remoteHost, identityFile string) *SSHContext {
	return &SSHContext{
		RemoteUser:   remoteUser,
		RemoteHost:   remoteHost,
		IdentityFile: identityFile,
	}
}
