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

func (proc *ProcState) GetFile(fd FileDescriptor) interface{File} {
  return proc.fileTable[fd]
}

func (proc *ProcState) Open(name string, flags AccessFlag,
mode [3]FileMode) (FileDescriptor, error) {
  file, err := proc.OpenX(name, flags, mode)
  if err != nil { return FileDescriptor(-1), err }

  fd := proc.getUnusedFd()
  proc.fileTable[fd] = file
  return fd, nil
}

func (proc *ProcState) OpenX(name string, flags AccessFlag,
mode [3]FileMode) (interface{File}, error) {
  var err error = nil
  file, present := proc.cwd[name]

  if invalidPath(name) {
    return nil, errors.New("Invalid path.")
  }

  if present {
    switch file.(type) {
    case *DataFile:
    default:
      err = errors.New("Cannot open file of this type.")
    }
  } else {
    switch {
      case (flags & O_CREAT) != 0:
        file = initDataFile()
        proc.cwd[name] = file
      default:
        err = errors.New("File not found.")
    }
  }

  return file.(interface{File}), err
}

func InitProc() *ProcState {
  initGlobalState()
  state := new(ProcState)
  state.cwd = globalState.root
  state.initFileTableAndLastFD()
  return state
}
