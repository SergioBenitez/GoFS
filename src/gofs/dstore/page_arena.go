package dstore

type PageArena struct {
  alloc int
  size int
  pages []*[4096]byte
}

func (a *PageArena) grow() {
  // Here goes the code to grow the arena!
}

// NOTE! Page is not guaranteed to be zeroed!
func (a *PageArena) AllocatePage() *[4096]byte {
  if a.alloc >= a.size { a.grow() }
  if a.alloc >= a.size { panic("Out of memory @ pageArena!") }

  page := a.pages[a.alloc]
  a.alloc += 1
  return page
}

func (a *PageArena) ReturnPage(page *[4096]byte) {
  if a.alloc <= 0 { panic("Over-freeing pages!") }

  a.alloc -= 1
  a.pages[a.alloc] = page
}

func InitPageArena(size int) *PageArena {
  arena := &PageArena{
    alloc: 0,
    size: size,
    pages: make([]*[4096]byte, size, size * 2),
  }

  // Allocating the first 'size' pages
  for i := 0; i < size; i++ {
    arena.pages[i] = new([4096]byte)
  }

  return arena
}
