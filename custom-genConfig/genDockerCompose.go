package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"strconv"
	"strings"
)

type DockerCompose struct {
	Version  string              `yaml:"version,omitempty"`
	Networks map[string]*Network `yaml:"networks,omitempty"`
	Services map[string]*Service `yaml:"services,omitempty"`
}

type Network struct {
	External *External `yaml:"external,omitempty"`
}

type External struct {
	Name string `yaml:"name,omitempty"`
}

type Service struct {
	Deploy      *Deploy             `yaml:"deploy,omitempty"`
	Hostname    string              `yaml:"hostname,omitempty"`
	Image       string              `yaml:"image,omitempty"`
	Ports       []string            `yaml:"ports,omitempty"`
	Networks    map[string]*ServNet `yaml:"networks,omitempty"`
	Environment []string            `yaml:"environment,omitempty"`
	WorkingDir  string              `yaml:"working_dir,omitempty"`
	Command     string              `yaml:"command,omitempty"`
	Volumes     []string            `yaml:"volumes,omitempty"`
}

type ServNet struct {
	Aliases []string `yaml:"aliases,omitempty"`
}

// Placement will be added future
type Deploy struct {
	Replicas      int            `yaml:"replicas,omitempty"`
	Placement     *Placement     `yaml:"placement,omitempty"`
	RestartPolicy *RestartPolicy `yaml:"restart_policy,omitempty"`
}

type Placement struct {
	Constraint []string `yaml:"constraints,omitempty"`
}

type RestartPolicy struct {
	Condition   string        `yaml:"condition,omitempty"`
	Delay       time.Duration `yaml:"delay,omitempty"`
	MaxAttempts int           `yaml:"max_attempts,omitempty"`
	Window      time.Duration `yaml:"window,omitempty"`
}

//var TAG = `:x86_64-1.0.0-beta`
// var TAG = `:x86_64-1.0.0`
var TAG = `:x86_64-1.1.0-alpha`

func GenDockerCompose(serviceName string, domainName string, orgName string, networkName string, nodeId string, num ...int) (*DockerCompose, error) {
	var dockerCompose = &DockerCompose{}
	dockerCompose.Version = "3"

	err := GenNetwork(dockerCompose, networkName)
	check(err)

	switch serviceName {
	case "peer", "couchdb":
		err = GenService(dockerCompose, domainName, orgName, serviceName, networkName, nodeId, num[0], num[1])
	default:
		err = GenService(dockerCompose, domainName, orgName, serviceName, networkName, nodeId, num[0])
	}

	return dockerCompose, nil
}

func GenDeploy(service *Service, nodeId string) error {
	var constraint []string
	constraint = append(constraint, "node.id == "+nodeId)
	deploy := &Deploy{
		Placement: &Placement{
			Constraint: constraint,
		},
		Replicas: 1,
		RestartPolicy: &RestartPolicy{
			Condition:   "on-failure",
			Delay:       5 * time.Second,
			MaxAttempts: 3,
		},
	}
	service.Deploy = deploy

	return nil
}

