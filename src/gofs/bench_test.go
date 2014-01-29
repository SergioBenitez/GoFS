package gofs

import (
  "testing"
)

// open, then close, 100 files
func BenchmarkOpen(b *testing.B) {
  p := InitProc()
  n := 100

  for j := 0; j < b.N; j += 1 {
    name := []byte("aaaaaaa")
    mode := UserMode()
    for i := 0; i < n; i += 1 {
      index := i / 26
      name[index] += 1

      fd, err := p.Open(string(name), O_CREAT, mode)
      if err != nil { b.Fatal("whoops") }
      p.Close(fd)
    }
  }
}
