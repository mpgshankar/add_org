# Exit on first error
set -e
export DOMAIN="example.com"
# export ORGNAME="org-test"
export ORGS=2
export PEERS=2
export SERVERUSER="dell"
export SERVERIP="192.168.5.153"
export SERVERCONFIGPATH="/home/dell/Cateina/Add-Org/fabric-samples/first-network"
export CLIENTUSER="$USER"
export CLIENTIP="192.168.5.137"
export CHANNEL_NAME="mychannel"
export CCNAME="mycc"
export CCV=2.0
export HOSTNAME=$(hostname)
export CLINODEID=$(docker node ls -f name=$HOSTNAME --format "{{.ID}}")

echo "Hello, '$USER'"

# IP=$(hostname -I)

# echo "My IP :- "$ip4
echo "Type the ORGNAME you want to add in the network, followed by [ENTER]:"
read ORGNAME

echo "Type the Number of organizations you want to add in the network, followed by [ENTER]:"
read ORGS

echo "Type the Number of peers in each organizations you want to add in the network, followed by [ENTER]:"
read PEERS

echo "Type the IP Address of server you want to add in the network, followed by [ENTER]:"
read CLIENTIP


export CLIENTCONFIGPATH="/home/dell/Cateina/Add-Org/fabric-samples/first-network"
# export CLIENTCONFIGPATH=/home/lenovo/Cateina/Add-Org/custom-genConfig
cd $CLIENTCONFIGPATH
# Generate peer, ca, cli, couchdb Crypto.yaml and configtx.yaml 
./custom-genConfig -domain $DOMAIN -orgName $ORGNAME -Orgs $ORGS -Peer $PEERS -NodeId $CLINODEID

# mkdir yaml

cd $ORGNAME-artifacts

mv hyperledger-* ../ 

# Create CryptoGen files using OrgName-crypto.yaml and Cryptogen tool
../../bin/cryptogen generate --config=./$ORGNAME-crypto.yaml
export FABRIC_CFG_PATH=$PWD

for ((number=1; number <= $ORGS; number++))
{
    ../../bin/configtxgen -printOrg $ORGNAME$number'MSP' > ../channel-artifacts/$ORGNAME$number.json
}
# Connect to Server to get Orderer config and copy files from SERVER to Client location
scp -r $SERVERUSER@$SERVERIP:$SERVERCONFIGPATH/crypto-config/ordererOrganizations ./crypto-config/ordererOrganizations

#Copy configtx files to server
cd ../channel-artifacts
scp $orgName*.json $SERVERUSER@$SERVERIP:$SERVERCONFIGPATH/channel-artifacts/

cd ../

docker stack deploy -c hyperledger-couchdb.yaml hyperledger-couchdb
docker stack deploy -c hyperledger-ca.yaml hyperledger-ca
docker stack deploy -c hyperledger-peer.yaml hyperledger-peer
docker stack deploy -c hyperledger-cli.yaml hyperledger-cli

# ssh $SERVERUSER@$SERVERIP 'bash -s' < ./$SERVERCONFIGPATH/server.sh $ORGNAME $ORGS $PEERS
# ssh $SERVERUSER@$SERVERIP "cd $SERVERCONFIGPATH; ./server.sh  \"$ORGNAME\"\"$ORGS\"\"$PEERS\"\"$DOMAIN\""
ssh $SERVERUSER@$SERVERIP "cd $SERVERCONFIGPATH; ./server.sh  $ORGNAME $ORGS $PEERS $DOMAIN"
if [ $? -ne 0 ]; then
    echo "ERROR !!!! ssh to server " failed
    exit 1
fi

./client.sh $ORGNAME $ORGS $PEERS $DOMAIN $CHANNEL_NAME $CCNAME $CCV

ssh $SERVERUSER@$SERVERIP "cd $SERVERCONFIGPATH; ./install-cc.sh  $ORGNAME $ORGS $CHANNEL_NAME $CCNAME $CCV"


exit 0
