# Exit on first error
set -e
export CONTAINERPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer
export SERVERCLIID=$(docker ps --filter="name=-cli" -aq | xargs)
echo $SERVERCLIID
export LANGUAGE="golang"
export CLI_DELAY=3
export CLI_TIMEOUT=10
export CHANNEL_NAME="mychannel"

export ORGNAME=$1
export ORGS=$2
export PEERS=$3
export DOMAIN=$4

export ORGS=$(($2 + 0))
export PEERS=$(($3 + 0))

docker cp server-script.sh $SERVERCLIID:$CONTAINERPATH
# docker exec -it $SERVERCLIID bash
docker exec $SERVERCLIID ./server-script.sh $CHANNEL_NAME $ORGNAME $ORGS $PEERS $DOMAIN
if [ $? -ne 0 ]; then
    echo "ERROR !!!! docker exec "$SERVERCLIID failed
    exit 1
fi
