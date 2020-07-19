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

type  ofile struct {
	Length uint64
	Date   uint64
	Data   []byte	
}

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

func ingestoriginalfile(f *os.File, sector int64)(fe ofile, err error){
        var fp oxfsgo.OBFS_FileHeader

        _,err = f.Seek(((sector/29)-1)*1024,0)
        if err == nil {
                err = binary.Read(f, binary.LittleEndian, &fp)
        }
        if err == nil {
		fe.Date=uint64(fp.Date)
		fe.Length=uint64((fp.Aleng*1024)+fp.Bleng-oxfsgo.OBFS_HeaderSize)
	}
	return fe, err
}

func ingestoriginaldir(f *os.File, sector int64, files map[string]ofile) (outfiles map[string]ofile, err error){
	var dp oxfsgo.OBFS_DirPage	
        

        _,err = f.Seek(((sector/29)-1)*1024,0)
        if err == nil {
                err = binary.Read(f, binary.LittleEndian, &dp)
        }
        if err == nil {
	        if dp.P0 != 0{
			files, err = ingestoriginaldir(f,int64(dp.P0),files)
		}
		for i:=int32(0);i<dp.M;i++{
			fmt.Println("file",string(dp.E[i].Name[:]))
			files[string(dp.E[i].Name[:])],err=ingestoriginalfile(f,int64(dp.E[i].Adr))
			if dp.E[i].P != 0 {	
				files, err = ingestoriginaldir(f,int64(dp.E[i].P),files)
			}
		}
	}
	return files, err
}

func ingestextendeddir(f *os.File, sector int64, files map[string]ofile) (outfiles map[string]ofile, err error){
        var dp *oxfsgo.OXFS_DirPage


        _,err = f.Seek((sector/29)-1,0)
        if err == nil {
                binary.Read(f, binary.LittleEndian, &dp)
        }
        if err == nil {
                fmt.Println("ingesting dirpage",sector,"mark",dp.Mark,"count",dp.M)

        }
        return files, err
}



func ingestfs(filename string, origfmt bool)(files map[string]ofile, err error){
	var f *os.File
	var kind  int

	files = make(map[string]ofile)

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
		if origfmt {
			files, err = ingestoriginaldir(f,29,files)
		}else{
                        files, err = ingestextendeddir(f,29,files)
		}
	}
	return files,err
}

func producefs(image string, original, force bool)(err error){
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
			if _,err:=ingestfs(*inPtr,*o2Ptr); err != nil {
		                fmt.Println(err)
				os.Exit(1)
			}else{
				if err=producefs(*outPtr,*x2Ptr,*forcePtr); err != nil {
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
