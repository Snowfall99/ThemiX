#!/bin/bash

NUM=$1
BATCH=$2
TIME=$3

for ((i = 0; i < NUM; i++)); do
{
  host1=$(jq ".nodes[$i].PublicIpAddress" clients.json)
  host=${host1//\"/}
  url1=$(jq ".nodes[$i].ServerURL" clients.json)
  url=${url1//\"/}
  port=6000
  user='ubuntu'
  key="~/.ssh/aws"
  id=$i
  node="node"$id

expect <<-END
spawn ssh -oStrictHostKeyChecking=no -i $key $user@$host "cd;cd client;./client  -url $url -batch $BATCH -time $TIME > output"
expect EOF
exit
END
} &
done

wait
