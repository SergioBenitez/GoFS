#include <stdlib.h>
#include <time.h>
#include <string.h>
#include "inc/file.h"

uint8_t *allocate_page(void);
void return_page(uint8_t *);
static inline uint8_t **get_block(Inode *, int num);
size_t inode_read_write(Inode *, const void *, void *, off_t, size_t);

uint8_t *
allocate_page() {
  return (uint8_t *)malloc(PAGE_SIZE);
}

void
return_page(uint8_t *page) {
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
  inode->link_count = 0;
  inode->file_count = 0;
  inode->size = 0;

  time_t now = time(NULL);
  inode->create_time = now;
  inode->access_time = now;
  inode->mod_time = now;
  return inode;
}

void
delete_inode_if_needed(Inode *inode) {
  if (inode->file_count == 0 && inode->link_count == 0)
    delete_inode(inode);
}

void
delete_inode(Inode *inode) {
  // Freeing used blocks
  for (int i = 0; i < MAX_BLOCKS; ++i) {
    uint8_t *block = inode->blocks[i];
    if (block != NULL) return_page(block);
  }

  free(inode);
}

void
inode_dec_file_ref(Inode *inode) {
  inode->file_count--;
  delete_inode_if_needed(inode);
}

void
inode_inc_file_ref(Inode *inode) {
  inode->file_count++;
}

void
inode_dec_link_ref(Inode *inode) {
  inode->link_count--;
  delete_inode_if_needed(inode);
}

void
inode_inc_link_ref(Inode *inode) {
  inode->link_count++;
}

/**
 * Both reads and writes from the inode. 
 * Exactly one of src or dst must be valid.
 *
 * Reads from inode to dst if dst != NULL
 * Writes to inode from src if src != NULL
 */
size_t
inode_read_write(Inode *inode, const void *src, void *dst, off_t o, size_t n) {
  if ((src == NULL && dst == NULL) || (src != NULL && dst != NULL))
    panic("Exactly one of src or dst must be valid!");

  // here, 'act' is a pseudonym for reading or writing to blocks
  // if src != NULL, we're reading from src and [writing] to blocks
  // if dst != NULL, we're writing to dst and [reading] from blocks
  size_t bytes_acted_on = 0;
  int start = o / PAGE_SIZE; // first block to act on
  off_t block_offset = o % PAGE_SIZE; // offset from first block
  int blocks_to_act_on = ceil_div(block_offset + n, PAGE_SIZE);
  for (int i = 0; i < blocks_to_act_on; ++i) {
    // Resetting the block offset after first pass since we want to read from
    // the beginning of the block after the first time.
    if (block_offset != 0 && i > 0) block_offset = 0;

    // Finding our block, adding offset, allocating on write if necessary
    uint8_t **block_ptr = get_block(inode, start + i);
    if (src != NULL && *block_ptr == NULL) *block_ptr = allocate_page();
    if (dst != NULL && *block_ptr == NULL) panic("Reading nonexisting data!");
    uint8_t *block = *block_ptr + block_offset;

    // Figuring out how many bytes (num_bytes) to read from / write to the block
    // Need to account for offsets from first and last blocks
    size_t num_bytes;
    if (i == blocks_to_act_on - 1) num_bytes = n - bytes_acted_on;
    else num_bytes = PAGE_SIZE - block_offset;

    // if src != NULL, then writing to block, else reading from block
    if (src != NULL) {
      memcpy(block, src, num_bytes);
      src += num_bytes;
    } else {
      memcpy(dst, block, num_bytes);
      dst += num_bytes;
    }

    bytes_acted_on += num_bytes;
  }

  return bytes_acted_on;
}

size_t
inode_read(Inode *inode, void *dst, off_t offset, size_t n) {
  if (dst == NULL || offset >= (off_t)inode_size(inode)) return 0;

  return inode_read_write(inode, NULL, dst, offset, n);
}

size_t
inode_write(Inode *inode, const void *src, off_t offset, size_t n) {
  if (src == NULL) return 0;

  size_t written = inode_read_write(inode, src, NULL, offset, n);
  if (offset + written > inode->size)
    inode->size = offset + written;

  return written;
}

size_t
inode_size(Inode *inode) {
  return inode->size;
}
