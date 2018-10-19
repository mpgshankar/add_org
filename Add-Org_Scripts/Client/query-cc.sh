#!/bin/bash
# Exit on first error
set -e

export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

export CHANNEL_NAME=$1
export ORGNAME=$2
export ORGS=$3
export CCNAME=$4
export PEERS=$5
export DOMAIN=$6
export CCV=$7

for ((number=1; number <= $ORGS; number++))
{
    export orgName=$ORGNAME$ORGS
    echo $orgName
    echo "Query Script Running for "$orgName" ....."  
    for ((numPeers=0; numPeers < $PEERS; numPeers++))
    {
    	export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/$orgName.$DOMAIN/peers/peer$numPeers.$orgName.$DOMAIN/tls/ca.crt

    	export CORE_PEER_ADDRESS=peer$numPeers.$orgName.$DOMAIN:7051

    	peer chaincode query -C $CHANNEL_NAME -n $CCNAME -v $CCV -c '{"Args":["query","b"]}'
    	if [[ $? -ne 0 ]] ; then
        	echo "ERROR !!!! peer chaincode install peer0 "$orgName
        	exit 1
    	fi
    }  

}
