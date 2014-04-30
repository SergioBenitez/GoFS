package main

import (
  "gofs"
  "math/rand"
  "time"
  "os"
  "strconv"
  // "fmt"
)

const NUM int = 100

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

func ceilDiv(x int, y int) int {
  return (x + y - 1) / y
}

func openManyC(p *gofs.ProcState, n int, f func(gofs.FileDescriptor, string)) []gofs.FileDescriptor {
  fds := make([]gofs.FileDescriptor, n)
  mode := gofs.UserMode()
  filename := make([]byte, ceilDiv(n, 26))
  for i := range filename { filename[i] = '@' }

  for i := range fds {
    var err error
    filename[i / 26] += 1
    fds[i], err = p.Open(string(filename), gofs.O_CREAT, mode)
    if err != nil { panic("bad open") }
    f(fds[i], string(filename))
  }

  return fds
}

func newProc() *gofs.ProcState {
  // fmt.Println("new proc")
  // Should we be clearing the global state?
  gofs.ClearGlobalState()
  // runtime.GC()
  gofs.InitGlobalState()

  return gofs.InitProc()
}

func BenchmarkOWbC(reps int) {
  size := 40960

  p := newProc()
  for j := 0; j < reps; j++ {
    content := randBytes(size)
    openManyC(p, NUM, func(fd gofs.FileDescriptor, _ string) {
      p.Write(fd, content)
      p.Close(fd)
    })
    // runtime.GC()
  }
}

func BenchmarkOWbCU(reps int) {
  size := 40960

  p := newProc()
  for j := 0; j < reps; j++ {
    content := randBytes(size)
    openManyC(p, NUM, func(fd gofs.FileDescriptor, s string) {
      p.Write(fd, content)
      p.Close(fd)
      p.Unlink(s)
    })
    // runtime.GC()
  }
}

func main() {
  // fmt.Println("Hello!");
  if len(os.Args) < 3 { panic("Need 2 args.") }
  rep1, _ := strconv.Atoi(os.Args[1])
  rep2, _ := strconv.Atoi(os.Args[2])
  BenchmarkOWbC(rep1);
  BenchmarkOWbCU(rep2);
  // fmt.Println("Test");
}

// func main() {
//   proc := gofs.InitProc()

//   fmt.Println("------first time through-------\n")
//   fd, err := proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
//   if err != nil { fmt.Println(err) }

//   rand1 := randBytes(2500)
//   rand2 := randBytes(4800)

//   proc.Write(fd, rand1)
//   proc.Write(fd, rand2)

//   buffer := make([]byte, 11)
//   proc.Seek(fd, 0, gofs.SEEK_SET)
//   n, err := proc.Read(fd, buffer)
//   if err !=  nil {
//     panic(err)
//   } else {
//     fmt.Println("good read", n)
//   }

//   fmt.Println("Need:", rand1[:11])
//   fmt.Println("Got:", buffer, "\n")

//   buffer = make([]byte, 11)
//   proc.Seek(fd, 2600, gofs.SEEK_SET)
//   proc.Read(fd, buffer)

//   fmt.Println("Need:", rand2[100:111])
//   fmt.Println("Got:", buffer, "\n")

//   proc.Close(fd)

//   fmt.Println("\n\n------second time through-------\n")
//   fd, err = proc.Open("file", gofs.O_RDWR, gofs.UserMode())
//   if err != nil { fmt.Println(err) }

//   buffer = make([]byte, 25)
//   fmt.Println("Reading...")
//   proc.Read(fd, buffer)
//   fmt.Println(buffer)
//   proc.Close(fd)
//   proc.Unlink("file")

//   fmt.Println("\n\n------second time through-------\n")
//   fd, err = proc.Open("file", gofs.O_RDWR | gofs.O_CREAT, gofs.UserMode())
//   if err != nil { fmt.Println(err) }

//   buffer = make([]byte, 25)
//   proc.Read(fd, buffer)
//   fmt.Println(buffer)
// }

