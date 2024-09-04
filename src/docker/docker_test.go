package docker_test

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"io"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/onurcevik/deploy-utilities/src/docker"
)

// MockDockerClient is a mock implementation of the Docker client
type MockDockerClient struct {
	client.APIClient
	mock.Mock
}

func (m *MockDockerClient) RegistryLogin(ctx context.Context, auth registry.AuthConfig) (registry.AuthenticateOKBody, error) {
	args := m.Called(ctx, auth)
	return args.Get(0).(registry.AuthenticateOKBody), args.Error(1)
}

func (m *MockDockerClient) ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, ref, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockDockerClient) BuildCachePrune(ctx context.Context, opts types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*types.BuildCachePruneReport), args.Error(1)
}

func (m *MockDockerClient) ImagesPrune(ctx context.Context, filters filters.Args) (image.PruneReport, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(image.PruneReport), args.Error(1)
}

func (m *MockDockerClient) ContainersPrune(ctx context.Context, filters filters.Args) (container.PruneReport, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(container.PruneReport), args.Error(1)
}

func (m *MockDockerClient) NetworksPrune(ctx context.Context, filters filters.Args) (network.PruneReport, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(network.PruneReport), args.Error(1)
}

func (m *MockDockerClient) VolumesPrune(ctx context.Context, filters filters.Args) (volume.PruneReport, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(volume.PruneReport), args.Error(1)
}

func TestLoginDocker(t *testing.T) {
	mockClient := new(MockDockerClient)
	d := docker.Docker{
		Client: mockClient,
	}

	user := "testuser"
	password := "testpassword"
	registryUri := "https://testregistry123456.com"

	authConfig := registry.AuthConfig{
		Username:      user,
		Password:      password,
		ServerAddress: registryUri,
	}

	//expectedEncodedAuth := "encodedAuth"
	mockClient.On("RegistryLogin", mock.Anything, authConfig).Return(registry.AuthenticateOKBody{}, nil)

	err := d.LoginDocker(user, password, registryUri)
	require.NoError(t, err)
	//assert.Equal(t, expectedEncodedAuth, d.EncodedAuth)

	mockClient.AssertExpectations(t)
}

// TestPullDockerImage tests PullDockerImage method
func TestPullDockerImage(t *testing.T) {
	mockClient := new(MockDockerClient)
	d := docker.Docker{
		Client: mockClient,
	}

	imageRef := "testimage:latest"

	// Create a buffer with some dummy data to simulate Docker pull output
	dummyData := []byte("dummy image pull data")
	mockReadCloser := io.NopCloser(bytes.NewReader(dummyData))
	mockClient.On("ImagePull", mock.Anything, imageRef, mock.Anything).Return(mockReadCloser, nil)

	err := d.PullDockerImage(imageRef)
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestPruneAll(t *testing.T) {
	mockClient := new(MockDockerClient)
	ctx := context.Background()
	d := docker.Docker{
		Client: mockClient,
		Ctx:    ctx,
	}

	pruneFilters := filters.NewArgs()
	pruneFilters.Add("dangling", "false")

	mockClient.On("BuildCachePrune", ctx, mock.Anything).Return(&types.BuildCachePruneReport{SpaceReclaimed: 100}, nil)
	mockClient.On("ImagesPrune", ctx, pruneFilters).Return(image.PruneReport{SpaceReclaimed: 200}, nil)
	mockClient.On("ContainersPrune", ctx, pruneFilters).Return(container.PruneReport{SpaceReclaimed: 300}, nil)
	mockClient.On("NetworksPrune", ctx, pruneFilters).Return(network.PruneReport{}, nil)
	mockClient.On("VolumesPrune", ctx, pruneFilters).Return(volume.PruneReport{SpaceReclaimed: 400}, nil)

	spaceReclaimed, err := d.PruneAll()
	require.NoError(t, err)
	assert.Equal(t, uint64(1000), spaceReclaimed)

	mockClient.AssertExpectations(t)
}
