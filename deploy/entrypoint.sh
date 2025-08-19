#!/bin/sh
# Persist runtime environments
printenv >> /etc/environment
exec ./main
