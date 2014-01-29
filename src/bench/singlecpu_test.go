package bench

import (
  "gofs"
  "testing"
)

const NUM = 100

func ceilDiv(x int, y int) int {
  return (x + y - 1) / y
}

func openMany(b *testing.B, p *gofs.ProcState, n int) []gofs.FileDescriptor {
  return openManyC(b, p, n, func(gofs.FileDescriptor, string) { })
}

func openManyC(b *testing.B, p *gofs.ProcState, n int, 
f func(gofs.FileDescriptor, string)) []gofs.FileDescriptor {
  b.StopTimer()
  fds := make([]gofs.FileDescriptor, n)
  mode := gofs.UserMode()
  filename := make([]byte, ceilDiv(n, 26))
  for i := range filename { filename[i] = 'a' }
  b.StartTimer()

  for i := range fds {
    var err error
    fds[i], err = p.Open(string(filename), gofs.O_CREAT, mode)
    if err != nil { b.Fatal("bad open") }
    f(fds[i], string(filename))
    filename[i / 26] += 1
  }

  return fds
}

func closeAll(b *testing.B, p *gofs.ProcState, fs []gofs.FileDescriptor) {
  var err error
  for _, fd := range fs {
    err = p.Close(fd)
    if err != nil { b.Fatal("bad close") }
  }
}

func unlinkAll(b *testing.B, p *gofs.ProcState, fs []gofs.FileDescriptor) {
  b.StopTimer()
  filename := make([]byte, ceilDiv(len(fs), 26))
  for i := range filename { filename[i] = 'a' }
  b.StartTimer()

  for i := range fs {
    err := p.Unlink(string(filename))
    if err != nil { b.Fatal("bad unlink") }
    filename[i / 26] += 1
  }
}

func newProc(b *testing.B) *gofs.ProcState {
  b.StopTimer()
  p := gofs.InitProc()
  b.StartTimer()
  return p
}

func BenchmarkOtC(b *testing.B) {
  for j := 0; j < b.N; j++ {
    p := newProc(b)
    fds := openMany(b, p, NUM)
    closeAll(b, p, fds)
  }
}

func BenchmarkOC(b *testing.B) {
  for j := 0; j < b.N; j++ {
    p := newProc(b)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, _ string) {
      p.Close(fd)
    })
  }
}

func BenchmarkOCInLn(b *testing.B) {
  for j := 0; j < b.N; j++ {
    b.StopTimer()
    p := gofs.InitProc()
    fds := make([]gofs.FileDescriptor, NUM)
    mode := gofs.UserMode()
    filename := make([]byte, ceilDiv(NUM, 26))
    for i := range filename { filename[i] = 'a' }
    b.StartTimer()

    for i := range fds {
      var err error
      fds[i], err = p.Open(string(filename), gofs.O_CREAT, mode)
      if err != nil { b.Fatal("bad open") }
      p.Close(fds[i])
      filename[i / 26] += 1
    }
  }
}

func BenchmarkOtCtU(b *testing.B) {
  for j := 0; j < b.N; j++ {
    p := newProc(b)
    fds := openMany(b, p, NUM)
    closeAll(b, p, fds)
    unlinkAll(b, p, fds)
  }
}

func BenchmarkOCU(b *testing.B) {
  for j := 0; j < b.N; j++ {
    p := newProc(b)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, s string) {
      p.Close(fd)
      p.Unlink(s)
    })
  }
}
