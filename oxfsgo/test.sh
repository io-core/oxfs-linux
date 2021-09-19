#!/bin/bash

unzip RISCimg.zip

echo "oxfstool version " `./oxfstool -V`

#test expanding an unpadded image
mkdir -p onp
./oxfstool -o2f -i io.img -o onp

#echo "expanded not padded image:"
#ls -alh onp

./oxfstool -f2o -i onp -o io2.img -s 8M

#test creating an unpadded image
mkdir -p onp2
./oxfstool -o2f -i io2.img -o onp2

#round-trip test
diff -r onp onp2

#test expanding a padded image
mkdir -p op
./oxfstool -o2f -i RISC.img -o op

#echo "expanded padded image:"
#ls -alh op

echo "tool directory:"
ls -alh


