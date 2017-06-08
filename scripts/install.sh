#!/bin/bash

if [ -f /opt/singlayer/singlayer ]; then
	mv /opt/singlayer/singlayer /opt/singlayer/singlayer.old
fi

buffalo build -o /opt/singlayer/singlayer
