package utils_test

import (
	"github.com/onurcevik/deploy-utilities/src/utils"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// execCommand is a variable to allow mocking
var remoteexecCommand = exec.Command

// MockCmd is a struct to mock exec.Cmd
type MockCmd struct {
	mock.Mock
}

// Run is a mock implementation of exec.Cmd's Run method
func (m *MockCmd) Run() error {
	args := m.Called()
	return args.Error(0)
}

// MockExecCommand is the mocked version of exec.Command
func MockExecCommand(name string, arg ...string) *exec.Cmd {
	mockCmd := &MockCmd{}
	mockCmd.On("Run").Return(nil)

	// Create a dummy *exec.Cmd
	return &exec.Cmd{
		Path: name,
		Args: append([]string{name}, arg...),
		Process: &os.Process{
			Pid: 1234, // Dummy PID
		},
	}
}

func init() {
	remoteexecCommand = MockExecCommand
}

// TestRemoteExec tests the RemoteExec function
func TestRemoteExec(t *testing.T) {
	// Mock exec.Command
	defer func() { remoteexecCommand = exec.Command }()

	sshCtx := utils.SSHContext{
		RemoteUser:   "testuser",
		RemoteHost:   "testhost",
		IdentityFile: "",
	}

	err := utils.RemoteExec(sshCtx, "ls -l")
	assert.NoError(t, err)
}

// TestRemoteExec tests the RemoteExec function
func TestSCP(t *testing.T) {
	// Mock exec.Command
	defer func() { remoteexecCommand = exec.Command }()

	sshCtx := utils.SSHContext{
		RemoteUser:   "testuser",
		RemoteHost:   "testhost",
		IdentityFile: "",
	}

	err := utils.SCP(sshCtx, "", "", true)
	assert.NoError(t, err)
}
