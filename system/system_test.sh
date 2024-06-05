start=$SECONDS
for i in $(seq 1 50); do curl -s -o /dev/null "http://$GATEWAY_URL/hi_every_one/"; done
duration=$(( SECONDS - start ))
echo $duration
