package gofs

import (
  "fmt"
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
  data []byte
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

  if file.seek > int64(len(file.data)) {
    return 0, errors.New("EOF")
  }

  fmt.Println("Reading from", file.seek, "to", file.seek + int64(len(p)))
  fmt.Println("File size is:", len(file.data))

  read := copy(p, file.data[file.seek:])
  file.seek += int64(read)
  file.lastAccessTime = time.Now()
  return read, nil
}

func (file *DataFile) Write(p []byte) (n int, err error) {
  if err := file.checkAccess(Read); err != nil { return 0, err }

  needed := file.seek + int64(len(p))
  fmt.Println("Needed:", needed, "Have:", cap(file.data))

  if needed - int64(cap(file.data)) > 0 {
    newData := make([]byte, needed, needed * 2)
    copy(newData, file.data)
    file.data = newData
  }

  file.data = file.data[:needed]
  written := copy(file.data[file.seek:], p)

  file.seek += int64(written)
  file.lastAccessTime = time.Now()
  file.lastModTime = time.Now()

  fmt.Println("Wrote:", written, "File size:", len(file.data), "\n")
  return written, nil
}

func (file *DataFile) Open() error {
  file.status = Open
  return nil
}

func (file *DataFile) Close() error {
  // should do more to ensure that the handle isn't reused...
  // perhaps setting status flags 'opened/closed' is good enough? 
  file.seek = 0
  file.lastAccessTime = time.Now()
  file.status = Closed
  return nil
}

func (file *DataFile) Seek(offset int64, whence int) (int64, error) {
  if err := file.checkAccess(Read); err != nil { return 0, err }

  switch whence {
    case SEEK_SET:
      file.seek = offset;
    case SEEK_CUR:
      file.seek += offset;
    case SEEK_END:
      file.seek = int64(len(file.data)) + offset;
  }

  return file.seek, nil
}

func initDataFile() *DataFile {
  return &DataFile{
    data: make([]byte, 0, 4096),
    lastModTime: time.Now(),
    lastAccessTime: time.Now(),
    createTime: time.Now(),
  }
}
