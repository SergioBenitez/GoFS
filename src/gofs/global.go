package gofs

import (
  "errors"
  "gofs/dstore"
  "os"
)

const USE_FILE_ARENA = true
const FILE_ARENA_SIZE = 100
const PAGE_ARENA_SIZE = 256 * 4 // 4MB

var fileArena *FileArena

type FileArena struct {
  files [FILE_ARENA_SIZE]*DataFile
  used  int
  size  int
}

func initDirectory(parent Directory) Directory {
  dir := make(Directory)
  dir["."] = dir
  if parent == nil {
    dir[".."] = dir
  } else {
    dir[".."] = parent
  }
  return dir
}

func (dir Directory) parent() Directory {
  return dir[".."].(Directory)
}

func ClearGlobalState() {
  globalState = nil
  fileArena = nil
  dstore.GlobalPageArena = nil
}

func ArenaAllocateDataFile(inode *Inode) (*DataFile, error) {
  if fileArena.used >= fileArena.size {
    panic("Out of arena memory!")
    return nil, errors.New("Out Of Memory!")
  }

  file := fileArena.files[fileArena.used]
  file.seek = 0
  file.inode = inode
  file.status = Open

  fileArena.used += 1
  return file, nil
}

func ArenaReturnDataFile(file *DataFile) error {
  if fileArena.used <= 0 {
    return errors.New("Over-Freeing")
  }

  file.inode = nil
  file.status = Closed

  fileArena.used -= 1
  fileArena.files[fileArena.used] = file
  return nil
}

func InitGlobalState() {
  if globalState == nil {
    globalState = new(GlobalState)
    globalState.root = initDirectory(nil)
    globalState.stdIn = os.Stdin
    globalState.stdOut = os.Stdout
    globalState.stdErr = os.Stderr
  }

  // Creating FileArena
  if fileArena == nil {
    fileArena = &FileArena{
      used: 0,
      size: FILE_ARENA_SIZE,
    }

    for i := 0; i < FILE_ARENA_SIZE; i++ {
      fileArena.files[i] = &DataFile{}
    }
  }

  if dstore.GlobalPageArena == nil {
    dstore.GlobalPageArena = dstore.InitPageArena(PAGE_ARENA_SIZE)
  }
}
