#!/bin/bash

redisHealth() {
    redis-cli PING | grep PONG > /dev/null
    if [[ $? -ne 0 ]]; then
        echo "redis ping failed"
        exit 1
    fi
}

consulHealth() {
    consul info | grep "version = 1.2.1"
    if [[ $? -ne 0 ]]; then
        echo "consul not connected"
        exit 1
    fi
}

help() {
    echo "Usage: ./manage.sh redisHealth    => check health of redis"
    echo "       ./manage.sh consulHealth   => check healthy of consul agent"
}

until
    cmd=$1
    if [[ -z "$cmd" ]]; then
        help
    fi
    shift 1
    $cmd "$@"
    [ "$?" -ne 127 ]
do
    help
    exit
done
