package utils

import (
	"fmt"
	"log"
	"os/exec"
)

func SCP(sshCtx SSHContext, fromPath, toPath string, errorIgnore bool) error {
	var command string
	if sshCtx.IdentityFile == "" {
		command = fmt.Sprintf("scp -r %s %s@%s:%s", fromPath, sshCtx.RemoteUser, sshCtx.RemoteHost, toPath)
	} else {
		command = fmt.Sprintf("scp -i %s -r %s %s@%s:%s", sshCtx.IdentityFile, fromPath, sshCtx.RemoteUser, sshCtx.RemoteHost, toPath)
	}

	err := exec.Command("sh", "-c", command).Run()
	if err != nil {
		log.Printf("file transportation; from=%s to=%s\nresult: %v", fromPath, toPath, err)
		if !errorIgnore {
			return fmt.Errorf("file transportation failed with error: %v", err)
		}
	} else {
		log.Printf("file transportation; from=%s to=%s\nresult: success", fromPath, toPath)
	}

	return nil
}
