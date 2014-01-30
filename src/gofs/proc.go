package gofs

import (
  "sync/atomic"
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
func (proc *ProcState) initFileTableAndLastFD() {
  table := make(FileTable)
  table[0] = globalState.stdIn
  table[1] = globalState.stdOut
  table[2] = globalState.stdErr
  proc.fileTable = table;
  proc.lastFd = 2;
}

// Allocates a new file descriptor atomically.
func (proc *ProcState) getUnusedFd() FileDescriptor {
  var thing *int64 = (*int64)(&proc.lastFd)
  newthing := atomic.AddInt64(thing, 1)
  return FileDescriptor(newthing)
}

// Fetches the file object given a file descriptor
func (proc *ProcState) getFile(fd FileDescriptor) (interface{File}, error) {
  file, present := proc.fileTable[fd]
  if present { return file, nil } 
  return nil, errors.New("fd not found")
}

// Opens a file without returning a file descriptor.
// What happens if the filename in path is empty? IE: path = a/b/c/
func (proc *ProcState) openFile(path string, flags AccessFlag,
mode [3]FileMode) (interface{File}, error) {
  var err error
  dir, file, _ := proc.resolveFilePath(path)

  if file != nil {
    switch file.(type) {
    case *DataFile:
      file.(*DataFile).Open()
    default:
      err = errors.New("Cannot open file of this type.")
    }
  } else {
    switch {
      case (flags & O_CREAT) != 0:
        file = initDataFile()
        dir[path] = file
        file.(*DataFile).Open()
      default:
        err = errors.New("File not found.")
    }
  }

  return file.(interface{File}), err
}

func (proc *ProcState) Mkdir(path string) error {
  parentDir, dirName, err := proc.resolveDirPath(path)
  if (err != nil) { return err }

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

  dstDir[baseName] = file
  return nil
}



// Opens a file and returns a file descriptor.
func (proc *ProcState) Open(path string, flags AccessFlag,
mode [3]FileMode) (FileDescriptor, error) {
  file, err := proc.openFile(path, flags, mode)
  if err != nil { return FileDescriptor(-1), err }

  fd := proc.getUnusedFd()
  proc.fileTable[fd] = file
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
 * Resource freeing happens below. The memory of a file is freed after it is
 * both closed by all processes and unlinked from all directories. This is
 * because a file can only be referenced from two different locations: 1) a
 * file table, and 2) a directory.
 */

func (proc *ProcState) Unlink(path string) error {
  delete(proc.cwd, path)
  return nil
}

// Below are FD interfaces to the file calls.
func (proc *ProcState) Close(fd FileDescriptor) error {
  file, err := proc.getFile(fd)
  if err != nil {
    return errors.New("fd not found")
  }
  delete(proc.fileTable, fd)
  return file.Close()
}

func init() {
  initGlobalState()
}

func InitProc() *ProcState {
  state := new(ProcState)
  state.cwd = globalState.root
  state.initFileTableAndLastFD()
  return state
}
