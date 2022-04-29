#!/bin/bash

N=$1
key='~/.ssh/dumbo'

for ((i = 0; i < N; i++)); do 
{
  pubIP0=$(jq ".nodes[$i].PublicIpAddress" nodes.json)
  pubIP=${pubIP0//\"/}
  
  ssh -oStrictHostKeyChecking=no -i $key ubuntu@${pubIP} "git clone https://github.com/yylluu/dumbo.git"
} &
done

wait
