#!/bin/bash

psname="xsbadmin_linux64"

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

./$psname &