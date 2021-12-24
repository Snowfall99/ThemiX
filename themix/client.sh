echo request number = $1
./single.sh $1 http://127.0.0.1:4100/client &
./single.sh $1 http://127.0.0.1:4110/client &
./single.sh $1 http://127.0.0.1:4120/client &
./single.sh $1 http://127.0.0.1:4130/client &
./single.sh $1 http://127.0.0.1:4140/client &
