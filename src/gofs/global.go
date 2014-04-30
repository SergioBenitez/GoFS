package gofs

import (
  "os"
  "errors"
  "gofs/dstore"
)

const USE_FILE_ARENA = true

const FILE_ARENA_SIZE = 100
const PAGE_ARENA_SIZE = 500

var pageArena *dstore.PageArena;
var fileArena *FileArena

type FileArena struct {
  files [FILE_ARENA_SIZE]*DataFile
  used int
  size int
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
  pageArena = nil
}

func ArenaAllocateDataFile(inode *Inode) (*DataFile, error) {
  if fileArena.used >= fileArena.size { 
    panic("Out of arena memory!")
    return nil, errors.New("Out Of Memory!")
  }

  file := fileArena.files[fileArena.used]
  file.seek = 0
  file.inode = inode;
  file.status = Open;

  fileArena.used += 1
  return file, nil
}

func ArenaReturnDataFile(file *DataFile) error {
  if fileArena.used <= 0 { return errors.New("Over-Freeing") }

  file.inode = nil;
  file.status = Closed;

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

  // Creating Page Arena
  if pageArena == nil {
    pageArena = dstore.InitPageArena(PAGE_ARENA_SIZE)
  }
}
