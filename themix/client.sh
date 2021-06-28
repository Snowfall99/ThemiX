echo request number = $1
./single.sh $1 http://127.0.0.1:11200/client &
./single.sh $1 http://127.0.0.1:11300/client &
./single.sh $1 http://127.0.0.1:11400/client &
./single.sh $1 http://127.0.0.1:11500/client &
