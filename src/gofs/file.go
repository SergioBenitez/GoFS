package gofs

import (
  // "fmt"
  "gofs/dstore"
  "time"
  "errors"
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

type FileStatus uint
const (
  Closed FileStatus = iota
  Open
)

type FileAccess uint
const (
  Read FileAccess = iota
  Write
  Seek
)

type DataFile struct {
  data interface{dstore.DataStore}
  seek int64

  perms uint
  ownerId uint
  groupId uint

  lastModTime time.Time
  lastAccessTime time.Time
  createTime time.Time

  status FileStatus
}

func (file *DataFile) checkAccess(acc FileAccess) error {
  switch file.status {
    case Closed:
      return errors.New("File is closed.")
  }
  return nil
}

func (file *DataFile) Read(p []byte) (int, error) {
  if err := file.checkAccess(Read); err != nil { return 0, err }

  if file.seek > file.data.Size() {
    return 0, errors.New("EOF")
  }

  read, err := file.data.Read(file.seek, p)
  file.lastAccessTime = time.Now()
  file.seek += int64(read)
  return read, err
}

func (file *DataFile) Write(p []byte) (int, error) {
  if err := file.checkAccess(Write); err != nil { return 0, err }
  
  wrote, err := file.data.Write(file.seek, p)
  file.lastAccessTime = time.Now()
  file.lastModTime = time.Now()
  file.seek += int64(wrote)

  return wrote, err
}

func (file *DataFile) Open() error {
  file.status = Open
  return nil
}

func (file *DataFile) Close() error {
  file.seek = 0
  file.lastAccessTime = time.Now()
  file.status = Closed
  return nil
}

func (file *DataFile) Seek(offset int64, whence int) (int64, error) {
  if err := file.checkAccess(Seek); err != nil { return 0, err }

  switch whence {
    case SEEK_SET:
      file.seek = offset;
    case SEEK_CUR:
      file.seek += offset;
    case SEEK_END:
      file.seek = file.data.Size() + offset;
  }

  return file.seek, nil
}

func initDataFile() *DataFile {
  return &DataFile{
    data: dstore.InitArrayStore(4096),
    lastModTime: time.Now(),
    lastAccessTime: time.Now(),
    createTime: time.Now(),
  }
}
