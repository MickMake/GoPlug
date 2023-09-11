#!/bin/bash

echo "########################################"
echo "# Building example plugin1"
pushd plugin1
go build -buildmode=plugin -o plugin1.so -gcflags 'all=-N -l' plugin1.go
popd

echo "########################################"
echo "# Building example plugin2"
pushd plugin2
go build -buildmode=plugin -o plugin2.so -gcflags 'all=-N -l' plugin2.go
# This will force a rebuild during tests.
touch farter2.go
go build -buildmode=plugin -o failplugin.so -gcflags 'all=-N -l' failplugin.go
popd

echo "########################################"
echo "# Building example plugin3"
go build -buildmode=plugin -o plugin3.so -gcflags 'all=-N -l' plugin3.go

