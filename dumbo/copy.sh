#!/bin/bash

N=$1
NAME=$2
key='~/.ssh/dumbo'

for ((i = 0; i < N; i++)); do
{
  pubIP0=$(jq ".nodes[$i].PublicIpAddress" nodes.json)
  pubIP=${pubIP0//\"/}

  scp -i $key ubuntu@${pubIP}:/home/ubuntu/dumbo/log/consensus-node-$i.log ./log/$NAME-$i.log &
}
done

wait
