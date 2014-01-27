package gofs

import (
  "strings"
)

func invalidPath(path string) bool {
  return strings.HasPrefix(path, ".")
}
