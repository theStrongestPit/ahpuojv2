#!/bin/bash

#!/bin/bash
# 设置go运行环境
# echo 'export GOROOT=/usr/lib/go-1.10' >> /etc/profile
# alter user 'root'@'%' identified by 'jsj123zxc' password expire never;
# alter user 'root'@'%' identified with mysql_native_password by 'jsj123zxc';
# flush privileges;

# 创建用户
/usr/sbin/useradd -m -u 1536 judge

# 创建web目录文件夹
mkdir /home/judge/web
mkdir /home/judge/web/config
mkdir /home/judge/web/log
mkdir /home/judge/web/upload
mkdir /home/judge/web/upload/avatars
mkdir /home/judge/web/upload/images
mkdir /home/judge/web/static

# 创建评测机目录文件夹
mkdir /home/judge/hustoj_core
cd /home/judge/hustoj_core
mkdir etc log
mkdir run0 run1 run2 run3
chown judge run0 run1 run2 run3

# 赋予评测机源码部署脚本执行权限并部署
chmod +x /home/judge/core/make.sh
cd /home/judge/core/ && ./make.sh
cd /usr/bin && rm awk && cp -s mawk awk

# 添加各种配置文件
cp /home/judge/install/config.ini /home/judge/web/config
cp /home/judge/install/auth_model.conf /home/judge/web/config
cp /home/judge/install/judge.conf.example /home/judge/hustoj_core/etc/judge.conf
cp /home/judge/install/java0.policy.example /home/judge/hustoj_core/etc/java0.policy

# 各种权限设置
chmod 775 -R /home/judge/hustoj_core/data
chmod 775 -R /home/judge/web/upload
chown -r judge -R /home/judge/hustoj_core/data
chown -r judge -R /home/judge/web/upload

# 修改配置文件
CPU=$(grep "cpu cores" /proc/cpuinfo | head -1 | awk '{print $4}')
cd /home/judge/hostoj_core
if [ $(grep -c "client_max_body_size" /etc/nginx/nginx.conf) -eq 0 ]; then
  sed -i "s:include /etc/nginx/mime.types;:client_max_body_size    80m;\n\tinclude /etc/nginx/mime.types;:g" /etc/nginx/nginx.conf
fi

sed -i "s/OJ_RUNNING=1/OJ_RUNNING=$CPU/g" /home/judge/hustoj_core/etc/judge.conf
sed -i "s/OJ_USER_NAME=root/OJ_USER_NAME=$MYSQL_USER/g" /home/judge/hustoj_core/etc/judge.conf
sed -i "s/OJ_PASSWORD=root/OJ_PASSWORD=$MYSQL_PASSWORD/g" /home/judge/hustoj_core/etc/judge.conf

# 初始化mysql数据库 修改mysql8.0的密码验证方式以使用navicat链接
mysql -h db -uroot -p$MYSQL_PASSWORD <<EOF
alter user 'root'@'%' identified by 'jsj123zxc' password expire never;
alter user 'root'@'%' identified with mysql_native_password by 'jsj123zxc';
flush privileges;
EOF

# 初始化数据库结构
mysql -h db -uroot -p$MYSQL_PASSWORD </home/judge/install/ahpuoj.sql
mysql -h db -uroot -p$MYSQL_PASSWORD <<EOF
CREATE DATABASE IF NOT EXISTS casbin;
use casbin;
CREATE TABLE IF NOT EXISTS casbin_rule (
  p_type varchar(100) DEFAULT NULL,
  v0 varchar(100) DEFAULT NULL,
  v1 varchar(100) DEFAULT NULL,
  v2 varchar(100) DEFAULT NULL,
  v3 varchar(100) DEFAULT NULL,
  v4 varchar(100) DEFAULT NULL,
  v5 varchar(100) DEFAULT NULL,
  KEY IDX_casbin_rule_v0 (v0),
  KEY IDX_casbin_rule_v1 (v1),
  KEY IDX_casbin_rule_v2 (v2),
  KEY IDX_casbin_rule_v3 (v3),
  KEY IDX_casbin_rule_v4 (v4),
  KEY IDX_casbin_rule_v5 (v5),
  KEY IDX_casbin_rule_p_type (p_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
INSERT INTO casbin_rule VALUES ('p', 'admin', '/*', '*', NULL, NULL, NULL);
INSERT INTO casbin_rule VALUES ('p', 'subadmin', '/problem/*', '*', NULL, NULL, NULL);
INSERT INTO casbin_rule VALUES ('p', 'subadmin', '/contest/*', '*', NULL, NULL, NULL);
INSERT INTO casbin_rule VALUES ('p', 'subadmin', '/tag/*', '*', NULL, NULL, NULL);
INSERT INTO casbin_rule VALUES ('p', 'subadmin', '/user/*', '*', NULL, NULL, NULL);
INSERT INTO casbin_rule VALUES ('p', 'subadmin', '/generator/*', '*', NULL, NULL, NULL);
EOF

# 语言运行环境设置
echo /usr/lib/jvm/java-8-openjdk-amd64/lib/amd64/jli/ >>/etc/ld.so.conf
ldconfig
ln -s /usr/bin/python3.7 /usr/bin/python3

service supervisor start
service nginx start
/usr/bin/judged
