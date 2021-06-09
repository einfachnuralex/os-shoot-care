package utils

import (
	"github.com/einfachnuralex/os-shoot-care/shootcare"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"log"
)

func CreateOSClients(ga *shootcare.Gandalf) {
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatalf("get env: %v", err)
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	// Create network
	ga.Network, err = openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatalf("create provider: %v", err)
	}
	// create compute
	ga.Compute, err = openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatalf("create provider: %v", err)
	}
	// create storage
	ga.Storage, err = openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatalf("create provider: %v", err)
	}
}
