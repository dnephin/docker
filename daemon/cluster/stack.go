package cluster

import (
	"fmt"

	"github.com/docker/docker/pkg/namesgenerator"
)

// CreateStack(name, bundle string) (string, error) // TODO: add config

func (c *Cluster) CreateStack(name, bundleRef string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if !c.isActiveManager() {
		return "", c.errNoManager()
	}

	if name == "" {
		name = namesgenerator.GetRandomName(0)
	}

	if bundleRef == "" {
		return "", fmt.Errorf("bundle name cannot be empty")
	}

	_, err := c.config.Backend.ResolveBundleManifest(bundleRef)
	if err != nil {
		return "", err
	}

	return "", fmt.Errorf("not implemented")
}
