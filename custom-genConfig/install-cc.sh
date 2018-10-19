# Exit on first error
set -e
export CONTAINERPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer
export SERVERCLIID=$(docker ps --filter="name=cli" -aq | xargs)
echo $SERVERCLIID
export LANGUAGE="golang"
export CLI_DELAY=3
export CLI_TIMEOUT=10
export ORGNAME=$1
export ORGS=$2
export CHANNEL_NAME=$3
export CCNAME=$4
export CCV=$5

export ORGS=$(($2 + 0))
# export PEERS=$(($3 + 0))

docker cp install-instantiate-cc.sh $SERVERCLIID:$CONTAINERPATH
# docker exec -it $SERVERCLIID bash
docker exec $SERVERCLIID ./install-instantiate-cc.sh $CHANNEL_NAME $ORGNAME $ORGS $CCNAME $CCV
if [ $? -ne 0 ]; then
    echo "ERROR !!!! docker exec "$SERVERCLIID failed
    exit 1
fi
