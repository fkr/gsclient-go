package main

import (
	"bufio"
	"os"

	"github.com/gridscale/gsclient-go"
	log "github.com/sirupsen/logrus"
)

const LocationUuid = "45ed677b-3702-4b36-be2a-a2eab9827950"
const webServerFirewallTemplateUuid = "82aa235b-61ba-48ca-8f47-7060a0435de7"

type ServiceType string

const (
	Server   ServiceType = "server"
	Storage  ServiceType = "storage"
	Network  ServiceType = "network"
	IP       ServiceType = "ip"
	ISOImage ServiceType = "isoimage"
)

//enhancedClient inherits all methods from gsclient.Client
//We need this to implement a new additional method
type enhancedClient struct {
	*gsclient.Client
}

func main() {
	uuid := os.Getenv("GRIDSCALE_UUID")
	token := os.Getenv("GRIDSCALE_TOKEN")
	config := gsclient.NewConfiguration(
		"https://api.gridscale.io",
		uuid,
		token,
		true,
	)
	client := enhancedClient{
		gsclient.NewClient(config),
	}
	log.Info("gridscale client configured")

	log.Info("Create server: Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	serverCreateRequest := gsclient.ServerCreateRequest{
		Name:         "go-client-server",
		Memory:       1,
		Cores:        1,
		LocationUuid: LocationUuid,
	}
	cServer, err := client.CreateServer(serverCreateRequest)
	if err != nil {
		log.Fatal("Create server has failed with error", err)
	}
	log.WithFields(log.Fields{
		"server_uuid": cServer.ObjectUuid,
	}).Info("Server successfully created")
	defer client.deleteService(Server, cServer.ObjectUuid)

	//get a server to interact with
	server, err := client.GetServer(cServer.ObjectUuid)
	if err != nil {
		log.Error("Get server has failed with error", err)
		return
	}

	log.Info("Start server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Turn on server
	err = client.StartServer(server.Properties.ObjectUuid)
	if err != nil {
		log.Error("Start server has failed with error", err)
		return
	}
	log.Info("Server successfully started")

	log.Info("Stop server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//Turn off server
	err = client.StopServer(server.Properties.ObjectUuid)
	if err != nil {
		log.Error("Stop server has failed with error", err)
		return
	}
	log.Info("Server successfully stop")

	log.Info("Update server: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	err = client.UpdateServer(server.Properties.ObjectUuid, gsclient.ServerUpdateRequest{
		Name:   "updated server",
		Memory: 1,
	})
	if err != nil {
		log.Error("Update server has failed with error", err)
		return
	}
	log.Info("Server successfully updated")

	//Get events of server
	events, err := client.GetServerEventList(server.Properties.ObjectUuid)
	if err != nil {
		log.Error("Get events has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"events": events,
	}).Info("Events successfully retrieved")

	//Create storage, network, IP, and ISO-image to attach to the server
	log.Info("Create storage, Network, IP, ISO-image: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cStorage, err := client.CreateStorage(gsclient.StorageCreateRequest{
		Capacity:     1,
		LocationUuid: LocationUuid,
		Name:         "go-client-storage",
	})
	if err != nil {
		log.Error("Create storage has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"storage_uuid": cStorage.ObjectUuid,
	}).Info("Storage successfully created")
	defer client.deleteService(Storage, cStorage.ObjectUuid)

	cNetwork, err := client.CreateNetwork(gsclient.NetworkCreateRequest{
		Name:         "go-client-network",
		LocationUuid: LocationUuid,
	})
	if err != nil {
		log.Error("Create network has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"network_uuid": cNetwork.ObjectUuid,
	}).Info("Network successfully created")
	defer client.deleteService(Network, cNetwork.ObjectUuid)

	cIp, err := client.CreateIp(gsclient.IpCreateRequest{
		Name:         "go-client-ip",
		Family:       4,
		LocationUuid: LocationUuid,
	})
	if err != nil {
		log.Error("Create IP has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"IP_uuid": cIp.ObjectUuid,
	}).Info("IP successfully created")
	defer client.deleteService(IP, cIp.ObjectUuid)

	cISOimage, err := client.CreateISOImage(gsclient.ISOImageCreateRequest{
		Name:         "go-client-iso",
		SourceUrl:    "http://tinycorelinux.net/10.x/x86/release/TinyCore-current.iso",
		LocationUuid: LocationUuid,
	})
	if err != nil {
		log.Error("Create ISO-image has failed with error", err)
		return
	}
	log.WithFields(log.Fields{
		"isoimage_uuid": cISOimage.ObjectUuid,
	}).Info("ISO-image successfully created")
	defer client.deleteService(ISOImage, cISOimage.ObjectUuid)

	//Attach storage, network, IP, and ISO-image to a server
	err = client.LinkStorage(server.Properties.ObjectUuid, cStorage.ObjectUuid, false)
	if err != nil {
		log.Error("Link storage has failed with error", err)
		return
	}
	log.Info("Storage successfully attached")
	defer client.unlinkService(Storage, server.Properties.ObjectUuid, cStorage.ObjectUuid)

	err = client.LinkNetwork(
		server.Properties.ObjectUuid,
		cNetwork.ObjectUuid,
		webServerFirewallTemplateUuid,
		false,
		1,
		nil,
		gsclient.FirewallRules{},
	)
	if err != nil {
		log.Error("Link network has failed with error", err)
		return
	}
	log.Info("Network successfully linked")
	defer client.unlinkService(Network, server.Properties.ObjectUuid, cNetwork.ObjectUuid)

	err = client.LinkIp(server.Properties.ObjectUuid, cIp.ObjectUuid)
	if err != nil {
		log.Error("Link IP has failed with error", err)
		return
	}
	log.Info("IP successfully linked")
	defer client.unlinkService(IP, server.Properties.ObjectUuid, cIp.ObjectUuid)

	err = client.LinkIsoImage(server.Properties.ObjectUuid, cISOimage.ObjectUuid)
	if err != nil {
		log.Error("Link ISO-image has failed with error", err)
		return
	}
	log.Info("ISO-image successfully linked")
	defer client.unlinkService(ISOImage, server.Properties.ObjectUuid, cISOimage.ObjectUuid)

	log.Info("Unlink and delete: press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *enhancedClient) deleteService(serviceType ServiceType, id string) {
	switch serviceType {
	case Server:
		//turn off server before deleting
		err := c.StopServer(id)
		if err != nil {
			log.Error("Stop server has failed with error", err)
			return
		}
		err = c.DeleteServer(id)
		if err != nil {
			log.Error("Delete server has failed with error", err)
			return
		}
		log.Info("Server successfully deleted")
	case Storage:
		err := c.DeleteStorage(id)
		if err != nil {
			log.Error("Delete storage has failed with error", err)
			return
		}
		log.Info("Storage successfully deleted")
	case Network:
		err := c.DeleteNetwork(id)
		if err != nil {
			log.Error("Delete network has failed with error", err)
			return
		}
		log.Info("Network successfully deleted")
	case IP:
		err := c.DeleteIp(id)
		if err != nil {
			log.Error("Delete IP has failed with error", err)
			return
		}
		log.Info("IP successfully deleted")
	case ISOImage:
		err := c.DeleteISOImage(id)
		if err != nil {
			log.Error("Delete ISO-image has failed with error", err)
			return
		}
		log.Info("ISO-image successfully deleted")
	default:
		log.Error("Unknown service type")
		return
	}
}

func (c *enhancedClient) unlinkService(serviceType ServiceType, serverId, serviceId string) {
	switch serviceType {
	case Storage:
		err := c.UnlinkStorage(serverId, serviceId)
		if err != nil {
			log.Error("Unlink storage has failed with error", err)
			return
		}
		log.Info("Storage successfully unlinked")
	case Network:
		err := c.UnlinkNetwork(serverId, serviceId)
		if err != nil {
			log.Error("Unlink network has failed with error", err)
			return
		}
		log.Info("Network successfully unlinked")
	case IP:
		err := c.UnlinkIp(serverId, serviceId)
		if err != nil {
			log.Error("Unlink IP has failed with error", err)
			return
		}
		log.Info("IP successfully unlinked")
	case ISOImage:
		err := c.UnlinkIsoImage(serverId, serviceId)
		if err != nil {
			log.Error("Unlink ISO-image has failed with error", err)
			return
		}
		log.Info("ISO-image successfully unlinked")
	default:
		log.Error("Unknown service type")
		return
	}
}