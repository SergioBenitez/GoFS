package gofs

import (
  "time"
  "fmt"
)

/**
* Ideally, a directory would be a 'File' of type 'Directory' that you can read
* from except it'd just be a map from path (string) to inode (File).
*
* Open("...") would grab a File and put a pointer to it in a table and return
* the index (or whatever) of that table. The rest of the file system functions
* would simply look up the File and call methods on it.
*
* class File ...
* instance RegularFile : File ... // not the right syntax
* readFile :: (File a) => a -> ByteString
* lookup :: (File a) => FileDescriptor -> a
* read :: FileDescriptor -> ByteString
* read = readFile . lookup
*
* Every Open() call creates a new entry in the file table. That is, besides the
* underlying file contents, two different open calls for the same file share
* nothing (file pointer, permissions, etc). However, by sharing the file
* descriptor, two processes can modify the same entry in the file table.
*/

type DataFile struct {
  data []byte
  size uint

  perms uint
  ownerId uint
  groupId uint

  lastModTime time.Time
  lastAccessTime time.Time
  createTime time.Time
}

func (file *DataFile) Read(p []byte) (n int, err error) {
  fmt.Println("Should read.")
  n, err = 0, nil
  return
}

func (file *DataFile) Write(p []byte) (n int, err error) {
  fmt.Println("Should write.")
  n, err = 0, nil
  return
}

func (file *DataFile) Close() error {
  fmt.Println("Should close.")
  return nil
}

func (file *DataFile) Seek(offset int64, whence int) (int64, error) {
  fmt.Println("Should seek.")
  return 0, nil
}

func initDataFile() *DataFile {
  return &DataFile{}
}
