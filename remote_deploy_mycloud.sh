#/bin/bash

# scp  index.html ubuntu@www.jiezi19971225.cn:/home/ubuntu/docker-deploy/web/
# scp  admin_index.html ubuntu@www.jiezi19971225.cn:/home/ubuntu/docker-deploy/web/
# scp -r dist ubuntu@www.jiezi19971225.cn:/home/ubuntu/docker-deploy/web/
# scp -r static ubuntu@www.jiezi19971225.cn:/home/ubuntu/docker-deploy/web/
# scp -r config ubuntu@www.jiezi19971225.cn:/home/ubuntu/docker-deploy/web/
# scp ahpuoj ubuntu@www.jiezi19971225.cn:/home/ubuntu/docker-deploy/web/ahpuoj
# scp -r /home/jiezi19971225/ojdataback/data ubuntu@www.jiezi19971225.cn:/home/ubuntu/docker-deploy/
# scp -r core ubuntu@www.jiezi19971225.cn:/home/ubuntu/docker-deploy/
# scp -r docker-deploy root@172.16.0.3:/root/ahpuojv2-docker

scp dist/index.html root@172.16.0.3:/root/ahpuojv2docker/web/
scp dist/admin_index.html root@172.16.0.3:/root/ahpuojv2docker/web/
scp -r dist root@172.16.0.3:/root/ahpuojv2docker/web/
# scp -r static root@172.16.0.3:/root/ahpuojv2docker/web/
# scp -r config root@172.16.0.3:/root/ahpuojv2docker/web/
# scp -r core root@172.16.0.3:/root/ahpuojv2docker/
# scp ahpuoj root@172.16.0.3:/root/ahpuojv2docker/web/ahpuoj
