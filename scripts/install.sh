#!/bin/bash

singlayer_root=/opt/singlayer

mkdir -p $singlayer_root

if [ -f $singlayer_root/singlayer ]; then
	mv $singlayer_root/singlayer $singlayer_root/singlayer.old
fi

buffalo build -o $singlayer_root/singlayer

mkdir -p $singlayer_root/templates
cp -a templates/order.templ $singlayer_root/templates/
