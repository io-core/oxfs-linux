package main

import (
  "fmt"
  "flag"
  "os"
  "sort"
  "strings"
  "encoding/binary"
//  "path"
  "path/filepath"
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

const PADOFFSET = (524288*512)+1024

type  ofile struct {
	Length uint64
	Date   uint64
	Data   []byte	
}

type iblock struct {
	A	int64
	E	[256] uint32
}

type dirTree struct {
	P0 *dirTree
	Name []string
	P []*dirTree
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
		if kind == UNKNOWN {
                	_,err = f.Seek(PADOFFSET,0)
		}
	        if err == nil {
	                _, err = f.Read(buf)
                }
                if err == nil {
	                if binary.LittleEndian.Uint32(buf) == oxfsgo.OBFS_DirMark {
	                        kind = PADDEDORIGINAL
				size = size - (PADOFFSET)
	                }
	                if binary.LittleEndian.Uint32(buf) == oxfsgo.OXFS_DirMark {
	                        kind = PADDEDEXTENDED
                                size = size - (PADOFFSET)
	                }
		}
	}
	return kind, size, err
}

func populateDir( fileSet []string, files map[string]ofile, n int) * dirTree{
	var dT dirTree	
	c:=len(fileSet)
	if c>n {
		sz:=c/(n+1)
		fmt.Println("sz:",sz)
		dT.Name = make([]string, n)
		dT.P = make([]*dirTree, n)
                
		dT.P0=populateDir(fileSet[0:sz],files,n)
		for i:=1;i<=n;i++{
			e:=i*sz
			dT.Name[i-1]=fileSet[e]
			z:=(i+1)*sz
			if i==n {
				z=c
			}
			dT.P[i-1]=populateDir(fileSet[e+1:z],files,n)
		}
	}else{
	        dT.Name = make([]string, c)
	        dT.P = make([]*dirTree, c)
                for e:=0;e<c;e++{
			dT.Name[e]=fileSet[e]
		}
	}	


	return &dT
}

func produceFile(f *os.File, e ofile, name string, outfmt int, thisSector int)(usedSector int, nextFree int, err error){
	nextFree = thisSector + 29
        if nextFree < 65*29 {
                nextFree = 65*29
        }
        if outfmt == ORIGINAL || outfmt == PADDEDORIGINAL {
                var hdrPage oxfsgo.OBFS_FileHeader
                hdrPage.Mark = oxfsgo.OBFS_HeaderMark
                hdrPage.Aleng = int32(len(e.Data)+oxfsgo.OBFS_HeaderSize)/1024
                hdrPage.Bleng = int32(len(e.Data)+oxfsgo.OBFS_HeaderSize)%1024
                for x, ch := range []byte(name){
                        hdrPage.Name[x]=ch
                }
		if len(e.Data)>=1024-oxfsgo.OBFS_HeaderSize {
			copy(hdrPage.Fill[0:1024-oxfsgo.OBFS_HeaderSize],e.Data[0:1024-oxfsgo.OBFS_HeaderSize])
		}else{
                        copy(hdrPage.Fill[0:len(e.Data)],e.Data[:])
		}
		hdrPage.Sec[0]=oxfsgo.OBFS_DiskAdr(thisSector)
	
	
                if outfmt == ORIGINAL {
                        _,err = f.Seek( (int64(thisSector/29)-1)*1024,0)
                }else{
                        _,err = f.Seek( PADOFFSET + ((int64(thisSector/29))*1024)-1,0)
                }
                if err == nil {
                        fmt.Print(".",thisSector)
                        err = binary.Write(f, binary.LittleEndian, hdrPage)
                }
	}

	return thisSector, nextFree, err
}

func produceDir(f *os.File, dT *dirTree, files map[string]ofile, outfmt int, fw *os.File, thisSector int )( usedSector int, nextFree int, err error){
	var thatSector int
	nextFree = thisSector + 29
	if nextFree < 65*29 {
		nextFree = 65*29
	}

	if outfmt == ORIGINAL || outfmt == PADDEDORIGINAL {
		var dirPage oxfsgo.OBFS_DirPage
		dirPage.Mark = oxfsgo.OBFS_DirMark
                dirPage.M = int32(len(dT.Name))
		if dT.P0 != nil{
			thatSector, nextFree, err = produceDir(f, dT.P0, files, outfmt, fw, nextFree )
			dirPage.P0 = oxfsgo.OBFS_DiskAdr(thatSector)
		}
		for i, _ := range dT.P{
//			fmt.Println(thisSector," -- ",dT.Name[i])
			for x, ch := range []byte(dT.Name[i]){
				dirPage.E[i].Name[x]=ch
			}
                        thatSector, nextFree, err = produceFile(f, files[dT.Name[i]], dT.Name[i], outfmt, nextFree )
                        dirPage.E[i].Adr = oxfsgo.OBFS_DiskAdr(thatSector)
			if dT.P[i] != nil {
				thatSector, nextFree, err = produceDir(f, dT.P[i], files, outfmt, fw, nextFree )
	                        dirPage.E[i].P = oxfsgo.OBFS_DiskAdr(thatSector)
			}
		}
		
	        if outfmt == ORIGINAL {
	                _,err = f.Seek( (int64(thisSector/29)-1)*1024,0)
	        }else{
	                _,err = f.Seek( PADOFFSET + ((int64(thisSector/29))*1024)-1,0)
	        }
	        if err == nil {
			fmt.Print(".",thisSector)
	                err = binary.Write(f, binary.LittleEndian, dirPage)
		}

	}else{
	}


	return thisSector, nextFree, err
}

