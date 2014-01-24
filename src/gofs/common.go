package gofs

import "time"

type FileDescriptor uint
type FileIndex uint

type FileDescriptorTable map[FileDescriptor]*DescriptorInfo

type FileInfo struct {
  inode *Inode
}

type DescriptorInfo struct {
  accessFlags AccessFlag
  pointer uint
  fileInfo *FileInfo
}

type ProcState struct {
  fdTable FileDescriptorTable
  lastFd FileDescriptor
}

type Inode struct {
  data []byte
  size uint

  refCount uint // do we need this?
  fileType uint // something? generics would be great here

  perms uint
  ownerId uint
  groupId uint

  lastModTime time.Time
  lastAccessTime time.Time
  createTime time.Time
}

type FileMode uint
const (
  M_EXEC FileMode = 1 << iota
  M_WRITE 
  M_READ
)

type Whence uint
const (
  SEEK_SET Whence = 1 << iota
  SEEK_CUR
  SEEK_END
)

type AccessFlag uint
const (
  O_RDONLY AccessFlag = 1 << iota
  O_WRONLY
  O_RDWR
  O_NONBLOCK
  O_APPEND
  O_CREAT
  O_TRUNC
  O_EXCL
  O_SHLOCK
  O_EXLOCK
  O_NOFOLLOW
  O_SYMLINK
  O_EVTONLY
  O_CLOEXEC
)

