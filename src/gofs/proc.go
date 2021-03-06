package gofs

import (
  "errors"
)

var globalState *GlobalState

/**
* This file contains the code to manage the state of a process in GoFS.
* Specifically, it provide the Open call and manages the file descriptor mapping
* from fd to File.
*/

func UserMode() [3]FileMode {
  return [3]FileMode{M_READ | M_WRITE | M_EXEC, M_READ, M_READ}
}

// Sets up the initial file table to point to std out, in, and err.
func (proc *ProcState) initFileDescriptorTableAndLastFD() {
  table := make(FileDescriptorTable)
  table[0] = globalState.stdIn
  table[1] = globalState.stdOut
  table[2] = globalState.stdErr
  proc.fileDescriptorTable = table;

  // lastFd keeps track of the index of the last used FD
  proc.lastFd = 0
  for i := 0; i < MAX_DESCRIPTORS; i++ {
    proc.freeDescriptors[i] = FileDescriptor(i + 3); // since 0, 1, 2 are taken
  }
}

// Allocates a new file descriptor (not) atomically.
func (proc *ProcState) getUnusedFd() (fd FileDescriptor) {
  // Below is what we used to do for the atomic stuff
  // var thing *int64 = (*int64)(&proc.lastFd)
  // newthing := atomic.AddInt64(thing, 1)
  if proc.lastFd >= MAX_DESCRIPTORS { panic("Out of FDs!") }
  fd = proc.freeDescriptors[proc.lastFd]
  proc.lastFd += 1
  return
}

func (proc *ProcState) returnFd(fd FileDescriptor) {
  if proc.lastFd <= 0 { panic("Overfreeing FDs!") }
  proc.lastFd -= 1
  proc.freeDescriptors[proc.lastFd] = fd;
}

// Fetches the file object given a file descriptor
func (proc *ProcState) getFile(fd FileDescriptor) (interface{File}, error) {
  file, present := proc.fileDescriptorTable[fd]
  if present { return file, nil } 
  return nil, errors.New("fd not found")
}

// Opens a file without returning a file descriptor.
// What happens if the filename in path is empty? IE: path = a/b/c/
func (proc *ProcState) openFile(path string, flags AccessFlag,
mode [3]FileMode) (interface{File}, error) {
  var err error; var inode *Inode
  dir, filename, _ := proc.resolveDirPath(path)
  file, ok := dir[filename]

  // Finding our *Inode, if possible.
  if ok {
    switch file.(type) {
    case *Inode:
      inode = file.(*Inode)
    default:
      return nil, errors.New("Cannot open file of this type.")
    }
  } else {
    switch {
      case (flags & O_CREAT) != 0:
        inode = initInode()
        dir[filename] = inode
      default:
        return nil, errors.New("File not found.")
    }
  }

  // We're here? We found it! Otherwise, would have err.
  file = initDataFile(inode)
  return file.(interface{File}), err
}

func (proc *ProcState) Mkdir(path string) error {
  parentDir, dirName, err := proc.resolveDirPath(path)
  if (err != nil) { return err }

  _, exists := parentDir[dirName]
  if exists { return errors.New("Destination already exists.") }

  parentDir[dirName] = initDirectory(parentDir)
  return nil
}

func (proc *ProcState) Chdir(path string) error {
  dir, _, err := proc.resolveDirPath(path)
  if err != nil { return err }

  proc.cwd = dir
  return nil
}

func (proc *ProcState) Link(src string, dst string) error {
  var err error

  _, file, err := proc.resolveFilePath(src)
  if err != nil { return err }

  dstDir, baseName, err := proc.resolveDirPath(dst)
  if err != nil { return err }

  _, exists := dstDir[baseName]
  if exists { return errors.New("Destination file already exists.") }

  switch inode := file.(type) {
  case *Inode:
    inode.incrementLinkCount()
  }

  dstDir[baseName] = file
  return nil
}

func (proc *ProcState) Rename(src string, dst string) error {
  err := proc.Link(src, dst)
  if err != nil { return err }

  err = proc.Unlink(src)
  return err
}

// Opens a file and returns a file descriptor.
func (proc *ProcState) Open(path string, flags AccessFlag,
mode [3]FileMode) (FileDescriptor, error) {
  file, err := proc.openFile(path, flags, mode)
  if err != nil { return FileDescriptor(-1), err }

  fd := proc.getUnusedFd()
  proc.fileDescriptorTable[fd] = file
  return fd, nil
}

func (proc *ProcState) Read(fd FileDescriptor, p []byte) (n int, err error) {
  file, err := proc.getFile(fd)
  if err != nil { return 0, err }
  return file.Read(p)
}

func (proc *ProcState) Write(fd FileDescriptor, p []byte) (n int, err error) {
  file, err := proc.getFile(fd)
  if err != nil { return 0, err }
  return file.Write(p)
}

func (proc *ProcState) Seek(fd FileDescriptor, offset int64, whence int) (int64, error) {
  file, err := proc.getFile(fd)
  if err != nil { return 0, err }
  return file.Seek(offset, whence)
}

/**
 * Resource freeing happens below. The memory of an inode is freed after it is
 * referenced by no open files and unlinked from all directories. This is
 * because an inode can only be referenced from two different locations: 
 * 1) a file, and 2) a directory.
 */

func (proc *ProcState) Unlink(path string) error {
  dir, name, err := proc.resolveDirPath(path)
  if err != nil { return err }

  file, ok := dir[name]
  if !ok { return errors.New("Cannot unlink nonexisting file.") }

  switch inode := file.(type) {
  case *Inode:
    inode.decrementLinkCount()
  }

  delete(dir, name)
  return nil
}

func (proc *ProcState) Close(fd FileDescriptor) error {
  file, err := proc.getFile(fd)
  if err != nil { return errors.New("fd not found") }

  proc.returnFd(fd)
  delete(proc.fileDescriptorTable, fd)
  return file.Close()
}

func init() {
  InitGlobalState()
}

func InitProc() *ProcState {
  state := new(ProcState)
  state.cwd = globalState.root
  state.initFileDescriptorTableAndLastFD()
  return state
}
