#!/bin/sh

RED='\033[0;31m'

while true; do
    echo 'Building...'
    if go build ; then
        ./tfbi $@ &
        PID=$!
        echo "Running app in the background: PID=$PID"
    else
        unset PID
        echo -e "${RED}Build Failed!"
    fi

    inotifywait \
        -e create -e delete -e modify -e move \
        *.go go.* ./internal/**

    if [ $PID ] ; then
        echo "Kill background service: PID=$PID"
        kill $PID
    fi
done
