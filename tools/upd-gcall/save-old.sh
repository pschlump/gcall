#!/bin/bash

SaveDir=~/.upd-gcall-old

# echo $SaveDir
mkdir -p $SaveDir

ls $SaveDir >/tmp/,a
lfn=$( tail -1 /tmp/,a )
if [ "${lfn}" == "" ] ; then
	pwd >/tmp/,b
	date >>/tmp/,b
	echo "===" >>/tmp/,b
	cat gcall.addr.json >>/tmp/,b
	mv /tmp/,b $SaveDir/0000001.gcall.addr.json
else
	seq=$( echo ${lfn} | sed -e 's/.gcall.*//' )
	let "seq=seq+1"
	#echo seq=$seq
	pseq=$( printf "%07d\n" "${seq}" )
	#echo pseq=$pseq

	pwd >/tmp/,b
	date >>/tmp/,b
	echo "===" >>/tmp/,b
	cat gcall.addr.json >>/tmp/,b
	mv /tmp/,b "$SaveDir/${pseq}.gcall.addr.json"
fi

cp gcall.addr.cfg gcall.addr.json

