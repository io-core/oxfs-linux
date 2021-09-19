#!/bin/bash

if [ ! -f RISC.img ]; then
  unzip RISCimg.zip
fi

echo "oxfstool version " `../oxfstool -V`

#test expanding an unpadded image
if [ -d onp ]; then
  rm -rf onp
fi
mkdir -p onp
../oxfstool -o2f -i io.img -o onp

#echo "expanded not padded image:"
#ls -alh onp

if [ -f io2.img ]; then
  rm -f io2.img
fi
../oxfstool -f2o -i onp -o io2.img -s 8M

#test creating an unpadded image
if [ -d onp2 ]; then
  rm -rf onp2
fi
mkdir -p onp2
../oxfstool -o2f -i io2.img -o onp2

#round-trip test
diff -r onp onp2

#test expanding a padded image
if [ -d op ]; then
  rm -rf op
fi
mkdir -p op
../oxfstool -o2f -i RISC.img -o op

#echo "expanded padded image:"
#ls -alh op

echo "tool directory:"
ls -alh



