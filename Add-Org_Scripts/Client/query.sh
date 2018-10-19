# Exit on first error
set -e
export CONTAINERPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer

echo $CLIENTCLIID
export LANGUAGE="golang"
export CLI_DELAY=3
export CLI_TIMEOUT=10
export ORGNAME=$1
export ORGS=$2
export CHANNEL_NAME=$3
export CCNAME=$4
export PEERS=$5
export DOMAIN=$6
export CCV=$7

export ORGS=$(($2 + 0))
export PEERS=$(($5 + 0))
export CLIENTCLIID=$(docker ps --filter="name=cli_$ORGNAME" -aq | xargs)
docker cp query-cc.sh $CLIENTCLIID:$CONTAINERPATH
# docker exec -it $SERVERCLIID bash
docker exec $CLIENTCLIID ./query-cc.sh $CHANNEL_NAME $ORGNAME $ORGS $CCNAME $PEERS $DOMAIN $CCV
if [ $? -ne 0 ]; then
    echo "ERROR !!!! docker exec "$SERVERCLIID failed
    exit 1
fi

exit 0
