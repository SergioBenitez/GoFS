package main

import (
  "gofs"
  "fmt"
  "math/rand"
  "time"
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

// Gets calls automatically by the runtime.
func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
  proc := gofs.InitProc()

  fmt.Println("------first time through-------\n")
  fd, err := proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
  if err != nil { fmt.Println(err) }

  rand1 := randBytes(2500)
  rand2 := randBytes(4800)

  proc.Write(fd, rand1)
  proc.Write(fd, rand2)

  buffer := make([]byte, 11)
  proc.Seek(fd, 0, gofs.SEEK_SET)
  n, err := proc.Read(fd, buffer)
  if err !=  nil {
    panic(err)
  } else {
    fmt.Println("good read", n)
  }

  fmt.Println("Need:", rand1[:11])
  fmt.Println("Got:", buffer, "\n")

  buffer = make([]byte, 11)
  proc.Seek(fd, 2600, gofs.SEEK_SET)
  proc.Read(fd, buffer)

  fmt.Println("Need:", rand2[100:111])
  fmt.Println("Got:", buffer, "\n")

  proc.Close(fd)

  fmt.Println("\n\n------second time through-------\n")
  fd, err = proc.Open("file", gofs.O_RDWR, gofs.UserMode())
  if err != nil { fmt.Println(err) }

  buffer = make([]byte, 25)
  fmt.Println("Reading...")
  proc.Read(fd, buffer)
  fmt.Println(buffer)
  proc.Close(fd)
  proc.Unlink("file")

  fmt.Println("\n\n------second time through-------\n")
  fd, err = proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
  if err != nil { fmt.Println(err) }

  buffer = make([]byte, 25)
  proc.Read(fd, buffer)
  fmt.Println(buffer)
}
