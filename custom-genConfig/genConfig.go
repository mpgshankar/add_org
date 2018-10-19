package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var overlayNetwork = "hyperledger-ov"

func main() {
	var domain, orgName, nodeId string
	var numOrgs, numPeer, numOrderer, numKafka, numZookeeper int
	// domain = "example.com"
	flag.StringVar(&domain, "domain", domain, "Generate config file for a particular doamin")
	flag.StringVar(&orgName, "orgName", orgName, "Generate config file for a particular orgName")
	flag.IntVar(&numOrgs, "Orgs", numOrgs, "Choose number of Organizations except Orderer's Organization. CA will be created per each organization")
	flag.IntVar(&numPeer, "Peer", numPeer, "Choose number of peers per organizations")
	flag.StringVar(&nodeId, "NodeId", nodeId, "Generate config file for a particular node")
	// flag.IntVar(&numOrderer, "Orderer", 1, "Choose number of orderers (if set, need to specify number of Kafka nodes)")
	// flag.IntVar(&numKafka, "Kafka", 3, "Choose number of kafka nodes")
	// flag.IntVar(&numZookeeper, "Zookeeper", 3, "Choose number of zookeeper nodes")

	flag.Parse()

	// Generate crypto-config.yaml
	crypto, err := GenCrypto(domain, orgName, numOrgs, numPeer, numOrderer)
	fmt.Println("Generating YAML file from crypto config....")
	cryptoYAML, err := yaml.Marshal(&crypto)
	check(err)

	// Generate configtx.yaml
	configtx, err := GenConfigtx(domain, orgName, numOrgs, numOrderer, numKafka)
	check(err)
	fmt.Println("Generating YAML file from configtx config....")
	configtxYAML, err := yaml.Marshal(&configtx)
	check(err)

	path := orgName + "-artifacts"

	// Write files to $PWD
	newpath := filepath.Join(".", path)
	os.MkdirAll(newpath, os.ModePerm)
	pwd, err := filepath.Abs(newpath)
	check(err)
	err = ioutil.WriteFile(pwd+"/"+orgName+"-crypto.yaml", []byte(cryptoYAML), 0644)
	check(err)
	err = ioutil.WriteFile(pwd+"/"+"configtx.yaml", []byte(configtxYAML), 0644)
	check(err)

	// Genearte docker composer file
	var composeOutput *DockerCompose
	var serviceList []string

	// if numOrderer == 1 {
	serviceList = make([]string, 4)
	serviceList = []string{"ca", "couchdb", "peer", "cli"}
	// } else {
	// 	serviceList = make([]string, 7)
	// 	serviceList = []string{"zookeeper", "kafka", "orderer", "ca", "couchdb", "peer", "cli"}
	// }

	for _, service := range serviceList {
		switch service {
		case "peer":
			composeOutput, err = GenDockerCompose(service, domain, orgName, overlayNetwork, nodeId, numPeer, numOrgs)
			check(err)
		case "zookeeper":
			composeOutput, err = GenDockerCompose(service, domain, orgName, overlayNetwork, nodeId, numZookeeper)
			check(err)
		case "kafka":
			composeOutput, err = GenDockerCompose(service, domain, orgName, overlayNetwork, nodeId, numKafka)
			check(err)
		case "orderer":
			composeOutput, err = GenDockerCompose(service, domain, orgName, overlayNetwork, nodeId, numOrderer)
			check(err)
		case "ca":
			composeOutput, err = GenDockerCompose(service, domain, orgName, overlayNetwork, nodeId, numOrgs)
			check(err)
		case "couchdb":
			composeOutput, err = GenDockerCompose(service, domain, orgName, overlayNetwork, nodeId, numPeer, numOrgs)
			check(err)
		case "cli":
			composeOutput, err = GenDockerCompose(service, domain, orgName, overlayNetwork, nodeId, 1)
			check(err)
		default:
			panic("Service Name isn't specified!!!")
		}
		fmt.Println("Generating Docker Compose file for " + service + "....")
		composeYAML, err := yaml.Marshal(composeOutput)
		check(err)
		err = ioutil.WriteFile(pwd+"/"+"hyperledger-"+service+".yaml", []byte(composeYAML), 0644)
		check(err)
	}

	fmt.Println("Output files are located on " + pwd)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
