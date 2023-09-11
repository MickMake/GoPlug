#!/bin/bash

echo "########################################"
echo "# Building master example"
pushd master
go build
popd

echo "########################################"
echo "# Building example plugins"
pushd plugins
./make.sh
popd
