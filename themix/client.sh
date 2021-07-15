echo request number = $1
for count in {1..100}
do
    ./single.sh $1 http://127.0.0.1:11200/client &
    ./single.sh $1 http://127.0.0.1:11210/client &
    ./single.sh $1 http://127.0.0.1:11220/client &
    ./single.sh $1 http://127.0.0.1:11230/client &
    echo $count
done
