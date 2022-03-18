start=$SECONDS
for (( i=1; i<=$1; i++ ))
do
        /opt/homebrew/opt/curl/bin/curl $3 --data 'a'
done
end=$SECONDS
duration=$(( end - start ))
echo "Themix client $2 took $duration seconds to complete"
