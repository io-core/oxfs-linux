package oxfsgo

// 32-bit oberon file system  141.2 GiB Max Volume         
const OBFS_FnLength    = 32                                  
const OBFS_SecTabSize  = 64                                  
const OBFS_ExTabSize   = 12  //64+12*256 = 3MB max file size 
const OBFS_SectorSize  = 1024                                
const OBFS_IndexSize   = 256    //SectorSize / 4             
const OBFS_HeaderSize  = 352                                 
const OBFS_DirRootAdr  = 29                                  
const OBFS_DirPgSize   = 24                                  
const OBFS_N = 12               //DirPgSize / 2              
const OBFS_DirMark    = 0x9B1EA38D
const OBFS_HeaderMark = 0x9BA71D86
const OBFS_FillerSize = 52                                   

type    OBFS_DiskAdr         int32
type    OBFS_FileName       [OBFS_FnLength]byte            // 672 data bytes in zeroth sector of file
type    OBFS_SectorTable    [OBFS_SecTabSize]OBFS_DiskAdr  // 65,184 byte max file size without using extension table
type    OBFS_ExtensionTable [OBFS_ExTabSize]OBFS_DiskAdr   // 3,210,912 max file size with addition of extension table


type  OBFS_FileHeader struct { // (*first page of each file on disk*)
        Mark uint32
        Name [32]byte
        Aleng, Bleng, Date int32
        Ext  OBFS_ExtensionTable
        Sec OBFS_SectorTable
        fill [OBFS_SectorSize - OBFS_HeaderSize]byte
}

type    OBFS_FileHd *OBFS_FileHeader
type    OBFS_IndexSector [OBFS_IndexSize]OBFS_DiskAdr
type    OBFS_DataSector [OBFS_SectorSize]byte

type    OBFS_DirEntry struct { //  (*B-tree node*)
        Name [32]byte
        Adr  OBFS_DiskAdr  // (*sec no of file header*)
        P    OBFS_DiskAdr  // (*sec no of descendant in directory*)
}

type    OBFS_DirPage struct {
        Mark  uint32
        M     int32
        P0    OBFS_DiskAdr //  (*sec no of left descendant in directory*)
        fill  [OBFS_FillerSize]byte
        E  [OBFS_DirPgSize]OBFS_DirEntry
}

func Config() string {
  return "oxfs config"
}
