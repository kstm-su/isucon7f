cp /var/log/nginx/access.log /var/log/nginx/access.log.1
cp /var/log/kataribe.log /var/log/kataribe.log.1
cat /var/log/nginx/access.log | /opt/kataribe -conf /opt/kataribe.toml > /var/log/kataribe.log

>/var/log/nginx/access.log

