package main

import (
  "gofs"
  "fmt"
  "math/rand"
)

func randBytes(n int) []byte {
  if (n & 0x3) != 0 { panic("n must be a multiple of 4") }

  out := make([]byte, n)
  for i := 0; i < n / 4; i += 1 {
    rand32 := rand.Uint32()
    for j := 0; j < 4; j += 1 {
      out[i * 4 + j] = byte((rand32 >> (uint(j) * 8)) & 0xFF)
    }
  }

  return out
}

func main() {
  proc := gofs.InitProc()

  fmt.Println("------first time through-------\n")
  fd, err := proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
  if err != nil { fmt.Println(err) }

  rand1 := randBytes(2500)
  rand2 := randBytes(4800)
  fmt.Println(rand1)

  file := proc.GetFile(fd)
  file.Write(rand1)
  file.Write(rand2)

  buffer := make([]byte, 11)
  file.Seek(0, gofs.SEEK_SET)
  file.Read(buffer)
  fmt.Println(string(buffer))

  buffer = make([]byte, 11)
  file.Read(buffer)
  fmt.Println(string(buffer))
  file.Close()

  fmt.Println("\n\n------second time through-------\n")
  fd, err = proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
  if err != nil { fmt.Println(err) }

  file = proc.GetFile(fd)
  buffer = make([]byte, 25)
  file.Read(buffer)
  fmt.Println(string(buffer))
}
