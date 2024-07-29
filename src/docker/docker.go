package docker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"io"
	"net/http"
	"os"
)

// Docker struct is used to pass Context and EncodedAuth around client.APIClient interface is used to make mock testing easier
type Docker struct {
	Ctx         context.Context
	Client      client.APIClient
	EncodedAuth string
}

// ClientOption defines the type for functional options.
type ClientOption func(*client.Client) error

// WithHost sets the host for the Docker client.
func WithHost(host string) ClientOption {
	return func(c *client.Client) error {
		newClient, err := client.NewClientWithOpts(client.WithHost(host))
		if err != nil {
			return err
		}
		*c = *newClient
		return nil
	}
}

// WithTLS sets the TLS configuration for the Docker client.
func WithTLS(certFile, keyFile, caFile string) ClientOption {
	return func(c *client.Client) error {
		// Create the TLS configuration.
		tlsConfig, err := NewTLSConfig(certFile, keyFile, caFile)
		if err != nil {
			return err
		}

		// Create a custom http.Client with the TLS configuration.
		customTransport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		customHTTPClient := &http.Client{
			Transport: customTransport,
		}

		// Create a new Docker client with the custom http.Client.
		newClient, err := client.NewClientWithOpts(
			client.FromEnv,
			client.WithAPIVersionNegotiation(),
			client.WithHTTPClient(customHTTPClient),
		)
		if err != nil {
			return err
		}
		*c = *newClient
		return nil
	}
}

// NewClient creates a new Docker client with the given options.
func NewClient(options ...ClientOption) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create docker client handle: %w", err)
	}

	// Apply the options
	for _, option := range options {
		if err := option(cli); err != nil {
			return nil, fmt.Errorf("could not apply option: %w", err)
		}
	}

	return cli, nil
}

// NewTLS creates a TLS configuration for the Docker client.
func NewTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load client certificate and key: %w", err)
	}

	// Load CA cert
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate file: %w", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("could not append CA certificate to pool")
	}

	// Create TLS config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	return tlsConfig, nil
}

// LoginDocker logs in registry with given credentials and sets encodedauth to Docker object for further use in other functions
func (d *Docker) LoginDocker(user, password, registryUri string) error {
	auth := registry.AuthConfig{
		Username:      user,
		Password:      password,
		ServerAddress: registryUri,
	}
	_, err := d.Client.RegistryLogin(nil, auth)
	if err != nil {
		return err
	}
	ea, err := registry.EncodeAuthConfig(auth)
	if err != nil {
		return err
	}
	d.EncodedAuth = ea
	return nil
}

// PullDockerImage pulls a Docker image from a registry
func (d *Docker) PullDockerImage(imageRef string) error {
	out, err := d.Client.ImagePull(d.Ctx, imageRef, image.PullOptions{
		All:           false,
		RegistryAuth:  "",
		PrivilegeFunc: nil,
		Platform:      "",
	})
	if err != nil {
		return fmt.Errorf("error pulling Docker image: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return fmt.Errorf("error reading output: %w", err)
	}

	return nil
}

// PruneAll prunes all unused and dangling docker objects
func (d *Docker) PruneAll() (uint64, error) {
	var spaceReclaimed uint64
	pruneFilters := filters.NewArgs()
	pruneFilters.Add("dangling", "false")

	bc, err := d.Client.BuildCachePrune(d.Ctx, types.BuildCachePruneOptions{
		All:         true,
		KeepStorage: 0,
		Filters:     pruneFilters,
	})
	if err != nil {
		return spaceReclaimed, err
	}
	spaceReclaimed += bc.SpaceReclaimed

	im, err := d.Client.ImagesPrune(d.Ctx, pruneFilters)
	if err != nil {
		return spaceReclaimed, err
	}
	spaceReclaimed += im.SpaceReclaimed

	con, err := d.Client.ContainersPrune(d.Ctx, pruneFilters)
	if err != nil {
		return spaceReclaimed, err
	}
	spaceReclaimed += con.SpaceReclaimed

	_, err = d.Client.NetworksPrune(d.Ctx, pruneFilters)
	if err != nil {
		return spaceReclaimed, err
	}

	vol, err := d.Client.VolumesPrune(d.Ctx, pruneFilters)
	if err != nil {
		return spaceReclaimed, err
	}
	spaceReclaimed += vol.SpaceReclaimed

	return spaceReclaimed, nil
}
