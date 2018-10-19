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
export CCV=$5
export CHANNEL_NAME=$1
export CCNAME=$4
export CCPATH=github.com/chaincode/chaincode_example02/go/
export ORGNAME=$2
export ORGS=$3
export NEWORGMSP=test1

export EXISTMSP = ('Org1MSP.member', 'Org2MSP.member')

peer chaincode install -n $CCNAME -v $CCV -p $CCPATH
if [[ $? -ne 0 ]] ; then
    echo "ERROR !!!! peer chaincode install peer0 "$ORGNAME
    exit 1
fi

export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin\@org2.example.com/msp/ && export CORE_PEER_ADDRESS=peer0.org2.example.com:7051 && export CORE_PEER_LOCALMSPID="Org2MSP" && export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

peer chaincode install -n $CCNAME -v $CCV -p $CCPATH
if [[ $? -ne 0 ]] ; then
    echo "ERROR !!!! peer chaincode install peer0 "$ORGNAME
    exit 1
fi

export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin\@org1.example.com/msp/ && export CORE_PEER_ADDRESS=peer0.org1.example.com:7051 && export CORE_PEER_LOCALMSPID="Org1MSP" && export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

peer chaincode upgrade -o orderer0.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CCNAME -v $CCV -c '{"Args":["init","a","90","b","210"]}' -P "OR('Org1MSP.member', 'Org2MSP.member', 'test1MSP.member')"