func GenService(dockerCompose *DockerCompose, domainName string, orgName string, serviceName string, networkName string, nodeId string, num ...int) error {
	var total int
	if len(num) > 1 {
		total = num[0] * num[1]
	} else {
		total = num[0]
	}

	dockerCompose.Services = make(map[string]*Service, total)

	for i := 0; i < total; i++ {
		var serviceHost string
		var service *Service
		// var couchPortInt = rand.Intn(150)
		// if couchPortInt == 59 {
		// 	couchPortInt = rand.Intn(150)
		// }
		// couchPort := strconv.Itoa(couchPortInt) + "84:5984"
		couchPort := generatePort("couch", "5984")
		// if couchPort == "5984:5984" {
		// 	couchPort := strconv.Itoa(rand.Intn(150)) + "84:5984"
		// 	fmt.Println(couchPort)
		// }
		// var caPortInt = rand.Intn(150)
		// if caPortInt == 70 {
		// 	caPortInt = rand.Intn(150)
		// }
		// caPort := strconv.Itoa(caPortInt) + "54:7054"
		caPort := generatePort("ca", "7054")
		// if caPort == "7054:7054" {
		// 	caPort := strconv.Itoa(rand.Intn(150)) + "54:7054"
		// 	fmt.Println(caPort)
		// }

		// var otherPort string
		// var peerPortInt = rand.Intn(150)
		// if peerPortInt == 70 {
		// 	peerPortInt = rand.Intn(150)
		// }
		// fmt.Println("couchPort ==> ", couchPort)
		// fmt.Println("caPort ==> ", caPort)

		// peerPortStr = otherPort
		peerPort1 := generatePort("peer", "7051")
		// peerPort1 := strconv.Itoa(peerPortInt) + "51:7051"
		portALl := strings.Split(peerPort1, ":")
		portStart := portALl[0]
		portPrefix := ""
		if len(portStart) > 4 {
			portPrefix = string(portStart[0:3])
		} else {
			portPrefix = string(portStart[0:2])
		}
		portEnd := "7053"
		// peerPort2 := strconv.Itoa(peerPortInt) + "53:7053"
		peerPort2 := portPrefix + "53:" + portEnd
		// fmt.Println("peerPort1 ==> ", peerPort1)
		// fmt.Println("peerPort2 ==> ", peerPort2)

		switch serviceName {
		case "zookeeper":
			serviceHost = "zookeeper" + strconv.Itoa(i)
			service = &Service{
				Hostname: serviceHost,
			}
			service.Networks = make(map[string]*ServNet, 1)
			service.Networks[networkName] = &ServNet{
				Aliases: []string{serviceHost + "." + domainName},
			}
			service.Image = "hyperledger/fabric-zookeeper" + TAG
			var zookeeperArray []string
			for j := 0; j < total; j++ {
				zookeeperArray = append(zookeeperArray, "server."+strconv.Itoa(j+1)+"=zookeeper"+strconv.Itoa(j)+":2888:3888")
			}
			zookeeperList := arrayToString(zookeeperArray, " ")
			service.Environment = make([]string, 3)
			service.Environment[0] = "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=" + networkName
			service.Environment[1] = "ZOO_MY_ID=" + strconv.Itoa(i+1)
			service.Environment[2] = "ZOO_SERVERS=" + zookeeperList
			err := GenDeploy(service, nodeId)
			check(err)

		case "kafka":
			serviceHost = "kafka" + strconv.Itoa(i)
			service = &Service{
				Hostname: serviceHost + "." + domainName,
			}
			service.Networks = make(map[string]*ServNet, 1)
			service.Networks[networkName] = &ServNet{
				Aliases: []string{serviceHost + "." + domainName},
			}
			service.Image = "hyperledger/fabric-kafka" + TAG
			var zookeeperArray []string
			for j := 0; j < 3; j++ { // 3 is number of zookeeper nodes
				zookeeperArray = append(zookeeperArray, "zookeeper"+strconv.Itoa(j)+":2181")
			}
			zookeeperString := arrayToString(zookeeperArray, ",")
			service.Environment = make([]string, 8)
			service.Environment[0] = "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=" + networkName
			service.Environment[1] = "KAFKA_MESSAGE_MAX_BYTES=103809024"       // 99 MB
			service.Environment[2] = "KAFKA_REPLICA_FETCH_MAX_BYTES=103809024" // 99 MB
			service.Environment[3] = "KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false"
			service.Environment[4] = "KAFKA_DEFAULT_REPLICATION_FACTOR=3"
			service.Environment[5] = "KAFKA_MIN_INSYNC_REPLICAS=2"
			service.Environment[6] = "KAFKA_ZOOKEEPER_CONNECT=" + zookeeperString
			service.Environment[7] = "KAFKA_BROKER_ID=" + strconv.Itoa(i)
			err := GenDeploy(service, nodeId)
			check(err)

		case "orderer":
			serviceHost = "orderer" + strconv.Itoa(i)
			service = &Service{
				Hostname: serviceHost + "." + domainName,
			}
			service.Networks = make(map[string]*ServNet, 1)
			service.Networks[networkName] = &ServNet{
				Aliases: []string{serviceHost + "." + domainName},
			}
			service.Image = "hyperledger/fabric-orderer" + TAG
			service.Environment = make([]string, 14)
			service.Environment[0] = "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=" + networkName
			service.Environment[1] = "ORDERER_GENERAL_LOGLEVEL=debug"
			service.Environment[2] = "ORDERER_GENERAL_LISTENADDRESS=0.0.0.0"
			service.Environment[3] = "ORDERER_GENERAL_GENESISMETHOD=file"
			service.Environment[4] = "ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block"
			service.Environment[5] = "ORDERER_GENERAL_LOCALMSPID=OrdererMSP"
			service.Environment[6] = "ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp"
			service.Environment[7] = "ORDERER_GENERAL_TLS_ENABLED=true"
			service.Environment[8] = "ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key"
			service.Environment[9] = "ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt"
			service.Environment[10] = "ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]"
			service.Environment[11] = "ORDERER_KAFKA_RETRY_SHORTINTERVAL=1s"
			service.Environment[12] = "ORDERER_KAFAK_RETRY_SHORTTOTAL=30s"
			service.Environment[13] = "ORDERER_KAFKA_VERBOSE=true"

			service.WorkingDir = "/opt/gopath/src/github.com/hyperledger/fabric"
			service.Command = "orderer"

			service.Volumes = make([]string, 3)
			service.Volumes[0] = "./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block"
			service.Volumes[1] = "./crypto-config/ordererOrganizations/" + domainName + "/orderers/" + serviceHost + "." + domainName + "/msp:/var/hyperledger/orderer/msp"
			service.Volumes[2] = "./crypto-config/ordererOrganizations/" + domainName + "/orderers/" + serviceHost + "." + domainName + "/tls/:/var/hyperledger/orderer/tls"
			err := GenDeploy(service, nodeId)
			check(err)

		case "ca":
			serviceHost = "ca" + orgName + strconv.Itoa(i)
			service = &Service{
				Hostname: serviceHost + "." + domainName,
			}
			orgId := strconv.Itoa(i + 1)
			service.Networks = make(map[string]*ServNet, 1)
			service.Networks[networkName] = &ServNet{
				Aliases: []string{serviceName + "_peer_" + orgName + strconv.Itoa(i)},
			}
			service.Image = "hyperledger/fabric-ca" + TAG
			service.Ports = make([]string, 1)
			service.Ports[0] = caPort
			service.Environment = make([]string, 5)
			service.Environment[0] = "FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server"
			service.Environment[1] = "FABRIC_CA_SERVER_CA_NAME=ca-" + orgName + orgId
			service.Environment[2] = "FABRIC_CA_SERVER_TLS_ENABLED=true"
			service.Environment[3] = "FABRIC_CA_SERVER_TLS_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca." + orgName + orgId + "." + domainName + "-cert.pem"
			service.Environment[4] = "FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/CA" + orgName + orgId + "_PRIVATE_KEY"
			service.Command = "sh -c 'fabric-ca-server start --ca.certfile /etc/hyperledger/fabric-ca-server-config/ca." + orgName + orgId + "." + domainName + "-cert.pem --ca.keyfile /etc/hyperledger/fabric-ca-server-config/CA" + orgId + "_PRIVATE_KEY -b admin:adminpw -d'"
			service.Volumes = make([]string, 1)
			service.Volumes[0] = "./" + orgName + "-artifacts/crypto-config/peerOrganizations/" + orgName + orgId + "." + domainName + "/ca/:/etc/hyperledger/fabric-ca-server-config"
			err := GenDeploy(service, nodeId)
			check(err)

		case "couchdb":
			serviceHost = serviceName + orgName + strconv.Itoa(i)
			service = &Service{
				Hostname: serviceHost + "." + domainName,
			}
			service.Image = "hyperledger/fabric-couchdb"
			service.Ports = make([]string, 1)
			service.Ports[0] = couchPort
			service.Networks = make(map[string]*ServNet, 1)
			service.Networks[networkName] = &ServNet{
				Aliases: []string{serviceHost},
			}
			err := GenDeploy(service, nodeId)
			check(err)

		case "peer":
			peerNum := strconv.Itoa(i % num[0])
			orgNum := strconv.Itoa((i / num[0]) + 1)
			// fmt.Println(orgNum)
			serviceHost = "peer" + peerNum + "_" + orgName + orgNum
			hostName := "peer" + peerNum + "." + orgName + orgNum + "." + domainName
			service = &Service{
				Hostname: hostName,
			}
			service.Image = "hyperledger/fabric-peer" + TAG
			service.Ports = make([]string, 2)
			service.Ports[0] = peerPort1
			service.Ports[1] = peerPort2
			service.Networks = make(map[string]*ServNet, 1)
			service.Networks[networkName] = &ServNet{
				Aliases: []string{hostName},
			}
			service.Environment = make([]string, 17)
			service.Environment[0] = "CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock"
			service.Environment[1] = "CORE_LOGGING_LEVEL=DEBUG"
			service.Environment[2] = "CORE_PEER_TLS_ENABLED=true"
			service.Environment[3] = "CORE_PEER_GOSSIP_USELEADERELECTION=true"
			service.Environment[4] = "CORE_PEER_GOSSIP_ORGLEADER=false"
			service.Environment[5] = "CORE_PEER_PROFILE_ENABLED=true"
			service.Environment[6] = "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt"
			service.Environment[7] = "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key"
			service.Environment[8] = "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt"
			service.Environment[9] = "CORE_PEER_ID=" + hostName
			service.Environment[10] = "CORE_PEER_ADDRESS=" + hostName + ":7051"
			service.Environment[11] = "CORE_PEER_GOSSIP_EXTERNALENDPOINT=" + hostName + ":7051"
			service.Environment[12] = "CORE_PEER_LOCALMSPID=" + orgName + orgNum + "MSP"
			service.Environment[13] = "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=" + networkName
			service.Environment[14] = "CORE_LEDGER_STATE_STATEDATABASE=CouchDB"
			service.Environment[15] = "CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb" + orgName + strconv.Itoa(i) + ":5984"
			service.Environment[16] = "CORE_PEER_GOSSIP_BOOTSTRAP=peer0." + orgName + orgNum + "." + domainName + ":7051"
			//service.Environment[3]  = "CORE_PEER_ENDORSER_ENABLED=true"
			//service.Environment[6]  = "CORE_PEER_GOSSIP_SKIPHANDSHAKE=true"
			service.WorkingDir = "/opt/gopath/src/github.com/hyperledger/fabric/peer"
			service.Command = "peer node start"
			service.Volumes = make([]string, 3)
			service.Volumes[0] = "/var/run/:/host/var/run/"
			service.Volumes[1] = "./" + orgName + "-artifacts/crypto-config/peerOrganizations/" + orgName + orgNum + "." + domainName + "/peers/" + hostName + "/msp:/etc/hyperledger/fabric/msp"
			service.Volumes[2] = "./" + orgName + "-artifacts/crypto-config/peerOrganizations/" + orgName + orgNum + "." + domainName + "/peers/" + hostName + "/tls:/etc/hyperledger/fabric/tls"
			err := GenDeploy(service, nodeId)
			check(err)

		case "cli":
			serviceHost = orgName + "cli"
			service = &Service{}
			service.Image = "hyperledger/fabric-tools" + TAG
			service.Networks = make(map[string]*ServNet, 1)
			service.Networks[networkName] = &ServNet{
				Aliases: []string{orgName + "cli"},
			}
			service.Environment = make([]string, 12)
			service.Environment[0] = "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=" + networkName
			service.Environment[1] = "GOPATH=/opt/gopath"
			service.Environment[2] = "CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock"
			service.Environment[3] = "CORE_LOGGING_LEVEL=DEBUG"
			service.Environment[4] = "CORE_PEER_ID=cli"
			service.Environment[5] = "CORE_PEER_ADDRESS=peer0." + orgName + "1." + domainName + ":7051"
			service.Environment[6] = "CORE_PEER_LOCALMSPID=" + orgName + "1MSP"
			service.Environment[7] = "CORE_PEER_TLS_ENABLED=true"
			service.Environment[8] = "CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/" + orgName + "1." + domainName + "/peers/peer0." + orgName + "1." + domainName + "/tls/server.crt"
			service.Environment[9] = "CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/" + orgName + "1." + domainName + "/peers/peer0." + orgName + "1." + domainName + "/tls/server.key"
			service.Environment[10] = "CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/" + orgName + "1." + domainName + "/peers/peer0." + orgName + "1." + domainName + "/tls/ca.crt"
			service.Environment[11] = "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/" + orgName + "1." + domainName + "/users/Admin@" + orgName + "1." + domainName + "/msp"
			service.WorkingDir = "/opt/gopath/src/github.com/hyperledger/fabric/peer"
			service.Command = "sleep 3600"
			service.Volumes = make([]string, 4)
			service.Volumes[0] = "/var/run/:/host/var/run/"
			service.Volumes[1] = "./../chaincode/:/opt/gopath/src/github.com/chaincode"
			service.Volumes[2] = "./" + orgName + "-artifacts/crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/"
			service.Volumes[3] = "./scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/"
			// service.Volumes[4] = "./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts"
			err := GenDeploy(service, nodeId)
			check(err)

		default:
			log.Fatalf("You didn't specify service name!!..\n")
		}
		dockerCompose.Services[serviceHost] = service
	}
	return nil
}

