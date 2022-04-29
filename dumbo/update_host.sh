#!/bin/bash

N=$1
key='~/.ssh/dumbo'
host='tmp_hosts.config'

rm $host
for ((i = 0; i < N; i++)); do
{
  pubIP0=$(jq ".nodes[$i].PublicIpAddress" nodes.json)
  pubIP=${pubIP0//\"/}

  priIP0=$(jq ".nodes[$i].PrivateIpAddress" nodes.json)
  priIP=${priIP0//\"/}

  echo $i ${priIP} ${pubIP} $(( $((200 * $i)) + 5000 )) >> tmp_hosts.config
}
done
