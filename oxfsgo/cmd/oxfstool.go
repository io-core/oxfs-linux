package main

import (
  "fmt"
  "flag"
  "os"
  "sort"
  "encoding/binary"
  "github.com/io-core/oxfs-linux/oxfsgo"
)

const (
	UNKNOWN = iota
	ORIGINAL 
	PADDEDORIGINAL 
	EXTENDED 
	PADDEDEXTENDED 
	LOCALFILES 
)

type  ofile struct {
	Length uint64
	Date   uint64
	Data   []byte	
}

type iblock struct {
	A	int64
	E	[256] uint32
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


func getOriginalDataBlock(f *os.File, e uint64, fp *oxfsgo.OBFS_FileHeader, iblk *iblock)(block []byte, err error){
	block = make([]byte, 1024)
	if e < oxfsgo.OBFS_SecTabSize {
	        _,err = f.Seek((int64(fp.Sec[e])/29-1)*1024,0)
	        if err == nil {
	                _, err = f.Read(block)	
		}
	}else{
		x:=e-oxfsgo.OBFS_SecTabSize
		i:=int64(fp.Ext[x/256])
		r:=x%256
		if iblk.A == i {
//			fmt.Print("!")
		}else{
//			fmt.Print("?")
			iblk.A=i
	                _,err = f.Seek((iblk.A/29-1)*1024,0)
	                if err == nil {
				err = binary.Read(f, binary.LittleEndian, iblk.E)
	                }
		}
                if err == nil {
                        _,err = f.Seek((int64(iblk.E[r])/29-1)*1024,0)
                }
                if err == nil {
                	_, err = f.Read(block)
                }
	}

	return block, err
}

func ingestOriginalFile(f *os.File, sector int64)(fe ofile, err error){
        var fp oxfsgo.OBFS_FileHeader
	var iblk iblock
	var block []byte

	const offset = 1024-oxfsgo.OBFS_HeaderSize
	

        _,err = f.Seek(((sector/29)-1)*1024,0)
        if err == nil {
                err = binary.Read(f, binary.LittleEndian, &fp)
        }
        if err == nil {
		fe.Date=uint64(fp.Date)
		fe.Length=uint64((fp.Aleng*1024)+fp.Bleng-oxfsgo.OBFS_HeaderSize)
		fe.Data=make([]byte, fe.Length)
		sz:=uint64(offset)
		if sz > fe.Length {
			sz =  fe.Length
		}
		copy(fe.Data[0:sz],fp.Fill[:])
		e:=uint64(1)
		for i:=uint64(offset); i < fe.Length; i=i+1024 {
			block,err=getOriginalDataBlock(f,e,&fp,&iblk)
			if e*1024+offset<=fe.Length{
				copy(fe.Data[i:i+1024],block)
			}else{
				sz =  fe.Length-i
				copy(fe.Data[i:i+sz],block[0:sz])
			}
			e++
		}
//		fmt.Println()
	}
	return fe, err
}

func ingestOriginalBootImage(f *os.File)(kernel []byte, err error){

	return nil,err
}


func ingestExtendedBootImage(f *os.File)(kernel []byte, err error){

        return nil,err
}


func ingestOriginalDir(f *os.File, sector int64, files map[string]ofile) (outfiles map[string]ofile, err error){
	var dp oxfsgo.OBFS_DirPage	
        

        _,err = f.Seek(((sector/29)-1)*1024,0)
        if err == nil {
                err = binary.Read(f, binary.LittleEndian, &dp)
        }
        if err == nil {
	        if dp.P0 != 0{
			files, err = ingestOriginalDir(f,int64(dp.P0),files)
		}
		for i:=int32(0);i<dp.M;i++{
//			fmt.Println("file",string(dp.E[i].Name[:]))
			files[string(dp.E[i].Name[:])],err=ingestOriginalFile(f,int64(dp.E[i].Adr))
			if dp.E[i].P != 0 {	
				files, err = ingestOriginalDir(f,int64(dp.E[i].P),files)
			}
		}
	}
	return files, err
}

func ingestExtendedDir(f *os.File, sector int64, files map[string]ofile) (outfiles map[string]ofile, err error){
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



func ingestFS(filename string, infmt int)(files map[string]ofile, err error){
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
		if !(((kind == ORIGINAL) && (infmt==ORIGINAL) ) || ((kind == EXTENDED) && (infmt == EXTENDED) )){
			err = fmt.Errorf("wrong format for input disk image %s",filename)
		}
        }
        if err == nil{
                if infmt == ORIGINAL {
			_,err=ingestOriginalBootImage(f)
                }else{
                        _,err=ingestExtendedBootImage(f)
		}
        }
        if err == nil{
 		if infmt == ORIGINAL {
			files, err = ingestOriginalDir(f,29,files)
		}else{
                        files, err = ingestExtendedDir(f,29,files)
		}
	}
	return files,err
}

func producefs(name string, files map[string]ofile, outfmt int, force bool)(err error){
	var fi *os.File
        keys := make([]string, 0, len(files))
        for k := range files {
                keys = append(keys, k)
        }
        sort.Strings(keys)

	if outfmt == ORIGINAL{
		err = fmt.Errorf("produce ORIGINAL not implemented")
	}else if outfmt == EXTENDED{
		err = fmt.Errorf("produce EXTENDED not implemented")
	}else if outfmt == LOCALFILES{
		fi, err = os.Open(name)
	        if err == nil{
			fs, staterr := fi.Stat()
			fi.Close()
			switch {
			  case staterr != nil:
			  	err = staterr
			  case fs.IsDir():
			        for _, k := range keys {
			                fmt.Println(k, files[k].Date,files[k].Length)
//      	        		fmt.Println(string(files[k].Data))
			        }
	                  default:
	                        err = fmt.Errorf("destination for localfiles is not a directory")
	                }
		}
	}

        return err
}

func main() {

        inPtr := flag.String("i", "", "input disk image")
        outPtr := flag.String("o", "", "output disk image")
//	sizePtr := flag.String("s", "same", "output disk image size e.g. '64M', '1G', '8G', etc. or 'same'") 
	forcePtr := flag.Bool("f", false, "overwrite output disk image if it exists")	
	o2xPtr := flag.Bool("o2x", false, "convert from original to extended format")	
        x2oPtr := flag.Bool("x2o", false, "convert from extended to original format")
        o2fPtr := flag.Bool("o2f", false, "convert from original to local files")
        x2fPtr := flag.Bool("x2f", false, "convert from extended to local files")
        f2xPtr := flag.Bool("f2x", false, "convert from local files to extended format")
        f2oPtr := flag.Bool("f2o", false, "convert from local files to original format")
        checkPtr := flag.Bool("check", false, "check a disk image")

	flag.Parse()

	infmt := UNKNOWN
	outfmt := UNKNOWN
	c:=0
	if *o2xPtr { c++; infmt = ORIGINAL; outfmt = EXTENDED; }
        if *x2oPtr { c++; infmt = EXTENDED; outfmt = ORIGINAL; }
        if *o2fPtr { c++; infmt = ORIGINAL; outfmt = LOCALFILES; }
        if *x2fPtr { c++; infmt = EXTENDED; outfmt = LOCALFILES; }
        if *f2oPtr { c++; infmt = LOCALFILES; outfmt = ORIGINAL; }
        if *f2xPtr { c++; infmt = LOCALFILES; outfmt = EXTENDED; }
        if *checkPtr { c++; }

	if c != 1 {
                fmt.Println("specify one of o2x, x2o, o2f, x2f, f2o, f2x, or check")
                flag.PrintDefaults()
                os.Exit(1)
	}else if (*o2xPtr || *x2oPtr || *o2fPtr || *x2fPtr ) {
		if (*inPtr == "") || (*outPtr == ""){
                        fmt.Println("input and output image or location must be specified")
			flag.PrintDefaults()
                        os.Exit(1)
		}else{
			if files,err:=ingestFS(*inPtr, infmt); err != nil {
		                fmt.Println(err)
				os.Exit(1)
			}else{
				if err=producefs(*outPtr, files, outfmt, *forcePtr); err != nil {
                                	fmt.Println(err)
	                                os.Exit(1)
				}
                        }
		}
	}else if *checkPtr {
                if (*inPtr == "") || (*outPtr == ""){   
                        fmt.Println("input disk image must be specified")
                        flag.PrintDefaults()
                        os.Exit(1)
                }else{
	                fmt.Println("Checking:", *inPtr)
		}
	}

}
