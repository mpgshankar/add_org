# Exit on first error
set -e
export CONTAINERPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer

echo $CLIENTCLIID
export LANGUAGE="golang"
export CLI_DELAY=3
export CLI_TIMEOUT=10
export CHANNEL_NAME=$5
export CCNAME=$6
export CCV=$7

#export CLIENTUSER=$1
#export CLIENTIP=$2
export ORGNAME=$1
export ORGS=$2
export PEERS=$3
export DOMAIN=$4
export CLIENTCLIID=$(docker ps --filter="name=cli_$ORGNAME" -aq | xargs)
docker cp client-script.sh $CLIENTCLIID:$CONTAINERPATH
# docker exec -it $CLIENTCLIID bash
for ((number=1; number <= $ORGS; number++))
{
    export orgName=$ORGNAME$ORGS
    echo $orgName
    echo "Client Script Running for "$orgName" ....."    
    #docker exec -i $CLIENTCLIID bash < client-script.sh
    docker exec -i $CLIENTCLIID /bin/bash -c "./client-script.sh '$DOMAIN' '$orgName' '$ORGS' '$PEERS' '$CHANNEL_NAME' '$CCNAME' '$CCV'"
    if [ $? -ne 0 ]; then
        echo "ERROR !!!! docker exec "$CLIENTCLIID failed
        exit 1
    fi
    
    #sleep 60
    echo "Client Script Completed Running for "$orgName"!!"

}
