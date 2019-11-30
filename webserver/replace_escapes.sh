#!/bin/bash

arg1=$1
arg2=$2

sed -e "s/\x1b//g" $1 | sed -e "s/\x0D//g" > $2