func installBootImage(f *os.File, bootImage []byte, outfmt int)(err error){
	if outfmt == ORIGINAL || outfmt == EXTENDED{
                _,err = f.Seek( 1024,0)
	}else{
                _,err = f.Seek( PADOFFSET + 1024,0)
	}
	if err != nil {
		_, err = f.Write(bootImage)
	}
	return err
}

func produceDirTree( files map[string]ofile, outfmt int,fw *os.File) (err error){
	var nA []string

	nE:=len(files)
	if _, ok := files["_BOOTIMAGE_"]; ok {
		nE=len(files)-1
		_ = installBootImage(fw,files["_BOOTIMAGE_"].Data,outfmt)
	}
	if err == nil {
		nA = make([]string,nE)

		i:=0
		for fn, _ := range files {
			if fn != "_BOOTIMAGE_" {
				nA[i]=fn
				i++
			}
		}

		rnA := nA[:]
		sort.Strings(rnA)

		dT := populateDir(rnA[:],files,oxfsgo.OBFS_N+(oxfsgo.OBFS_N/2))

		_,_, err = produceDir(fw,dT,files,outfmt,fw,29)
	}
	
	return err
}

func producefs(name string, files map[string]ofile, outfmt int, force bool)(err error){
	var fi *os.File
        keys := make([]string, 0, len(files))
        for k := range files {
                keys = append(keys, k)
        }
        sort.Strings(keys)

	if outfmt != LOCALFILES{

                fi, err = os.Open(name)
		if err != nil{ // assume because file does not exist
			err = nil
                }else{
                        fs, staterr := fi.Stat()
                        fi.Close()
                        switch {
                          case staterr != nil:
                                err = staterr
                          case fs.IsDir():
                                err = fmt.Errorf("destination disk image is a directory")
                          default:
                                err = fmt.Errorf("destination disk image already exists")
                        }
                }
		if err == nil {
			fw, err := os.Create( name )
			if err == nil {
				defer fw.Close()
				err = produceDirTree(files,outfmt,fw)
			}

		}
	}else{
		fi, err = os.Open(name)
	        if err == nil{
			fs, staterr := fi.Stat()
			fi.Close()
			switch {
			  case staterr != nil:
			  	err = staterr
			  case fs.IsDir():
			        for _, k := range keys {
				   if err == nil {
					fname:=name+"/"+k
					fw, err := os.Create( strings.Replace(fname, "\x00", "", -1))
					if err == nil {
						_, err = fw.Write(files[k].Data)
					}
					fw.Close()
				   }
			        }
	                  default:
	                        err = fmt.Errorf("destination for localfiles is not a directory")
	                }
		}
	}

        return err
}

func getOriginalDataBlock(f *os.File, pad int64, e uint64, fp *oxfsgo.OBFS_FileHeader, iblk *iblock)(block []byte, err error){
	block = make([]byte, 1024)
	if e < oxfsgo.OBFS_SecTabSize {
	        _,err = f.Seek(pad + (int64(fp.Sec[e])/29-1)*1024,0)
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
	                _,err = f.Seek(pad + (iblk.A/29-1)*1024,0)
	                if err == nil {
				err = binary.Read(f, binary.LittleEndian, iblk.E)
	                }
		}
                if err == nil {
                        _,err = f.Seek(pad + (int64(iblk.E[r])/29-1)*1024,0)
                }
                if err == nil {
                	_, err = f.Read(block)
                }
	}

	return block, err
}

