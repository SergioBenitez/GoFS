package gofs

import (
  "io"
)

type FileDescriptor int64

type File interface {
  io.Reader // Read(p []byte) (n int, err error)
  io.Writer // Write(p []byte) (n int, err error)
  io.Closer // Close() error
  io.Seeker // Seek(offset int64, whence int) (int64, error)
}

type Directory map[string]interface{}

// This is not like Unix's FileTable that is global. This FileTable is per
// process. The *File is what's global. Basically, replaces the indexed layer
// of indirection through the file pointer.
type FileTable map[FileDescriptor]interface{File}

type ProcState struct {
  fileTable FileTable
  cwd Directory
  lastFd FileDescriptor
}

type GlobalState struct {
  root Directory
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

type Whence int
const (
  SEEK_SET = 0
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
