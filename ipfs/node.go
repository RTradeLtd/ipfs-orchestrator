package ipfs

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
)

const (
	keyNetworkID = "network_id"
	keyJobID     = "job_id"

	keyBootstrapPeers = "bootstrap_peers"
	keyDataDir        = "data_dir"

	keyPortSwarm   = "ports.swarm"
	keyPortAPI     = "ports.api"
	keyPortGateway = "ports.gateway"

	keyResourcesDisk   = "resources.disk"
	keyResourcesMemory = "resources.memory"
	keyResourcesCPUs   = "resources.cpus"
)

// NodeInfo defines metadata about an IPFS node
type NodeInfo struct {
	NetworkID string `json:"network_id"`
	JobID     string `json:"job_id"`

	Ports     NodePorts     `json:"ports"`
	Resources NodeResources `json:"resources"`

	// Metadata set by node client:
	// DockerID is the ID of the node's Docker container
	DockerID string `json:"docker_id"`
	// ContainerName is the name of the node's Docker container
	ContainerName string `json:"container_id"`
	// DataDir is the path to the directory holding all data relevant to this
	// IPFS node
	DataDir string `json:"data_dir"`
	// BootstrapPeers lists the peers this node was bootstrapped onto upon init
	BootstrapPeers []string `json:"bootstrap_peers"`
}

// NodePorts declares the exposed ports of an IPFS node
type NodePorts struct {
	Swarm   string `json:"swarm"`   // default: 4001
	API     string `json:"api"`     // default: 5001
	Gateway string `json:"gateway"` // default: 8080
}

// NodeResources declares resource quotas for this node
type NodeResources struct {
	DiskGB   int `json:"disk"`
	MemoryGB int `json:"memory"`
	CPUs     int `json:"cpus"`
}

func newNode(id, name string, attributes map[string]string) (NodeInfo, error) {
	// check if container is a node
	if !isNodeContainer(name) {
		return NodeInfo{DockerID: id, ContainerName: name}, fmt.Errorf("unknown name format %s", name)
	}

	// parse bootstrap state
	var peers []string
	json.Unmarshal([]byte(attributes[keyBootstrapPeers]), &peers)

	// parse resource data
	var (
		disk, _ = strconv.Atoi(attributes[keyResourcesDisk])
		mem, _  = strconv.Atoi(attributes[keyResourcesMemory])
		cpus, _ = strconv.Atoi(attributes[keyResourcesCPUs])
	)

	// create node metadata to return
	return NodeInfo{
		NetworkID: attributes[keyNetworkID],
		JobID:     attributes[keyJobID],

		Ports: NodePorts{
			Swarm:   attributes[keyPortSwarm],
			API:     attributes[keyPortAPI],
			Gateway: attributes[keyPortGateway],
		},
		Resources: NodeResources{
			DiskGB:   disk,
			MemoryGB: mem,
			CPUs:     cpus,
		},

		DockerID:       id,
		ContainerName:  name,
		DataDir:        attributes[keyDataDir],
		BootstrapPeers: peers,
	}, nil
}

func (n *NodeInfo) withDefaults() {
	if n.Resources.CPUs == 0 {
		n.Resources.CPUs = 4
	}
	if n.Resources.DiskGB == 0 {
		n.Resources.DiskGB = 100
	}
	if n.Resources.MemoryGB == 0 {
		n.Resources.MemoryGB = 4
	}

	// set container name from network name
	if n.ContainerName == "" {
		n.ContainerName = toNodeContainerName(n.NetworkID)
	}

	// set name for node operations - container name can be used interchangably
	// with container name
	if n.DockerID == "" {
		n.DockerID = n.ContainerName
	}
}

func (n *NodeInfo) labels(peers []string, dataDir string) map[string]string {
	var peerBytes, _ = json.Marshal(peers)
	return map[string]string{
		keyNetworkID: n.NetworkID,
		keyJobID:     n.JobID,

		keyPortSwarm:   n.Ports.Swarm,
		keyPortAPI:     n.Ports.API,
		keyPortGateway: n.Ports.Gateway,

		keyBootstrapPeers: string(peerBytes),
		keyDataDir:        dataDir,

		keyResourcesCPUs:   strconv.Itoa(n.Resources.CPUs),
		keyResourcesDisk:   strconv.Itoa(n.Resources.DiskGB),
		keyResourcesMemory: strconv.Itoa(n.Resources.MemoryGB),
	}
}

func (n *NodeInfo) updateFromContainerDetails(c *types.Container) {
	if c == nil {
		return
	}

	// check container ID
	n.DockerID = c.ID

	// check ports
	if len(c.Ports) > 0 {
		for _, p := range c.Ports {
			var public = strconv.Itoa(int(p.PublicPort))
			var private = strconv.Itoa(int(p.PrivatePort))
			switch private {
			case containerSwarmPort:
				n.Ports.Swarm = public
			case containerAPIPort:
				n.Ports.API = public
			case containerGatewayPort:
				n.Ports.Gateway = public
			}
		}
	}
}
