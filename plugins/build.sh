#!/bin/bash

cd plugins
tsc -t ESNEXT --removeComments *.ts &> /dev/null
exit 0
