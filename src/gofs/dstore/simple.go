package dstore

import "errors"

type ArrayStore struct {
  data []byte
}

func (s *ArrayStore) Read(o int64, p []byte) (int, error) {
  if o > s.Size() { return 0, errors.New("EOF") }

  return copy(p, s.data[o:]), nil
}

func (s *ArrayStore) Write(o int64, p []byte) (int, error) {
  needed := o + int64(len(p))

  if needed - int64(cap(s.data)) > 0 {
    newData := make([]byte, needed, needed * 2)
    copy(newData, s.data)
    s.data = newData
  }

  s.data = s.data[:needed]
  return copy(s.data[o:], p), nil
}

func (s *ArrayStore) Size() int64 {
  return int64(len(s.data))
}

func InitArrayStore(alloc uint64) *ArrayStore {
  return &ArrayStore{
    data: make([]byte, 0, alloc),
  }
}
