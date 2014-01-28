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

/*
 * Sets up the initial file table to point to std out, in, and err.
 */
func (proc *ProcState) initFileTableAndLastFD() {
  table := make(FileTable)
  table[0] = globalState.stdIn
  table[1] = globalState.stdOut
  table[2] = globalState.stdErr
  proc.fileTable = table;
  proc.lastFd = 2;
}

func (proc *ProcState) getUnusedFd() FileDescriptor {
  var thing *int64 = (*int64)(&proc.lastFd)
  newthing := atomic.AddInt64(thing, 1)
  return FileDescriptor(newthing)
}

func (proc *ProcState) GetFile(fd FileDescriptor) (interface{File}, error) {
  file, present := proc.fileTable[fd]
  if present { return file, nil } 
  return nil, errors.New("fd not found")
}

func (proc *ProcState) Open(path string, flags AccessFlag,
mode [3]FileMode) (FileDescriptor, error) {
  file, err := proc.OpenX(path, flags, mode)
  if err != nil { return FileDescriptor(-1), err }

  fd := proc.getUnusedFd()
  proc.fileTable[fd] = file
  return fd, nil
}

func (proc *ProcState) OpenX(path string, flags AccessFlag,
mode [3]FileMode) (interface{File}, error) {
  var err error = nil
  file, present := proc.cwd[path]

  if invalidPath(path) {
    return nil, errors.New("Invalid path.")
  }

  if present {
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
        proc.cwd[path] = file
        file.(*DataFile).Open()
      default:
        err = errors.New("File not found.")
    }
  }

  return file.(interface{File}), err
}

/**
 * Resource freeing happens below. The memory of a file is freed after it is
 * both closed by all processes and unlinked from all directories. This is
 * because a file can only be referenced from two different locations: 1) a
 * file table, and 2) a directory.
 */

func (proc *ProcState) Close(fd FileDescriptor) error {
  file, err := proc.GetFile(fd)
  if err != nil {
    return errors.New("fd not found")
  }
  delete(proc.fileTable, fd)
  return file.Close()
}

func (proc *ProcState) Unlink(path string) error {
  delete(proc.cwd, path)
  return nil
}

func InitProc() *ProcState {
  initGlobalState()
  state := new(ProcState)
  state.cwd = globalState.root
  state.initFileTableAndLastFD()
  return state
}
