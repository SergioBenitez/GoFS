#include <stdlib.h>
#include <time.h>
#include <string.h>
#include "inc/file.h"

uint8_t *allocate_page(void);
void return_page(char *);
static inline uint8_t **get_block(Inode *, int num);

uint8_t *
allocate_page() {
  return (uint8_t *)malloc(PAGE_SIZE);
}

void
return_page(char *page) {
  free(page); 
}

static inline uint8_t **
get_block(Inode *inode, int num) {
  if (num >= MAX_BLOCKS) panic("Exceeding file size.");
  return &inode->blocks[num];
}

Inode *
new_inode() {
  Inode *inode = (Inode *)malloc(sizeof(Inode));
  memset(inode->blocks, 0, MAX_BLOCKS * sizeof(uint8_t *));

  inode->type = F_DATA;
  inode->link_count = 1;
  inode->file_count = 0;
  inode->size = 0;

  time_t now = time(NULL);
  inode->create_time = now;
  inode->access_time = now;
  inode->mod_time = now;
  return inode;
}

void
delete_inode(Inode *inode) {
  // Need to account for reference counting
  free(inode);
}


size_t
inode_read(Inode *inode, void *dst, off_t offset, size_t n) {
  /* printf("Reading %zu bytes to %p at %lld.\n", n, dst, offset); */
  /* printf("Inode is %zu bytes.\n", inode_size(inode)); */
  if (dst == NULL || offset >= (off_t)inode_size(inode)) return 0;

  size_t read = 0;
  off_t block_offset = offset % PAGE_SIZE;
  int start = offset / PAGE_SIZE;
  int blocks_to_read = ceil_div(block_offset + n, PAGE_SIZE);
  for (int i = 0; i < blocks_to_read; ++i) {
    if (block_offset != 0 && i > 0) block_offset = 0;

    uint8_t *block = *get_block(inode, start + i) + block_offset;
    if (block == NULL) panic("Reading where no data exists!");

    size_t bytes_to_read;
    if (i == blocks_to_read - 1) bytes_to_read = n - read;
    else bytes_to_read = PAGE_SIZE - block_offset;

    memcpy(dst, block, bytes_to_read);
    dst += bytes_to_read;
    read += bytes_to_read;
  }

  /* printf("Read %zu bytes.\n", read); */
  return read;
}

size_t
inode_write(Inode *inode, const void *src, off_t offset, size_t n) {
  if (src == NULL) return 0;
  /* printf("Writing %zu bytes from %p at %lld.\n", n, src, offset); */

  size_t written = 0;
  off_t block_offset = offset % PAGE_SIZE;
  int start = offset / PAGE_SIZE;
  int blocks_to_write = ceil_div(block_offset + n, PAGE_SIZE);
  for (int i = 0; i < blocks_to_write; ++i) {
    if (block_offset != 0 && i > 0) block_offset = 0;

    uint8_t **block_ptr = get_block(inode, start + i);
    /* printf("Block starts at %p. || ", *block_ptr); */
    if (*block_ptr == NULL) *block_ptr = allocate_page();
    uint8_t *block = *block_ptr + block_offset;
    /* printf("Block starts at %p.\n", *block_ptr); */

    size_t bytes_to_write;
    if (i == blocks_to_write - 1) bytes_to_write = n - written;
    else bytes_to_write = PAGE_SIZE - block_offset;

    // printf("Copying %zu bytes to %p from %p\n", bytes_to_write, block, src);
    memcpy(block, src, bytes_to_write);
    src += bytes_to_write;
    written += bytes_to_write;
  }

  /* printf("Wrote %zu bytes.\n", written); */
  if (offset + n > inode->size) inode->size = offset + n;
  return written;
}

size_t
inode_size(Inode *inode) {
  return inode->size;
}
