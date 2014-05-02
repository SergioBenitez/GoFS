package dstore

import (
  "errors"
)

const ENTRIES = 256

// Up to 257MB if ENTRIES = 256, PAGE_SIZE = 4096
// = ENTRIES * PAGE_SIZE + ENTRIES^2 * PAGE_SIZE
type PageStore struct {
  single *[ENTRIES][]byte            // 1MB
  double *[ENTRIES]*[ENTRIES][]byte   // 256MB
  pagesUsed int
  lastEntryBytesUsed int
}

func ceilDiv(x int, y int) int {
  return (x + y - 1) / y
}

func (s *PageStore) getEntry(num int) *[]byte {
  if num < ENTRIES {
    if s.single == nil { s.single = new([ENTRIES][]byte) }
    return &s.single[num]
  }

  doubleEntry := num - ENTRIES
  slot := doubleEntry / ENTRIES
  entryOffset := doubleEntry % ENTRIES
  if s.double == nil { s.double = new([ENTRIES]*[ENTRIES][]byte) }
  if s.double[slot] == nil { s.double[slot] = new([ENTRIES][]byte) }
  return &s.double[slot][entryOffset]
}

func (s *PageStore) Read(o int, p []byte) (int, error) {
  if o >= s.Size() { return 0, errors.New("EOF") }

  offset := o % PAGE_SIZE
  start := o / PAGE_SIZE
  entriesToRead := ceilDiv(len(p) + offset, PAGE_SIZE)

  read := 0
  for entry := 0; entry < entriesToRead; entry++ {
    page := *s.getEntry(start + entry)
    if (page == nil) { panic("Page is nil?") }
    read += copy(p[read:], page[offset:])
    if offset != 0 { offset = 0 }
  }

  return read, nil
}

func (s *PageStore) Write(o int, p []byte) (int, error) {
  offset := o % PAGE_SIZE
  start := o / PAGE_SIZE
  entriesToWrite := ceilDiv(len(p) + offset, PAGE_SIZE)

  written := 0
  for entry := 0; entry < entriesToWrite; entry++ {
    page := s.getEntry(start + entry)
    if *page == nil { *page = GlobalPageArena.AllocatePage() }
    if *page == nil { panic("Page was not allocated!") }

    written += copy((*page)[offset:], p[written:])
    if offset != 0 { offset = 0 }
  }

  if (start + entriesToWrite) >= s.pagesUsed {
    s.pagesUsed = start + entriesToWrite
    s.lastEntryBytesUsed = (len(p) + offset) % PAGE_SIZE
  }

  return written, nil
}

func (s *PageStore) Size() int {
  if s.pagesUsed == 0 { return 0 }
  return (s.pagesUsed - 1) * PAGE_SIZE + s.lastEntryBytesUsed
}

// Releases all pages in a singly-indirect block of pages
func (s *PageStore) ReleaseSinglePages(pages *[ENTRIES][]byte) {
  for _, value := range pages {
    if value != nil { GlobalPageArena.ReturnPage(value) }
  }
}

func (s *PageStore) ReleasePages() {
  if s.single != nil { s.ReleaseSinglePages(s.single) }

  if s.double == nil { return }

  for i, single := range s.double {
    if i > s.pagesUsed - len(s.single) { break }
    if single != nil { s.ReleaseSinglePages(single) }
  }
}

func InitPageStore() *PageStore {
  return &PageStore{
    pagesUsed: 0,
    lastEntryBytesUsed: 0,
  }
}
