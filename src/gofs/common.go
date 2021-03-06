package gofs

import (
  "io"
  "gofs/dstore"
  "time"
)

type Directory map[string]interface{}

// This is the per process FileDescriptor Table
type FileDescriptor int16
type FileDescriptorTable map[FileDescriptor]interface{File}

// Global structure keeps track of all open files via an array of *File objects.
// These *File objects contain the necessary file information.
// type FileTable []interface{File}

type File interface {
  io.Closer // Close() error
  io.Reader // Read(p []byte) (n int, err error)
  io.Writer // Write(p []byte) (n int, err error)
  io.Seeker // Seek(offset int64, whence int) (int64, error)
}

type Inode struct {
  data interface{dstore.DataStore}

  perms uint
  ownerId uint
  groupId uint

  lastModTime time.Time
  lastAccessTime time.Time
  createTime time.Time

  linkCount int
  fileCount int
}

const MAX_DESCRIPTORS = 1024;
type ProcState struct {
  fileDescriptorTable FileDescriptorTable
  freeDescriptors [MAX_DESCRIPTORS]FileDescriptor
  lastFd FileDescriptor
  cwd Directory
}

type GlobalState struct {
  root Directory
  // fileTable FileTable
  stdIn interface{File}
  stdOut interface{File}
  stdErr interface{File}
}

type FileMode uint
const (
  M_EXEC FileMode = 1 << iota
  M_WRITE
  M_READ
)

const (
  SEEK_SET = iota
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
