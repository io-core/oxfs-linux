package main

import (
  "fmt"
  "flag"
  "os"

//  "github.com/io-core/oxfs-linux/oxfsgo"
)



func ingest(filename string, origfmt bool)(item string,err error){
	var f *os.File

	_, err = os.Stat(filename)
	if err == nil {
	        f, err = os.Open(filename)
        }
	if err == nil{
		defer f.Close()
	        fmt.Println("opened",filename)

		
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
