package cmd

import (
	"fmt"
	"github.com/einfachnuralex/os-shoot-care/shootcare"
	"github.com/einfachnuralex/os-shoot-care/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

var ga shootcare.Gandalf
var Name string
var Project string

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get shoot resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := getAll(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	utils.CreateOSClients(&ga)
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&Name, "name", "n", "", "Name of shoot")
	getCmd.Flags().StringVarP(&Project, "project", "p", "", "Project ID")
	getCmd.MarkFlagRequired("name")
	getCmd.MarkFlagRequired("project")
}

func getAll() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "%s \n", "Network")
	fmt.Fprintf(w, "%s \t %s \t %s \t \n", "ID", "Subnets", "Net Name")
	net, err := ga.GetNetworkByName(Name, Project)
	if err != nil {
		fmt.Println("error net: ", err)
	}
	if net.ID != "" {
		fmt.Fprintf(w, "%s \t %s \t %s \t \n", net.ID, net.Subnets, net.Name)
	}
	w.Flush()

	fmt.Fprintf(w, "%s \n", "Server")
	fmt.Fprintf(w, "%s \t %s \t %s \t\n", "ID", "Name", "Status")
	vms, err := ga.GetInstancesByNetwork(net.ID, Project)
	if err != nil {
		fmt.Println("error net: ", err)
	}
	for _, vm := range vms {
		fmt.Fprintf(w, "%s \t %s \t %s \t\n", vm.ID, vm.Name, vm.Status)
	}
	w.Flush()

	fmt.Fprintf(w, "%s \n", "Ports")
	fmt.Fprintf(w, "%s \t %s \t %s \t %s \t\n", "ID", "IP", "Owner", "Server ID")
	ports, err := ga.GetPortsByNetwork(net.ID, Project)
	if err != nil {
		fmt.Println("error port: ", err)
	}
	for _, port := range ports {
		if !(port.DeviceOwner == "network:ha_router_replicated_interface") &&
			!(port.DeviceOwner == "network:dhcp") {
			fmt.Fprintf(w, "%s \t %s \t %s \t %s \t\n", port.ID, port.FixedIPs[0].IPAddress, port.DeviceOwner, port.DeviceID)
		}
	}
	w.Flush()
	return nil
}
