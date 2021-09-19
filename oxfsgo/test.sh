#!/bin/bash

unzip RISCimg.zip

echo "oxfstool version " `./oxfstool -V`

echo "tool directory:"
ls -alh

mkdir -p onp
./oxfstool -o2f -i io.img -o onp

echo "expanded not padded image:"
ls -alh onp

mkdir -p op
./oxfstool -o2f -i RISC.img -o op

echo "expanded padded image:"
ls -alh op

