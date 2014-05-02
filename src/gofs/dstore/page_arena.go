package dstore

// import "fmt"

var GlobalPageArena *PageArena

const PAGE_SIZE = 4096
const EXP_GROW_LIMIT = 65536 // 131,072 pages = 256MB
// const EXP_GROW_LIMIT = 131072 // 131,072 pages = 512MB
// const EXP_GROW_LIMIT = 262144 // 262,144 pages = 1GB
// const EXP_GROW_LIMIT = 1048576 // 4GB

type PageArena struct {
  alloc int       // number of allocated pages
  size int        // size in num of pages
  pages [][]byte
}

func allocatePages(num int) [][]byte {
  // Allocating containing array
  pages := make([][]byte, num, num)

  // Allocating all bytes at once
  allPages := make([]byte, num * PAGE_SIZE, num * PAGE_SIZE)

  // Setting the internal pointers to all num PAGE_SIZE slices
  for i := 0; i < num; i++ {
    pages[i] = allPages[i * PAGE_SIZE : (i + 1) * PAGE_SIZE]
  }

  return pages
}

/*
* Grows the page arena size in an interesting way.
* The arena is exponentially grown, doubling in size each time, until
* EXP_GROW_LIMIT pages have been allocated. From that point, only EXP_GROW_LIMIT
* pages are added (so first time after EXP_GROW_LIMIT it's doubled, then only
* 1.5x, then 1.25x, etc.).
*/
func (a *PageArena) grow() {
  var newSize int
  if a.size < EXP_GROW_LIMIT { newSize = a.size * 2
  } else { newSize = a.size + EXP_GROW_LIMIT }

  newPages := allocatePages(newSize - a.size)
  a.pages = append(a.pages, newPages...)
  a.size = newSize
}

// NOTE! Page is not guaranteed to be zeroed!
func (a *PageArena) AllocatePage() []byte {
  // fmt.Println("Allocating page. Pages so far:", a.alloc)

  if a.alloc >= a.size { a.grow() }
  if a.alloc >= a.size { panic("Out of memory @ pageArena!") }

  page := a.pages[a.alloc]
  a.alloc += 1
  return page
}

func (a *PageArena) ReturnPage(page []byte) {
  if a.alloc <= 0 { panic("Over-freeing pages!") }

  a.alloc -= 1
  a.pages[a.alloc] = page
}

func InitPageArena(size int) *PageArena {
  // fmt.Println("New arena with size", size)

  arena := &PageArena{
    alloc: 0,
    size: size,
    pages: allocatePages(size),
  }

  return arena
}
