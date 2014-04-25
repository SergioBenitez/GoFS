package gofs

import (
  "strings"
  "errors"
  "fmt"
)

func splitPath(path string) (dir string, base string) {
  index := -1
  for i := len(path) - 1; i >= 0; i-- {
    if path[i] == '/' {
      index = i
      break
    }
  }

  if index == -1 { return "", path }
  return path[:index], path[index + 1:]
}

func (proc *ProcState) resolveFilePath(path string) (Directory, interface{}, error) {
  dir, fileName, err := proc.resolveDirPath(path)
  if err != nil { return dir, nil, err }

  file, ok := dir[fileName]
  if !ok { return dir, nil, errors.New("File not found") }

  return dir, file, nil
}

// Returns the directory for a given path and the filename in that path.
// Example: a/b/c.txt returns the Directory for b in a and the string 'c.txt'
// Example: a/b/c/ return the c Directory and the string ""
// This isn't perfect yet: 
//  should handle multiple // in path ... it might do this
func (proc *ProcState) resolveDirPath(path string) (Directory, string, error) {
  // shouldn't do anything in this case
  if len(path) == 0 { return proc.cwd, "", nil }

  dirPath, fileName := splitPath(path)
  dirs := strings.Split(dirPath, "/")

  cwd := proc.cwd
  if path[0] == '/' {
    cwd = globalState.root 
    dirs = dirs[1:]
  }

  for _, name := range dirs {
    if name == "" { continue }

    dir := cwd[name]
    switch dir.(type) {
      case Directory:
        cwd = dir.(Directory)
      default:
        errString := fmt.Sprintf("Invalid path: %s in %s", name, path)
        return nil, "", errors.New(errString)
    }
  }

  return cwd, fileName, nil
}
