#!/bin/sh
cd "$(dirname "$0")"
java -Xmx4G -Xms4G -jar server.jar nogui
