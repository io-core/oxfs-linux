# oxfstool

A tool for creating and unpacking original and extended format oberon disk images
(functional for original format disk images)

Latest binary for Ubuntu: [oxfstool-on-Linux](https://github.com/io-core/oxfs-linux/releases/download/v0.1.2/oxfstool-on-Linux) 
Latest binary for macOS: [oxfstool-on-macOS](https://github.com/io-core/oxfs-linux/releases/download/v0.1.2/oxfstool-on-macOS) 

Usage of oxfstool:

specify one of o2x, x2o, o2f, x2f, f2o, f2x, or check

  -V	output the version of oxfstool
  -check
    	check a disk image
  -f	overwrite output disk image if it exists
  -f2o
    	convert from local files to original format
  -f2x
    	convert from local files to extended format
  -i string
    	input disk image
  -o string
    	output disk image
  -o2f
    	convert from original to local files
  -o2x
    	convert from original to extended format
  -s string
    	output disk image size e.g. '64M', '1G', '8G', etc. or 'same' (default "same")
  -x2f
    	convert from extended to local files
  -x2o
    	convert from extended to original format

# oxfs-linux
Oberon extended file system Linux kernel module
(not at all functional)

