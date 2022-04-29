#!/bin/bash

N=$1
B=$2
K=$3
key='~/.ssh/dumbo'

for ((i = 0; i < N; i++)); do
{
  pubIP0=$(jq ".nodes[$i].PublicIpAddress" nodes.json)
  pubIP=${pubIP0//\"/}

  ssh -i $key ubuntu@${pubIP} "export LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib; export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib; cd dumbo; nohup python3 run_socket_node.py --sid 'sidA' --id $i --N $N --f $(( (N-1)/3 )) --B $B --K $K --S 50 --T 2 --P "dumbo" --F 1000000 > node-$i.out" &
} &
done

wait
