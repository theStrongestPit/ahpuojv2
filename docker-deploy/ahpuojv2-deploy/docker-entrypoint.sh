#!/bin/bash
# /usr/bin/judged
service nginx start
update-alternatives --auto java
update-alternatives --auto javac
/usr/bin/judged
/bin/bash  
exit 0s