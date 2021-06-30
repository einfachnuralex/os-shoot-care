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
var VolumeID string

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

var checkCMD = &cobra.Command{
	Use:   "check",
	Short: "Check shoot resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := CheckAll(); err != nil {
			return err
		}
		return nil
	},
}

var checkVolume = &cobra.Command{
	Use:   "volume-attachments",
	Short: "Checks volume for attaches on server object",
	RunE: func(cmd *cobra.Command, args []string) error {
		return GetVolumeAttachments()
	},
}

func init() {
	utils.CreateOSClients(&ga)
	getCmd.AddCommand(checkVolume)
	rootCmd.AddCommand(getCmd, checkCMD)

	checkCMD.Flags().StringVarP(&Name, "name", "n", "", "Name of shoot")
	checkCMD.Flags().StringVarP(&Project, "project", "p", "", "Project ID")
	checkCMD.MarkFlagRequired("name")
	checkCMD.MarkFlagRequired("project")

	getCmd.Flags().StringVarP(&Name, "name", "n", "", "Name of shoot")
	getCmd.Flags().StringVarP(&Project, "project", "p", "", "Project ID")
	getCmd.MarkFlagRequired("name")
	getCmd.MarkFlagRequired("project")

	checkVolume.Flags().StringVarP(&VolumeID, "volume", "v", "", "ID of volume")
	checkVolume.MarkFlagRequired("volume")
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

func CheckAll() error {
	fmt.Println(Name, Project)
	vms, err := ga.GetInstancesByName(Name, Project)
	if err != nil {
		fmt.Println("error net: ", err)
	}
	//fmt.Println(vms)
	for _, vm := range vms {
		lovol, err := ga.GetServerLostVolumes(vm.ID)
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}
		if len(lovol) > 0 {
			fmt.Println("Server: ", vm.Name, vm.ID)
			for _, vol := range lovol {
				fmt.Println(vol.Name, vol.ID)
			}
		}
	}

	//vms2, err2 := ga.GetInstancesByName(Name, Project)
	//if err2 != nil {
	//	return err2
	//}
	//fmt.Println(vms2)
	//for _, vm := range vms2 {
	//	fmt.Println(vm.Name, vm.Tags)
	//}

	return nil
}

func GetVolumeAttachments() error {
	volumeAttachment, serverList, err := ga.GetVolumeAttachmentsForVolume(VolumeID)
	if err != nil {
		fmt.Println("error net: ", err)
	}

	if len(volumeAttachment) == 0 {
		fmt.Println("Volume reports no attachments")
	} else {
		fmt.Println("Volume reports server:")
	}
	for _, attachment := range volumeAttachment {
		fmt.Println(attachment.ServerID)
	}

	fmt.Println("")

	if len(serverList) == 0 {
		fmt.Println("No server reports attachment")
	} else {
		fmt.Println("Servers report attachments:")
	}
	for _, server := range serverList {
		timesAttached := 0
		for _, attachedVolume := range server.AttachedVolumes {
			if attachedVolume.ID == VolumeID {
				timesAttached++
			}
		}

		fmt.Printf("%s (%d times)", server.ID, timesAttached)
	}

	return nil
}
