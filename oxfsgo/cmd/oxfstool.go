package main

import (
  "fmt"
  "flag"

//  "github.com/io-core/oxfs-linux/oxfsgo"
)



func ingest(image string, original bool)(item string,err error){
	err = fmt.Errorf("ingest function not implemented")

	return "OK",err
}

func produce(image string, original bool)(err error){
        err = fmt.Errorf("produce function not implemented")

        return err
}

func main() {

        inPtr := flag.String("i", "", "input disk image")
        outPtr := flag.String("o", "", "output disk image")
	o2Ptr := flag.Bool("o2x", false, "convert from original to extended format")	
        x2Ptr := flag.Bool("x2o", false, "convert from extended to original format")
        checkPtr := flag.Bool("check", false, "check a disk image")

	flag.Parse()

	if ((*o2Ptr && (! *x2Ptr)) || (*x2Ptr && (! *o2Ptr))) && (! *checkPtr) {
		if (*inPtr == "") || (*outPtr == ""){
                        fmt.Println("input and output disk images must be specified")
			flag.PrintDefaults()
		}else{
			if *o2Ptr {
	        		fmt.Println("converting original format file system",*inPtr,"to extended format file system",*outPtr)
			}else{
                                fmt.Println("converting extended format file system",*inPtr,"to original format file system",*outPtr)
			}
			if _,err:=ingest(*inPtr,*o2Ptr); err != nil {
		                fmt.Println(err)
			}else{
				if err=produce(*outPtr,*x2Ptr); err != nil {
                                	fmt.Println(err)
				}
                        }
		}
	}else if (! *o2Ptr) && (! *x2Ptr) && *checkPtr {
                if (*inPtr == "") || (*outPtr == ""){   
                        fmt.Println("input disk image must be specified")
                        flag.PrintDefaults()
                }else{
	                fmt.Println("Checking:", *inPtr)
		}
	}else{
                fmt.Println("specify one of o2x, x2o, or check")
                flag.PrintDefaults()
	}

}
