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
if [[ $rv == 1 ]]  
then    
    echo "failed encoding"    
    exit 1
fi

#echo "expanded not padded image:"
#ls -alh onp

if [ -f io2.img ]; then
  rm -f io2.img
fi
../oxfstool -f2o -i onp -o io2.img -s 8M
if [[ $rv == 1 ]]  
then    
    echo "failed decoding"    
    exit 1
fi

#test creating an unpadded image
if [ -d onp2 ]; then
  rm -rf onp2
fi
mkdir -p onp2
../oxfstool -o2f -i io2.img -o onp2
if [[ $rv == 1 ]]  
then    
    echo "failed encoding"    
    exit 1
fi

#uncomment next line to force a failure:
#echo "failure!" >> onp2/failure.file

#round-trip test
diff -r onp onp2
rv=$?  
if [[ $rv == 1 ]]  
then    
    echo "failed comparison"    
    exit 1
fi

#test expanding a padded image
if [ -d op ]; then
  rm -rf op
fi
mkdir -p op
../oxfstool -o2f -i RISC.img -o op
if [[ $rv == 1 ]]  
then    
    echo "failed decoding"    
    exit 1
fi

#echo "expanded padded image:"
#ls -alh op

echo "tool directory:"
ls -alh



