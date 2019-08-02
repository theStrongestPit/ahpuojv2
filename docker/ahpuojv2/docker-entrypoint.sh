#!/bin/bash
# /usr/bin/judged
service nginx start
update-alternatives --auto java
update-alternatives --auto javac

/bin/bash  
exit 0