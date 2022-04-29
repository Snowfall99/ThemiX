#!/bin/bash

N=$1
key='~/.ssh/dumbo'

for ((i = 0; i < N; i++)); do 
{
  pubIP0=$(jq ".nodes[$i].PublicIpAddress" nodes.json)
  pubIP=${pubIP0//\"/}
  
ssh -oStrictHostKeyChecking=no -i $key ubuntu@${pubIP} "
sudo apt-get update;
sudo apt-get -y install make bison flex libgmp-dev libmpc-dev python3 python3-dev python3-pip libssl-dev
"

ssh -oStrictHostKeyChecking=no -i $key ubuntu@${pubIP} "
python3 -m pip install --upgrade pip;
sudo pip3 install gevent setuptools gevent numpy ecdsa pysocks gmpy2 zfec gipc pycrypto coincurve
"

ssh -oStrictHostKeyChecking=no -i $key ubuntu@${pubIP} "
wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz;
tar -xvf pbc-0.5.14.tar.gz;
cd pbc-0.5.14;
sudo ./configure;
sudo make;
sudo make install;
cd ..
"

ssh -oStrictHostKeyChecking=no -i $key ubuntu@${pubIP} "
sudo ldconfig /usr/local/lib;
cat <<EOF >/home/ubuntu/.profile
export LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
EOF
"

ssh -oStrictHostKeyChecking=no -i $key ubuntu@${pubIP} "
source /home/ubuntu/.profile;
export LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib;
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib;
"

ssh -oStrictHostKeyChecking=no -i $key ubuntu@${pubIP} "
git clone https://github.com/JHUISI/charm.git;
cd charm;
sudo ./configure.sh;
sudo make;
sudo make install;
"
} &
done 

wait
