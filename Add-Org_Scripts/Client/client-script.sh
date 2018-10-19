#!/bin/bash
# Exit on first error
set -e
echo "                                         "
echo " ____    _____      _      ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|"
echo "\___ \    | |     / _ \   | |_) |   | |  "
echo " ___) |   | |    / ___ \  |  _ <    | |  "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|  "
echo "                                         "

export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export CCV=$7
export CHANNEL_NAME=$5
export CCNAME=$6
export CCPATH=github.com/chaincode/chaincode_example02/go/
export DOMAIN=$1
export ORGNAME=$2
export ORGS=$3
export PEERS=$4

echo $CCV
echo $CHANNEL_NAME+"CHANNEL_NAME"
echo $CCNAME
echo $CCPATH
echo $DOMAIN+"DOMAIN"
echo $ORGNAME
echo $ORGS
echo $PEERS

peer channel fetch 0 mychannel.block -o orderer0.$DOMAIN:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
if [[ $? -ne 0 ]] ; then
    echo "ERROR !!!! peer channel fetch 0 mychannel.block"
    exit 1
fi

peer channel join -b mychannel.block
if [[ $? -ne 0 ]] ; then
    echo "ERROR !!!! peer channel join peer0 mychannel.block"
    exit 1
fi

export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/$ORGNAME.$DOMAIN/peers/peer1.$ORGNAME.$DOMAIN/tls/ca.crt

export CORE_PEER_ADDRESS=peer1.$ORGNAME.$DOMAIN:7051

peer channel join -b mychannel.block
if [[ $? -ne 0 ]] ; then
    echo "ERROR !!!! peer channel join peer1 mychannel.block"
    exit 1
fi

for ((numPeers=0; numPeers < $PEERS; numPeers++))
{
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/$ORGNAME.$DOMAIN/peers/peer$numPeers.$ORGNAME.$DOMAIN/tls/ca.crt

    export CORE_PEER_ADDRESS=peer$numPeers.$ORGNAME.$DOMAIN:7051

    peer chaincode install -n $CCNAME -v $CCV -p $CCPATH
    if [[ $? -ne 0 ]] ; then
        echo "ERROR !!!! peer chaincode install peer0 "$ORGNAME
        exit 1
    fi
}
