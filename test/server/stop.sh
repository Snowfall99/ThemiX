#!/bin/bash

ps -u `whoami` | grep main | awk '{system("kill -9 "$1)}'
