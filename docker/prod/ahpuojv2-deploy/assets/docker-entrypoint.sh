#!/bin/bash
update-alternatives --auto java
update-alternatives --auto javac
service supervisor start
service nginx start
/usr/bin/judged
/bin/bash
exit 0
