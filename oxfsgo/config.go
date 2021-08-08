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

// 64-bit oberon extended file system  2 ZiB Max Volume
const OXFS_FnLength    = 47
const OXFS_TabSize     = 4
const OXFS_SectorSize  = 4096
const OXFS_IndexSize   = 512    //SectorSize / 8
const OXFS_HeaderSize  = 64
const OXFS_DirRootAdr  = 29
const OXFS_DirPgSize   = 63
const OXFS_HdrPgSize   = 63
const OXFS_N = 31               //DirPgSize / 2
const OXFS_DirMark    = 0x9B1EA38E
const OXFS_HeaderMark = 0x9BA71D87
const OXFS_FillerSize = 32

type    OBFS_DiskAdr         int32
type    OBFS_FileName       [OBFS_FnLength]byte            // 672 data bytes in zeroth sector of file
type    OBFS_SectorTable    [OBFS_SecTabSize]OBFS_DiskAdr  // 65,184 byte max file size without using extension table
type    OBFS_ExtensionTable [OBFS_ExTabSize]OBFS_DiskAdr   // 3,210,912 max file size with addition of extension table

type    OXFS_DiskAdr         int64
type    OXFS_FileName       [OXFS_FnLength]byte
type    OXFS_SectorTable    [OXFS_TabSize]OXFS_DiskAdr     // 3 point to a file page, last points to a sector of pointers to file pages, except the last 
                                                           // which points to a sector of pointers to pointers to file pages, except the last
                                                           // which points to a sector of pointers to pointers to pointers to file pages, and so on

type  OBFS_FileHeader struct { // first page of each file on disk
        Mark uint32
        Name [32]byte
        Aleng, Bleng, Date int32
        Ext  OBFS_ExtensionTable
        Sec OBFS_SectorTable
        Fill [OBFS_SectorSize - OBFS_HeaderSize]byte
}

type  OXFS_FileHeader struct { // 63 in a HeaderPage
        Type uint32
	Perm uint32
	Date uint64
	Length uint64
	Owner uint32
	Group uint32
        Tab OXFS_SectorTable
}

type  OXFS_HeaderPage struct {
        Mark uint64
	Next uint64        // (*sec no of next non-full HeaderPage, zero if full *)
	Bmap uint64        // bitmap of used FileHeader entries in this page
	Fill [40]byte
	Headers [OXFS_HdrPgSize]OXFS_FileHeader
}

type    OBFS_DirEntry struct { //  (*B-tree node*)
        Name [32]byte
        Adr  OBFS_DiskAdr  // (*sec no of file header*)
        P    OBFS_DiskAdr  // (*sec no of descendant in directory*)
}

type    OXFS_DirEntry struct { //  (*B-tree node*)
        P    OXFS_DiskAdr  // (*sec no of descendant in directory*)
        Adr  OXFS_DiskAdr  // (*sec no of file header*)
        I    byte          // low six bits of I are the HeaderPage index to be used with the Adr field
        Name [47]byte      // top 2 bits of I are number of subsequent DirEntry slots pre-empted
                           // for a longer filename for this entry for a maximum name length of 303 bytes
}

type    OBFS_DirPage struct {
        Mark  uint32
        M     int32
        P0    OBFS_DiskAdr //  (*sec no of left descendant in directory*)
        Fill  [OBFS_FillerSize]byte
        E  [OBFS_DirPgSize]OBFS_DirEntry
}

type    OXFS_DirPage struct {
        Mark  uint64
        M     int64
	N     int64
        P0    OXFS_DiskAdr //  (*sec no of left descendant in directory*)
        Fill  [OXFS_FillerSize]byte
        E  [OXFS_DirPgSize]OXFS_DirEntry
}



func Config() string {
  return "oxfs config"
}
