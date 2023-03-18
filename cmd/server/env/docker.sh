# To start backend services.[mysql minio redis]

# mysql
sudo docker run mysql --name mot-server-mysql -e MYSQL_ROOT_PASSWORD=123456  -p 3306:3306 -d


# minio
mkdir -p ~/minio/data

docker run \
   -p 9000:9000 \
   -p 9090:9090 \
   --name minio \
   -v ~/minio/data:/data \
   -e "MINIO_ROOT_USER=ROOTNAME" \
   -e "MINIO_ROOT_PASSWORD=CHANGEME123" \
   quay.io/minio/minio server /data --console-address ":9090"

# TODO: redis

# TODO: rabbit mq