package gofs

import (
  "testing"
  "bytes"
  "fmt"
  "runtime"
  "path/filepath"
  "math/rand"
  "time"
)

func randBytes(n int) []byte {
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

// print at most n callers
// so, why do we have this? because t.Log doesn't print the stack
func printStack(t *testing.T, n int) {
  if n == 0 { n = 10 }
  stack := make([]string, 0, n)

  for i := 0; i < n; i += 1 {
    _, file, line, ok := runtime.Caller(i + 3)
    if !ok { break }
    stack = append(stack, fmt.Sprintf("%s:%d", filepath.Base(file), line))
  }

  for i := range stack {
    t.Log(stack[len(stack) - i - 1])
  }
}

func AssertNoErr(t *testing.T, err error) {
  if err == nil { return }
  printStack(t, 0)
  t.Fatal(err)
}

func AssertTrue(t *testing.T, val bool, msg string) {
  if val { return }
  printStack(t, 0)
  t.Fatal(msg)
}

func AssertEqualBytes(t *testing.T, b1 []byte, b2 []byte) {
  equal := bytes.Equal(b1, b2)
  str := fmt.Sprintf("b1[%d] != b2[%d]\nb1: %v\nb2: %v",
    len(b1), len(b2), b1, b2)
  AssertTrue(t, equal, str)
}

func (p *ProcState) safeOpen(t *testing.T, s string,
f AccessFlag, m [3]FileMode) FileDescriptor {
  fd, err := p.Open(s, f, m)
  AssertNoErr(t, err)
  return fd
}

func (p *ProcState) safeSeek(t *testing.T, fd FileDescriptor,
off int64, whence int) int64 {
  n, err := p.Seek(fd, off, whence)
  AssertNoErr(t, err)
  return n
}

func (p *ProcState) safeRead(t *testing.T, fd FileDescriptor, b []byte) int {
  n, err := p.Read(fd, b)
  AssertNoErr(t, err)
  return n
}

func (p *ProcState) safeWrite(t *testing.T, fd FileDescriptor, b []byte) int {
  n, err := p.Write(fd, b)
  AssertNoErr(t, err)
  return n
}

func (p *ProcState) safeClose(t *testing.T, fd FileDescriptor) {
  err := p.Close(fd)
  AssertNoErr(t, err)
}

func (p *ProcState) safeUnlink(t *testing.T, s string) {
  err := p.Unlink(s)
  AssertNoErr(t, err)
}

func (p *ProcState) safeChdir(t *testing.T, s string) {
  err := p.Chdir(s)
  AssertNoErr(t, err)
}

func (p *ProcState) safeMkdir(t *testing.T, s string) {
  err := p.Mkdir(s)
  AssertNoErr(t, err)
}

func (p *ProcState) safeLink(t *testing.T, s string, s2 string) {
  err := p.Link(s, s2)
  AssertNoErr(t, err)
}

func (p *ProcState) safeRename(t *testing.T, s string, s2 string) {
  err := p.Rename(s, s2)
  AssertNoErr(t, err)
}

func TestEmptyRead(t *testing.T) {
  p := InitProc()
  filename := "file"
  buffer := make([]byte, 24)
  blank := make([]byte, 24)

  fd := p.safeOpen(t, filename, O_RDONLY | O_CREAT, UserMode())

  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer, blank)

  p.safeSeek(t, fd, 0, SEEK_SET)
  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer, blank)

  p.safeClose(t, fd)
  p.safeUnlink(t, filename)
}

func TestWriteRead(t *testing.T) {
  p := InitProc()
  filename := "file"
  content := []byte("Hello, world!")
  buffer := make([]byte, 24)
  buffer2 := make([]byte, 24)
  blank := make([]byte, 24)

  fd := p.safeOpen(t, filename, O_RDWR | O_CREAT, UserMode())
  p.safeWrite(t, fd, content)
  p.safeSeek(t, fd, 0, SEEK_SET)
  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer[:len(content)], content)

  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer2, blank)

  p.safeClose(t, fd)
  p.safeUnlink(t, filename)
}

func TestSeek(t *testing.T) {
  p := InitProc()
  filename := "file"
  size := 9240
  content := randBytes(size)
  buffer := make([]byte, size)

  fd := p.safeOpen(t, filename, O_RDWR | O_CREAT, UserMode())
  p.safeWrite(t, fd, content)
  p.safeSeek(t, fd, 0, SEEK_SET)
  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer[:len(content)], content)

  // randomly seek 5000 times and verify 250 bytes
  bytes := 250
  buf := make([]byte, bytes)
  for i := 0; i < 5000; i += 1 {
    pos := rand.Int63n(int64(size - bytes))
    p.safeSeek(t, fd, pos, SEEK_SET)
    p.safeRead(t, fd, buf)
    AssertEqualBytes(t, content[pos : pos + int64(bytes)], buf)
  }

  p.safeClose(t, fd)
  p.safeUnlink(t, filename)
}

func TestMkDirAndLink(t *testing.T) {
  p := InitProc()
  filename := "file"
  size := 24

  buffer := make([]byte, size)
  content1 := randBytes(size)
  content2 := randBytes(size)

  // Writing file to root
  fd := p.safeOpen(t, filename, O_RDWR | O_CREAT, UserMode())
  p.safeWrite(t, fd, content1)
  p.safeSeek(t, fd, 0, SEEK_SET)
  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer[:len(content1)], content1)
  p.safeClose(t, fd)

  // Now switching directory and writing 'file' again with content2
  p.safeMkdir(t, "mydir")
  p.safeChdir(t, "mydir/")

  fd = p.safeOpen(t, filename, O_RDWR | O_CREAT, UserMode())
  p.safeWrite(t, fd, content2)
  p.safeSeek(t, fd, 0, SEEK_SET)
  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer[:len(content2)], content2)
  p.safeClose(t, fd)

  // Verifying first file wasn't changed
  p.safeChdir(t, "/")
  fd = p.safeOpen(t, filename, O_RDONLY, UserMode())
  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer[:len(content1)], content1)
  p.safeClose(t, fd)

  // linking /mydir/file2 to /file and checking contents
  buffer = make([]byte, size)
  p.safeLink(t, filename, "/mydir/file2")
  fd = p.safeOpen(t, "/mydir/file2", O_RDONLY, UserMode())
  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer[:len(content1)], content1)
  p.safeClose(t, fd)
}

func TestRename(t *testing.T) {
  p := InitProc()
  filename := "file"
  filename2 := "another"
  size := 24
  buffer := make([]byte, size)
  content := randBytes(size)

  // write to the first file
  fd := p.safeOpen(t, filename, O_RDWR | O_CREAT, UserMode())
  p.safeWrite(t, fd, content)
  p.safeClose(t, fd)

  // rename file 1 to file 2
  p.safeRename(t, filename, filename2)

  // open and read from the second file
  fd = p.safeOpen(t, filename2, O_RDONLY, UserMode())
  p.safeRead(t, fd, buffer)
  AssertEqualBytes(t, buffer[:len(content)], content)
  p.safeClose(t, fd)
}
