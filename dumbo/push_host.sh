#!/bin/bash

N=$1
key='~/.ssh/dumbo'

for ((i = 0; i < N; i++)); do
{
  pubIP0=$(jq ".nodes[$i].PublicIpAddress" nodes.json)
  pubIP=${pubIP0//\"/}

  ssh -o "StrictHostKeyChecking no" -i $key ubuntu@${pubIP} "rm /home/ubuntu/dumbo/hosts.config"
  scp -i $key tmp_hosts.config ubuntu@${pubIP}:/home/ubuntu/dumbo/hosts.config &  
} &
done

wait
