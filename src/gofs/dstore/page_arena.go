package dstore

import "fmt"

var GlobalPageArena *PageArena

const PAGE_SIZE = 4096

type PageArena struct {
  alloc int
  size int
  pages []*[PAGE_SIZE]byte
}

func (a *PageArena) grow() {
  // Here goes the code to grow the arena!
}

// NOTE! Page is not guaranteed to be zeroed!
func (a *PageArena) AllocatePage() *[PAGE_SIZE]byte {
  fmt.Println("Allocating page. Pages so far:", a.alloc)

  if a.alloc >= a.size { a.grow() }
  if a.alloc >= a.size { panic("Out of memory @ pageArena!") }

  page := a.pages[a.alloc]
  a.alloc += 1
  return page
}

func (a *PageArena) ReturnPage(page *[PAGE_SIZE]byte) {
  if a.alloc <= 0 { panic("Over-freeing pages!") }

  a.alloc -= 1
  a.pages[a.alloc] = page
}

func InitPageArena(size int) *PageArena {
  fmt.Println("New arena with size", size)

  arena := &PageArena{
    alloc: 0,
    size: size,
    pages: make([]*[PAGE_SIZE]byte, size, size * 2),
  }

  // Allocating the first 'size' pages
  for i := 0; i < size; i++ {
    arena.pages[i] = new([PAGE_SIZE]byte)
  }

  return arena
}
