# SPDX-License-Identifier: GPL-2.0-only
#
# Makefile for the Linux minix filesystem routines.
#

obj-m := minix.o

minix-objs := bitmap.o itree_v1.o itree_v2.o namei.o inode.o file.o dir.o

all: ko 

ko:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean

#obj-$(CONFIG_MINIX_FS) += minix.o
#
#minix-objs := bitmap.o itree_v1.o itree_v2.o namei.o inode.o file.o dir.o
