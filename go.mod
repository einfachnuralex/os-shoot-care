module github.com/einfachnuralex/os-shoot-care

go 1.16

require (
	github.com/gardener/gardener v1.25.1
	github.com/gophercloud/gophercloud v0.17.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stackitcloud/gophercloud-wrapper v0.0.0-20210701103103-2b49346ec4c5
	sigs.k8s.io/controller-runtime v0.9.2
)

replace (
	k8s.io/api => k8s.io/api v0.19.2
	k8s.io/client-go => k8s.io/client-go v0.19.2
)
