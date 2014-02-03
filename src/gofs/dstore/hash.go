package dstore

import "errors"

type HashStore struct {
  blockSize int
  data [][]byte
}

/**
* Should refactor this to have a COPY function from some interface that looks
* like [][]byte to [][]byte with some method "atIndex". IE:
*   copy(a, b IndexableBytes)
* 
* Then, can create simple wrapper around []byte where atIndex(n) returns the
* same []byte for all n.
*/

func (s *HashStore) Read(o int, p []byte) (n int, e error) {
  if o >= s.Size() { return 0, errors.New("EOF") }

  // copy from first block
  i, off := o / s.blockSize, o % s.blockSize
  if i < len(s.data) {
    n += copy(p, s.data[i][off:s.blockSize]) 
  }

  // copy from rest of blocks
  for i = i + 1; i < len(s.data) && n < len(p); i++ {
    n += copy(p[n:], s.data[i][:s.blockSize])
  }

  return
}

func (s *HashStore) Write(o int, p []byte) (n int, err error) {
  // determining how many blocks are needed
  needed := (o + len(p) + s.blockSize - 1) / s.blockSize

  // expanding blocks data array if needed
  if needed - cap(s.data) > 0 {
    newData := make([][]byte, needed, needed * 2)
    copy(newData, s.data)
    s.data = newData
  }

  // allocating empty data blocks if needed, writing to first block
  s.expandTo(((o + len(p)) / s.blockSize) + 1)
  i, off := o / s.blockSize, o % s.blockSize
  n += copy(s.data[i][off:s.blockSize], p)

  // updating headers
  s.data[i] = s.data[i][:n]

  // writing to rest of block and updating headers in process
  for i = i + 1; i < len(s.data) && n < len(p); i++ {
    written := copy(s.data[i][:s.blockSize], p[n:])
    s.data[i] = s.data[i][:written]
    n += written
  }

  return 0, nil
}

func (s *HashStore) expandTo(length int) {
  if cap(s.data) < length {
    panic("HashStore:expandTo: callers must ensure array has enough capacity.")
  }

  s.data = s.data[:length]
  for i := 0; i < length; i++ {
    s.data[i] = make([]byte, 0, s.blockSize)
  }
}

func (s *HashStore) Size() int {
  if len(s.data) == 0 { return 0 }

  lastI := len(s.data) - 1
  return s.blockSize * lastI + len(s.data[lastI])
}

func InitHashStore(blockSize int) *HashStore {
  return &HashStore{
    blockSize: blockSize,
    data: make([][]byte, 0, 4096),
  }
}
