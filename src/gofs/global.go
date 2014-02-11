package gofs

import "os"

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

func InitGlobalState() {
  if globalState == nil {
    globalState = new(GlobalState)
    globalState.root = initDirectory(nil)
    globalState.stdIn = os.Stdin
    globalState.stdOut = os.Stdout
    globalState.stdErr = os.Stderr
  }
}