func ingestOriginalFile(f *os.File, pad int64, sector int64)(fe ofile, err error){
        var fp oxfsgo.OBFS_FileHeader
	var iblk iblock
	var block []byte

	const offset = 1024-oxfsgo.OBFS_HeaderSize
	

        _,err = f.Seek(pad + ((sector/29)-1)*1024,0)
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
			block,err=getOriginalDataBlock(f,pad,e,&fp,&iblk)
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

func ingestOriginalBootImage(f *os.File, pad int64)(fe ofile, err error){
	var sz int
        block := make([]byte, 1024)
        _,err = f.Seek(pad + 1024,0)
        if err == nil {
                _, err = f.Read(block)
	}
        if err == nil {
		sz = int(block[16])+(int(block[17])*0x100)+(int(block[18])*0x10000)+(int(block[19])*0x1000000)
		_,err = f.Seek(pad + 1024,0)
        }
        if err == nil {
		fmt.Println("Boot Image Size:",sz)
		block = make([]byte, sz)
                _, err = f.Read(block)
        }
        if err == nil {

                fe.Date=0
                fe.Length=uint64(sz)
                fe.Data=block
	}
	return fe,err
}


func ingestExtendedBootImage(f *os.File, pad int64)(kernel []byte, err error){

        return nil,err
}


func ingestOriginalDir(f *os.File, pad int64, sector int64, files map[string]ofile) (outfiles map[string]ofile, err error){
	var dp oxfsgo.OBFS_DirPage	
        

        _,err = f.Seek(pad + ((sector/29)-1)*1024,0)
        if err == nil {
                err = binary.Read(f, binary.LittleEndian, &dp)
        }
        if err == nil {
	        if dp.P0 != 0{
			files, err = ingestOriginalDir(f,pad,int64(dp.P0),files)
		}
		for i:=int32(0);i<dp.M;i++{
//			fmt.Println("file",string(dp.E[i].Name[:]))
			files[string(dp.E[i].Name[:])],err=ingestOriginalFile(f,pad,int64(dp.E[i].Adr))
			if dp.E[i].P != 0 {	
				files, err = ingestOriginalDir(f,pad,int64(dp.E[i].P),files)
			}
		}
	}
	return files, err
}

func ingestExtendedDir(f *os.File, pad int64, sector int64, files map[string]ofile) ( _ map[string]ofile, err error){
        var dp *oxfsgo.OXFS_DirPage


        _,err = f.Seek(pad + (sector/29)-1,0)
        if err == nil {
                binary.Read(f, binary.LittleEndian, &dp)
        }
        if err == nil {
                fmt.Println("ingesting dirpage",sector,"mark",dp.Mark,"count",dp.M)

        }
        return files, err
}

func visit(fnames *[]string) filepath.WalkFunc {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
//            log.Fatal(err)
        }
        *fnames = append(*fnames, path)
        return nil
    }
}

func ingestFromFile(filename string)(fi ofile, err error){
        var f *os.File
        _, err = os.Stat(filename)
        if err == nil {
                f, err = os.Open(filename)
        }
        if err == nil {
	        fs, _ := f.Stat()
                fi.Length = uint64(fs.Size())
	        fi.Data = make([]byte, fi.Length)
                _, err = f.Read(fi.Data)
        }
	return fi, err 
}

func ingestFS(filename string, infmt int)(files map[string]ofile, err error){
	var f *os.File
	var kind  int
	var pad int64
	var fnames []string	

	files = make(map[string]ofile)

	_, err = os.Stat(filename)
	if err == nil {
	        f, err = os.Open(filename)
        }
        if err == nil{
                defer f.Close()
                fs, _ := f.Stat()
                        switch mode := fs.Mode(); {
                          case mode.IsDir():
				pre:=len(filename)
				err = filepath.Walk(filename, visit(&fnames))
				for _, fn := range fnames {
					//_, file := path.Split(fn)
					if len(fn)>pre+1{
						fmt.Println(fn[pre+1:])
						files[fn[pre+1:]],err=ingestFromFile(fn)
					}
				}
                                //                err = fmt.Errorf("don't know how to read directory %s",filename)
 
                          case mode.IsRegular():

				if err == nil{
					kind,_,err = identify(f)
					if kind == PADDEDORIGINAL || kind == PADDEDEXTENDED {
						pad = PADOFFSET
					}
				}
				if err == nil{		
					if !(((kind == ORIGINAL) && (infmt==ORIGINAL) ) || 
					     ((kind == EXTENDED) && (infmt == EXTENDED) ) ||  
					     ((kind == PADDEDORIGINAL) && (infmt==ORIGINAL) ) || 
			                     ((kind == PADDEDEXTENDED) && (infmt == EXTENDED) )   ){
						err = fmt.Errorf("wrong format for input disk image %s",filename)
					}
			        }

	

			        if err == nil{
			                if infmt == ORIGINAL {
						files["_BOOTIMAGE_"],err=ingestOriginalBootImage(f,pad)
						
			                }else{
			                        _,err=ingestExtendedBootImage(f,pad)
					}
			        }
			        if err == nil{
			 		if infmt == ORIGINAL {
						files, err = ingestOriginalDir(f,pad,29,files)
					}else{
			                        files, err = ingestExtendedDir(f,pad,29,files)
					}
				}
			}
	}
	return files,err
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
	}else if (*o2xPtr || *x2oPtr || *o2fPtr || *x2fPtr || *f2oPtr || *f2xPtr) {
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
