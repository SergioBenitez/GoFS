package gofs

/**
* TODO: Figure out this directory stuff. Ideally, a directory would be a 'File'
* of type 'Directory' that you can read from.
*
* TODO: Think about what's hapenning. It seems to me that this whole file
* descriptor thing doesn't make a lot of sense in the context of a HLL with
* first class support for object.  Instead, it should return a 'File' type. A
* simple compatbility layer could be built on top of this where we keep a
* filetable mapping fds to this file type.
*
* In this way, Open("...") would grab a File and put a pointer to it in
* a table and return the index (or whatever) of that table. The rest of the
* file system functions would simply look up the File and call methods on it.
*
* class File ...
* instance RegularFile : File ... // not the right syntax
* readFile :: (File a) => a -> ByteString
* lookup :: (File a) => FileDescriptor -> a
* read :: FileDescriptor -> ByteString
* read = readFile . lookup
*/

// type File interface {
//   Open(path string, flags AccessFlag, mode [3]FileMode) FileDescriptor
//   Close(FileDescriptor) bool
//   Read(FileDescriptor, numBytes uint) ([]byte, uint)
//   Write(FileDescriptor, bytes []byte) uint
//   Seek(FileDescriptor, offset uint, whence Whence) uint
// }

func (p *ProcState) getNewFD() FileDescriptor {
  // Totally a race condition. Need access to docs to check out atomic
  // operation support in Go.
  next := p.lastFd + 1
  p.lastFd += 1
  return next
}

func (p *ProcState) allocateDescriptorInfo(flags AccessFlag) *DescriptorInfo {
  info := new(DescriptorInfo)

  return info
}

func (p *ProcState) Open(path string,
  flags AccessFlag, mode [3]FileMode) FileDescriptor {

  fd := p.getNewFD()
  p.fdTable[fd] = p.allocateDescriptorInfo(flags)

  return fd
}

func (p *ProcState) Close(fd FileDescriptor) bool {
  delete(p.fdTable, fd)
  return true
}

func UserMode() [3]FileMode {
  return [3]FileMode{M_READ | M_WRITE | M_EXEC, M_READ, M_READ}
}

func InitProc() *ProcState {
  state := new(ProcState)
  state.fdTable = make(FileDescriptorTable)
  return state
}
