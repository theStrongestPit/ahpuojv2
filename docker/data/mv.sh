#/bin/bash

for dir in `ls $1`
do
	newdir=$[ $dir - 1000 ]
	echo $newdir
	mv $dir $newdir
done
