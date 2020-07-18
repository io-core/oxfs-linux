package main

import (
  "fmt"
  "flag"
  "os"
  "encoding/binary"
  "github.com/io-core/oxfs-linux/oxfsgo"
)

const	UNKNOWN = 0
const   ORIGINAL =  1
const   PADDEDORIGINAL =  2
const   EXTENDED =  3
const   PADDEDEXTENDED =  4

func identify(f *os.File) (kind int, size int64, err error) {
	fi, err := f.Stat()
	if err == nil {
		size = fi.Size()
		_,err = f.Seek(0,0)
	}
        if err == nil {
		buf := make([]byte, 4)
		_, err = f.Read(buf)
		if binary.LittleEndian.Uint32(buf) == oxfsgo.OBFS_DirMark {
			kind = ORIGINAL
		}
                if binary.LittleEndian.Uint32(buf) == oxfsgo.OXFS_DirMark {
                        kind = EXTENDED
		}
	}
	return kind, size, err
}

func ingest(filename string, origfmt bool)(item string,err error){
	var f *os.File
	var kind int

	_, err = os.Stat(filename)
	if err == nil {
	        f, err = os.Open(filename)
        }
	if err == nil{
		defer f.Close()
		kind,_,err = identify(f)
	}
	if err == nil{		
		if !(((kind == ORIGINAL) && origfmt ) || ((kind == EXTENDED) && (! origfmt) )){
			err = fmt.Errorf("wrong format for input disk image %s",filename)
		}
	}
	return "OK",err
}

func produce(image string, original, force bool)(err error){
        err = fmt.Errorf("produce function not implemented")

        return err
}

func main() {

        inPtr := flag.String("i", "", "input disk image")
        outPtr := flag.String("o", "", "output disk image")
	sizePtr := flag.String("s", "same", "output disk image size e.g. '64M', '1G', '8G', etc. or 'same'") 
	forcePtr := flag.Bool("f", false, "overwrite output disk image if it exists")	
	o2Ptr := flag.Bool("o2x", false, "convert from original to extended format")	
        x2Ptr := flag.Bool("x2o", false, "convert from extended to original format")
        checkPtr := flag.Bool("check", false, "check a disk image")

	flag.Parse()

	if ((*o2Ptr && (! *x2Ptr)) || (*x2Ptr && (! *o2Ptr))) && (! *checkPtr) {
		if (*inPtr == "") || (*outPtr == ""){
                        fmt.Println("input and output disk images must be specified")
			flag.PrintDefaults()
                        os.Exit(1)
		}else{
			if *o2Ptr {
	        		fmt.Println("converting original format file system",*inPtr,"to extended format file system",*outPtr,"target size",*sizePtr)
			}else{
                                fmt.Println("converting extended format file system",*inPtr,"to original format file system",*outPtr,"target size",*sizePtr)
			}
			if _,err:=ingest(*inPtr,*o2Ptr); err != nil {
		                fmt.Println(err)
				os.Exit(1)
			}else{
				if err=produce(*outPtr,*x2Ptr,*forcePtr); err != nil {
                                	fmt.Println(err)
	                                os.Exit(1)
				}
                        }
		}
	}else if (! *o2Ptr) && (! *x2Ptr) && *checkPtr {
                if (*inPtr == "") || (*outPtr == ""){   
                        fmt.Println("input disk image must be specified")
                        flag.PrintDefaults()
                        os.Exit(1)
                }else{
	                fmt.Println("Checking:", *inPtr)
		}
	}else{
                fmt.Println("specify one of o2x, x2o, or check")
                flag.PrintDefaults()
                os.Exit(1)
	}

}
