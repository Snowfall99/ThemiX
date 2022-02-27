echo request number = $1
./single.sh $1 1 http://127.0.0.1:11300/client &
./single.sh $1 2 http://127.0.0.1:11310/client &
./single.sh $1 3 http://127.0.0.1:11320/client &
./single.sh $1 4 http://127.0.0.1:11330/client &
./single.sh $1 5 http://127.0.0.1:11340/client &
