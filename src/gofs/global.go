package gofs

import (
  "os"
  "errors"
)

var fileArena *FileArena

type FileArena struct {
  files [64]*DataFile
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
}

func ArenaAllocateDataFile(inode *Inode) (*DataFile, error) {
  if fileArena.used >= fileArena.size { 
    panic("Out of memory!")
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
    var files [64]*DataFile
    for i := 0; i < 64; i++ {
      files[i] = &DataFile{}
    }

    fileArena = &FileArena{
      files: files,
      used: 0,
      size: 64,
    }
  }
}
