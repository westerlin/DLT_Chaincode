#!/bin/bash
if [ $# -eq 0 ] ; then
 echo "You have not set a logging level"
 echo "Choose from :"
 echo "  CRITICAL | ERROR | WARNING | NOTICE | INFO | DEBUG "
else
 export CORE_LOGGING_LEVEL=$1
 echo "Logging level set at $1"
fi