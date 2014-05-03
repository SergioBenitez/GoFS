#include "inc/file.h"
#include <stdlib.h>
#include <time.h>

#define UNUSED(x) (void)(x) 

Inode *
new_inode() {
  Inode *inode = (Inode *)malloc(sizeof(Inode));
  inode->link_count = 1;

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

FileHandle *
new_handle(Inode *inode) {
  FileHandle *handle = (FileHandle *)malloc(sizeof(FileHandle));

  inode->file_count++;
  handle->inode = inode;
  handle->status = F_OPEN;

  return handle;
}

void
delete_handle(FileHandle *handle) {
  // Need to account for reference counting
  free(handle);
}

size_t
file_read(FileHandle *handle, void *dst, size_t num) {
  UNUSED(handle);
  UNUSED(dst);
  UNUSED(num);
  return 0;
}

size_t
file_write(FileHandle *handle, const void *src, size_t num) {
  UNUSED(handle);
  UNUSED(src);
  UNUSED(num);
  return 0;
}
