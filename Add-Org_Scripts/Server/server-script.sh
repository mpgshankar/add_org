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

export ORGNAME=$2
export ORGS=$3
export PEERS=$4
export DOMAIN=$5


export EXORGNAMEMSP="Org"
export EXORGNAME="org"
export EXORGS=2
export EXPEERS=2
export EXDOMAIN="example.com"

configtxlator start &
if [[ $? -ne 0 ]] ; then
    exit 130
fi

apt update && apt install jq -y

export CONFIGTXLATOR_URL=http://127.0.0.1:7059
echo $CONFIGTXLATOR_URL

export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$DOMAIN/orderers/orderer0.$DOMAIN/msp/tlscacerts/tlsca.$DOMAIN-cert.pem

export CHANNEL_NAME=$1

for ((numOrgs=1; numOrgs <= $ORGS; numOrgs++))
{
        export ORGNAMEMSP=$ORGNAME$numOrgs
        peer channel fetch config config_block.pb -o orderer0.$DOMAIN:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
        if [[ $? -ne 0 ]] ; then
            echo "ERROR !!!! peer channel fetch config config_block.pb"
            exit 1
        fi

        curl -X POST --data-binary @config_block.pb "$CONFIGTXLATOR_URL/protolator/decode/common.Block" | jq . > config_block.json

        jq .data.data[0].payload.data.config config_block.json > config.json


        jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups":{"'$ORGNAMEMSP'MSP":.[1]}}}}}' config.json ./channel-artifacts/$ORGNAMEMSP.json >& updated_config.json

        curl -X POST --data-binary @config.json "$CONFIGTXLATOR_URL/protolator/encode/common.Config" > config.pb

        curl -X POST --data-binary @updated_config.json "$CONFIGTXLATOR_URL/protolator/encode/common.Config" > updated_config.pb

        curl -X POST -F channel=$CHANNEL_NAME -F "original=@config.pb" -F "updated=@updated_config.pb" "${CONFIGTXLATOR_URL}/configtxlator/compute/update-from-configs" > config_update.pb

        curl -X POST --data-binary @config_update.pb "$CONFIGTXLATOR_URL/protolator/decode/common.ConfigUpdate" | jq . > config_update.json

        echo '{"payload":{"header":{"channel_header":{"channel_id":"mychannel","type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json

        curl -X POST --data-binary @config_update_in_envelope.json "$CONFIGTXLATOR_URL/protolator/encode/common.Envelope" > config_update_in_envelope.pb

        peer channel signconfigtx -f config_update_in_envelope.pb
        if [[ $? -ne 0 ]] ; then
            echo "ERROR !!!! peer channel signconfigtx"
            exit 1
        fi

        export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/$EXORGNAME$EXORGS.$DOMAIN/users/Admin\@$EXORGNAME$EXORGS.$DOMAIN/msp/ && export CORE_PEER_ADDRESS=peer0.$EXORGNAME$EXORGS.$DOMAIN:7051 && export CORE_PEER_LOCALMSPID=$EXORGNAMEMSP$EXORGS"MSP" && export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/$EXORGNAME$EXORGS.$DOMAIN/peers/peer0.$EXORGNAME$EXORGS.$DOMAIN/tls/ca.crt

        peer channel update -f config_update_in_envelope.pb -c $CHANNEL_NAME -o orderer0.$DOMAIN:7050 --tls true --cafile $ORDERER_CA
        if [[ $? -ne 0 ]] ; then
            echo "ERROR !!!! peer channel update -f config_update_in_envelope.pb"
            exit 1
        fi

        mkdir $ORGNAMEMSP
        mv *.json $ORGNAMEMSP/
        mv *.pb $ORGNAMEMSP/
}

ls -la