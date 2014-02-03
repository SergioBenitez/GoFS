package gofs

import (
  "gofs/dstore"
  "time"
  "errors"
)

/**
* Here's the philosophical QOTD: Is a directory /really/ a file?
*
* My response: No.
*
* Unlike files, directories can't be read read, written, or seeked to in any
* meaninful way, at least not as they're implemented today. Further, directories
* aren't 'opened' or 'closed' like files are; no, they simply exist as part of
* the structure of the file system itself. Directories are 'made' and 'changed
* into', two operations that don't exist for files. Indeed, file and directory
* interface are highly orthogonal: directories are as much files as file
* systems are files: they aren't! So, let's not treat them like files.
*
* So, are devices files? Sometimes. I posit that if the device can be opened,
* closed, read, written, and have only a minimal additional interface, then we
* can call them files. Thankfully, most devices fit into this umbrella: shared
* memory, networking devices, including sockets, the console, and more. 
*
* Older notes:
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
  seek int

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
  if file.seek >= file.data.Size() { return 0, errors.New("EOF") }

  read, err := file.data.Read(file.seek, p)
  file.lastAccessTime = time.Now()
  file.seek += read
  return read, err
}

func (file *DataFile) Write(p []byte) (int, error) {
  if err := file.checkAccess(Write); err != nil { return 0, err }
  
  wrote, err := file.data.Write(file.seek, p)
  file.lastAccessTime = time.Now()
  file.lastModTime = time.Now()
  file.seek += wrote

  return wrote, err
}

func (file *DataFile) Open() error {
  file.seek = 0
  file.lastAccessTime = time.Now()
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
      file.seek = int(offset);
    case SEEK_CUR:
      file.seek += int(offset);
    case SEEK_END:
      file.seek = file.data.Size() + int(offset);
  }

  return int64(file.seek), nil
}

func initDataFile() *DataFile {
  return &DataFile{
    // data: dstore.InitHashStore(4096),
    data: dstore.InitArrayStore(4096),
    lastModTime: time.Now(),
    lastAccessTime: time.Now(),
    createTime: time.Now(),
  }
}
