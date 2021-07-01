package openstack

import (
	"github.com/einfachnuralex/os-shoot-care/shootcare"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"log"
)

func CreateClients(ga *shootcare.Gandalf) {
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
	//fmt.Println(ga.Compute.Microversion)
	////ga.Compute.Microversion = "2.60"
	//fmt.Println(ga.Compute.Microversion)
	//fmt.Println(ga.Compute.Endpoint)
	////ga.Compute.Endpoint = "https://platform.cloud.schwarz:8774/v2.6/"
	//fmt.Println(ga.Compute.Endpoint)
	// create storage
	ga.Storage, err = openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		log.Fatalf("create provider: %v", err)
	}
}

func CreateComputeClient() (*gophercloud.ServiceClient,error) {
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return nil, err
	}
	return openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
}