func GenNetwork(dockerCompose *DockerCompose, networkName string) error {
	network := &Network{
		External: &External{
			Name: networkName,
		},
	}

	dockerCompose.Networks = make(map[string]*Network, 1)
	dockerCompose.Networks[networkName] = network

	return nil
}

func arrayToString(array []string, delim string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array)), delim), "[]")
}

func generatePort(portFor string, basePort string) string {
	rand.Seed(time.Now().UnixNano())
	// port = "36367"
	var appender = rand.Intn(20)
	portString := ""
	for {
		var portInt = rand.Intn(130) + appender
		port := ""

		portStart, _ := strconv.Atoi(string(basePort[0:2]))
		portEnd := string(basePort[len(basePort)-2:])
		if portInt == portStart {
			portInt = rand.Intn(150)
			if portInt == portStart {
				portInt = rand.Intn(150)
				if portInt == portStart {
					portInt = rand.Intn(150)
				}
			}
		}

		if portFor == "ca" {
			port = strconv.Itoa(portInt) + portEnd
			portString = strconv.Itoa(portInt) + portEnd + ":" + basePort
		} else if portFor == "couch" {
			port = strconv.Itoa(portInt) + portEnd
			portString = strconv.Itoa(portInt) + portEnd + ":" + basePort
		} else if portFor == "orderer" {
			port = strconv.Itoa(portInt) + portEnd
			portString = strconv.Itoa(portInt) + portEnd + ":" + basePort
		} else if portFor == "peer" {
			port = strconv.Itoa(portInt) + portEnd
			portString = strconv.Itoa(portInt) + portEnd + ":" + basePort
		}

		portStatus := true

		_, err := strconv.ParseUint(port, 10, 16)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid port %q: %s", port, err)
			os.Exit(1)
		}

		ln, err := net.Listen("tcp", ":"+port)

		if err != nil {
			// fmt.Fprintf(os.Stderr, "Can't listen on port %q: %s", port, err)
			// os.Exit(1)
			portStatus = false
			// return portString
		}

		if portStatus {

			_ = ln.Close()
			break
			// portString
			// fmt.Printf("TCP Port %q is available", port)
			// os.Exit(0)
		}

	}
	return portString
}
