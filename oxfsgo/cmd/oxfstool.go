package main

import (
  "fmt"
  "flag"

  "github.com/io-core/oxfs-linux/oxfsgo"
)

func main() {

        inPtr := flag.String("i", "RISC.img", "Disk image to read")
        outPtr := flag.String("o", "OXFS.img", "Disk image to write")
	convPtr := flag.Bool("convert", false, "Convert a disk image")	
        checkPtr := flag.Bool("check", false, "Check a disk image")

	flag.Parse()

	fmt.Println("starting oxfstool")
	fmt.Println("Config:", oxfsgo.Config())
	if *convPtr {
	        fmt.Println("converting")
	        fmt.Println("From:", *inPtr)
	        fmt.Println("To:", *outPtr)
	}else if *checkPtr{
                fmt.Println("Checking:", *inPtr)
	}else{
                fmt.Println("Must specify convert or check")
	}

}
