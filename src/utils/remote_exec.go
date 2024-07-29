package utils

import (
	"fmt"
	"log"
	"os/exec"
)

func RemoteExec(sshCtx SSHContext, cmd string) error {
	var command string
	if sshCtx.IdentityFile == "" {
		command = fmt.Sprintf("ssh -i ~/DevOps/.keys/test-stage-key.pem -o StrictHostKeyChecking=no %s@%s \"%s\"", sshCtx.RemoteUser, sshCtx.RemoteHost, cmd)
	} else {
		command = fmt.Sprintf("ssh -i %s -o StrictHostKeyChecking=no %s@%s \"%s\"", sshCtx.IdentityFile, sshCtx.RemoteUser, sshCtx.RemoteHost, cmd)
	}

	err := exec.Command("sh", "-c", command).Run()
	if err != nil {
		log.Printf("command: %s\nresult: %v", command, err)
	}
	log.Printf("command: %s\nresult: success", command)

	return nil
}
