package bench

import (
  // "fmt"
  "gofs"
  "time"
  "math/rand"
  "testing"
  // "runtime"
)

const NUM = 100

func ceilDiv(x int, y int) int {
  return (x + y - 1) / y
}

func randBytes(b *testing.B, n int) []byte {
  b.StopTimer()
  defer b.StartTimer()

  rand.Seed(time.Now().UTC().UnixNano())
  if (n & 0x3) != 0 { panic("randBytes: n must be a multiple of 4") }

  out := make([]byte, n)
  for i := 0; i < n / 4; i += 1 {
    rand32 := rand.Uint32()
    for j := 0; j < 4; j += 1 {
      out[i * 4 + j] = byte((rand32 >> (uint(j) * 8)) & 0xFF)
    }
  }

  return out
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
  for i := range filename { filename[i] = '@' }
  b.StartTimer()

  for i := range fds {
    var err error
    filename[i / 26] += 1
    fds[i], err = p.Open(string(filename), gofs.O_CREAT, mode)
    if err != nil { b.Fatal("bad open") }
    f(fds[i], string(filename))
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
  for i := range filename { filename[i] = '@' }
  b.StartTimer()

  for i := range fs {
    filename[i / 26] += 1
    err := p.Unlink(string(filename))
    if err != nil { b.Fatal("bad unlink:", err) }
  }
}

func newProc(b *testing.B) *gofs.ProcState {
  b.StopTimer()
  defer b.StartTimer()

  // fmt.Println("new proc")
  // Should we be clearing the global state?
  gofs.ClearGlobalState()
  // runtime.GC()
  gofs.InitGlobalState()

  return gofs.InitProc()
}

func BenchmarkOC1(b *testing.B) {
  p := newProc(b)
  for j := 0; j < b.N; j++ {
    // for i := 0; i < 17108864; i++ { // string garbage?
      fd, err := p.Open("test", gofs.O_CREAT, gofs.UserMode())
      if err != nil { b.Fatal("bad open") }
      p.Close(fd)
    // }
    // runtime.GC()
  }
}

func BenchmarkOtC(b *testing.B) {
  // fmt.Println("Creating new process...")
  p := newProc(b)
  for j := 0; j < b.N; j++ {
    // fmt.Println("Opening many...")
    fds := openMany(b, p, NUM)
    // fmt.Println("Closing many...")
    closeAll(b, p, fds)
    // runtime.GC()
  }
}

func BenchmarkOC(b *testing.B) {
  p := newProc(b)
  for j := 0; j < b.N; j++ {
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, _ string) {
      p.Close(fd)
    })
    // runtime.GC()
  }
}

func BenchmarkOtCtU(b *testing.B) {
  p := newProc(b)
  for j := 0; j < b.N; j++ {
    fds := openMany(b, p, NUM)
    closeAll(b, p, fds)
    unlinkAll(b, p, fds)
    // runtime.GC()
  }
}

func BenchmarkOCU(b *testing.B) {
  p := newProc(b)
  for j := 0; j < b.N; j++ {
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, s string) {
      p.Close(fd)
      p.Unlink(s)
    })
    // runtime.GC()
  }
}

func BenchmarkOWsC(b *testing.B) {
  size := 1024

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, size)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, _ string) {
      p.Write(fd, content)
      p.Close(fd)
    })
    // runtime.GC()
  }
}

func BenchmarkOWsCU(b *testing.B) {
  size := 1024

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, size)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, s string) {
      p.Write(fd, content)
      p.Close(fd)
      p.Unlink(s)
    })
    // runtime.GC()
  }
}

func BenchmarkOWbC(b *testing.B) {
  size := 40960

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, size)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, _ string) {
      p.Write(fd, content)
      p.Close(fd)
    })
    // runtime.GC()
  }
}

func BenchmarkOWbCU(b *testing.B) {
  size := 40960

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, size)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, s string) {
      p.Write(fd, content)
      p.Close(fd)
      p.Unlink(s)
    })
    // runtime.GC()
  }
}

func BenchmarkOWMsC(b *testing.B) {
  size := 1024
  many := 4096

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, size)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, _ string) {
      for i := 0; i < many; i++ {
        p.Write(fd, content)
      }
      p.Close(fd)
    })
    // runtime.GC()
  }
}

func BenchmarkOWMsCU(b *testing.B) {
  size := 1024
  many := 4096

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, size)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, s string) {
      for i := 0; i < many; i++ {
        p.Write(fd, content)
      }
      p.Close(fd)
      p.Unlink(s)
    })
    // runtime.GC()
  }
}

// the following tests need
// size * many * NUM bytes
// for NUM = 100, size = 1MB, many = 32, this is 3.125GB

func BenchmarkOWMbC(b *testing.B) {
  size := 1048576
  many := 32

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, size)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, _ string) {
      for i := 0; i < many; i++ {
        p.Write(fd, content)
      }
      p.Close(fd)
    })
    // runtime.GC()
  }
}

func BenchmarkOWMbCU(b *testing.B) {
  size := 1048576
  many := 32

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, size)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, s string) {
      for i := 0; i < many; i++ {
        p.Write(fd, content)
      }
      p.Close(fd)
      p.Unlink(s)
    })
  }
}

// the following two tests need
// NUM * startSize * many * (many - 1) / 2 bytes
// of memory.
//
// for NUM = 100, startSize = 2, many = 4096, this is ~1.56GB

func BenchmarkOWMbbC(b *testing.B) {
  startSize := 2
  many := 4096

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, startSize * many)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, _ string) {
      for i := 1; i <= many; i++ {
        p.Write(fd, content[:startSize * i])
      }
      p.Close(fd)
    })
    // runtime.GC()
  }
}

func BenchmarkOWMbbCU(b *testing.B) {
  startSize := 2
  many := 4096

  p := newProc(b)
  for j := 0; j < b.N; j++ {
    content := randBytes(b, startSize * many)
    openManyC(b, p, NUM, func(fd gofs.FileDescriptor, s string) {
      for i := 1; i <= many; i++ {
        p.Write(fd, content[:startSize * i])
      }
      p.Close(fd)
      p.Unlink(s)
    })
    // runtime.GC()
  }
}
