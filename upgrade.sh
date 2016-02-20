#!/bin/bash

psname="xsbadmin"

echo "PS name is $psname"

ps=`pgrep $psname`
# echo $ps

if [ "$ps" ]
then 
	echo "$psname was running" 
	kill $ps
	echo "kill $psname OK"
else
	echo "$psname not running"
fi

git pull

process_file_name=${psname}"_linux64"

nohup ./$process_file_name &