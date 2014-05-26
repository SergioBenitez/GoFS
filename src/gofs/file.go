package gofs

import (
  // "fmt"
  "errors"
  "gofs/dstore"
  "time"
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
  inode  *Inode
  seek   int
  status FileStatus
}

func (file *DataFile) checkAccess(acc FileAccess) error {
  switch file.status {
  case Closed:
    return errors.New("File is closed.")
  }
  return nil
}

func (file *DataFile) Size() int {
  return file.inode.data.Size()
}

func (file *DataFile) Read(p []byte) (int, error) {
  if err := file.checkAccess(Read); err != nil { return 0, err }
  if file.seek >= file.Size() { return 0, errors.New("EOF") }

  read, err := file.inode.data.Read(file.seek, p)
  file.inode.lastAccessTime = time.Now()
  file.seek += read
  return read, err
}

func (file *DataFile) Write(p []byte) (int, error) {
  if err := file.checkAccess(Write); err != nil { return 0, err }

  wrote, err := file.inode.data.Write(file.seek, p)
  file.inode.lastAccessTime = time.Now()
  file.inode.lastModTime = time.Now()
  file.seek += wrote

  return wrote, err
}

// Open and Close should simply increment and decrement a reference count for
// when file descriptors are shared between processes so that each can Close()
// without affecting the other, and so that when all of them Close(), the handle
// is disgarded.

func (file *DataFile) Open() error {
  file.seek = 0
  file.inode.lastAccessTime = time.Now()
  file.status = Open
  return nil
}

func (file *DataFile) Close() error {
  file.seek = 0
  file.inode.lastAccessTime = time.Now()
  file.status = Closed

  file.inode.decrementFileCount()
  if USE_FILE_ARENA { return ArenaReturnDataFile(file) }
  return nil
}

func (file *DataFile) Seek(offset int64, whence int) (int64, error) {
  if err := file.checkAccess(Seek); err != nil {
    return 0, err
  }

  switch whence {
  case SEEK_SET:
    file.seek = int(offset)
  case SEEK_CUR:
    file.seek += int(offset)
  case SEEK_END:
    file.seek = file.Size() + int(offset)
  }

  return int64(file.seek), nil
}

func initDataFile(inode *Inode) *DataFile {
  inode.incrementFileCount()

  if USE_FILE_ARENA {
    file, err := ArenaAllocateDataFile(inode)
    if err != nil { panic("Out of arena memory!") }
    return file
  }

  return &DataFile{
    inode:  inode,
    seek:   0,
    status: Open,
  }
}

func (inode *Inode) destroyIfNeeded() {
  if inode.linkCount == 0 && inode.fileCount == 0 {
    switch data := inode.data.(type) {
    case *dstore.PageStore:
      data.ReleasePages()
    }
    // fmt.Println("Destroy!")
  }
}

func (inode *Inode) decrementLinkCount() {
  // fmt.Println("||||| -- Link Count:", inode.linkCount)
  inode.linkCount--
  inode.destroyIfNeeded()
}

func (inode *Inode) incrementLinkCount() {
  // fmt.Println("||||| ++ Link Count:", inode.linkCount)
  inode.linkCount++
}

func (inode *Inode) decrementFileCount() {
  // fmt.Println("||||| -- File Count:", inode.fileCount)
  inode.fileCount--
  inode.destroyIfNeeded()
}

func (inode *Inode) incrementFileCount() {
  // fmt.Println("||||| ++ File Count:", inode.fileCount)
  inode.fileCount++
}

func initInode() *Inode {
  store := dstore.InitPageStore()
  // store := dstore.InitArrayStore(0)

  return &Inode{
    data: store,
    linkCount: 1,
    fileCount: 0,
  }
}
