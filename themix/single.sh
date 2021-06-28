start=$SECONDS
for (( i=1; i<=$1; i++ ))
do
        curl -L $2 --data 'a'
done
end=$SECONDS
duration=$(( end - start ))
echo "Themix took $duration seconds to complete"
