# ThemiX

## Introduction
This is an implementation of `ThemiX: a novel timing-balanced consensus protocol`.

## Content
```
.
|- Makefile
|- README.md
|- acs
|   |-Makefile
|   |-README.md
|   |- src
|   |- script
|
|- src
|   |- client
|   |- crypto
|   |- Makefile
|   |- themix
|   |- transport
|
|- test
    |- client
    |- server
```

## Getting Started

### Build ThemiX
You can build the whole project as follows
```bash
cd src/
make
```
We set `n=5 and t=3` as default. You can refer to [Build Step by Step](#build-step-by-step) for more information.

### Build Step by Step
1. Build ThemiX
    ```bash
    # Download dependencies.
    cd src/themix
    go get -v -t -d ./...

    # Build ThemiX
    go build -o main server/cmd/main.go
    ```
    - Verification of client request signature is set to `false` as default. You can use flag `-sign true` to change it.
    - $2\Delta$ is set as 2500ms and $2\delta$ is set to 500ms as default. You can refer to [`server/instance.go`](src/themix/server/instance.go) line 99 and 100 to change them.
2. Generate BLS Keys and ECDSA Key
    ```bash
    # Download dependencies.
    cd src/crypto
    go get -v -t -d ./...

    # Generate bls keys.
    # Default n=5, t=3
    go run cmd/bls/main.go -n 5 -t 3

    # Generate ECDSA key.
    go run cmd/ecdsa/main.go
    ```
3. Build Client
    ```bash
    # Compile protobuf.
    cd src/client/proto
    make

    # Build Client.
    cd src/client
    go build -o client main.go
    ```
    `client` binary is supposed to be copied to directory `test/client` for further use.

## Testing

**First of all, you are supposed to change `key` in every shell script to the path of your own ssh key to access your AWS account.** Our access key is named `aws`, and you can change it to your own key.

### Run ThemiX on a cluster of AWS EC2 machines
1. **Create a cluster of AWS EC2 machines.** You may refer to [narwhal](https://github.com/facebookresearch/narwhal) for help. We create 5 `t3.2xlarge` instances in region `us-east-1`, `us-west-1`, `ap-northeast-1`, `ap-southeast-1` and `ap-east-1`, and we open `Port 5000-7000`.

2. **Fetch machine information from AWS.** First, `test/server/aws.py` is supposed to be modified:
      - line 5 should contain the regions of your instances
      - line 6 and line 7 is supposed to be set according to your own AWS secret keys
      - line 23 is the name of your AWS secret key to access other instances
    
    After modification, you can fetch instance information as following
    ```bash
    cd test/server
    python aws.py
    ```
    

3. **Generate config file(`node.json`) for every node.**
    ```bash
    python generate.py
    ```

4. **Deliver nodes.** Again, value `key` in every shell script is supposed to be the path of your AWS secret key to access other instances.
    ```bash
    chmod +x *.sh

    # Compress BLS keys.
    ./tarKeys.sh
    
    # Deliver to every node.
    # n is the number of nodes in the cluster.
    ./deliverNode.sh n
    ```

5. **Run nodes.**
   ```bash
   ./beginNode.sh n <batchsize>
   ```

6. **Stop nodes.**
   ```bash
   ./stopNode.sh n
   ```

### Run clients to send requests to ThemiX nodes.
1. Deliver client.
   ```bash
   chmod +x *.sh
   
   # n is the number of nodes in the cluster
   ./deliverNode.sh n
   ```

2. Run client for a period of time.
   ```bash
   ./beginNode.sh <size of payload> <size of batch> <running time>
   ```

3. Copy result from client node.
   ```bash
   ./createDir.sh n
   ./copyResult.sh n <name of log file>
   ```

4. Calculate throughput and latency.
   ```bash
   python cal.py <number of nodes> <batchsize> <path> <name of log file>
   ```

### A Brief introduction of test scripts

#### script/server

* aws.py: get machine information from AWS.
  > You may need to change line 5-7 and 23 to your own config.

* generate.py：generate configuration for every node.

* tarKeys.sh: compress BLS keys.

* deliverNode.sh: deliver node to remote machines.
  * `./deliverNode.sh <the number of remote machines>`

* beginNode.sh: run node on remote machines.
  * `./beginNode.sh <the number of remote machines>`

* stopNode.sh: stop node on remote machines.
  * `./stopNode.sh <the number of remote machines>`

* Simulate the crash of specific nodes.
  * You can refer to `crash33.sh`, and change `ADDR` to your specific situation.
  * `./crash33.sh`

* Simulate the crash of the last few nodes.
  * eg. simulate the crash of the last 33 nodes of 100 nodes.
    ```bash
    ./crash.sh 100 33
    ```

* rmLog.sh: remove log file on remote machines.
  * `./rmLog.sh <the number of remote machines>`

### script/client
* deliverNode.sh: deliver client to remote machines.
  * `./deliverNode.sh <the number of remote machines>`

* Get log files from remote machines
    ```bash
    mkdir log
    ./createDir.sh <the number of remote machines>
    ./copyResult.sh <the number of remote machines> <name of log files>
    ```

* beginNode.sh: run client on remote machines.
  * These args are supposed to be assigned.
    * -payload: the size of a single request
    * -time: running time
    * -batch: the size of batch

## ACS
We also implement ACS in directory `acs/`. You can refer to [`acs/README.md`](./acs/README.md) for more detailed information.
