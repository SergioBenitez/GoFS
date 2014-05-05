#include <stdlib.h>
#include <time.h>
#include "inc/file.h"

#define UNUSED(x) (void)(x) 

FileHandle *
new_handle(Inode *inode) {
  FileHandle *handle = (FileHandle *)malloc(sizeof(FileHandle));

  inode_inc_file_ref(inode);
  handle->inode = inode;
  handle->status = F_OPEN;
  handle->seek = 0;

  return handle;
}

void
delete_handle(FileHandle *handle) {
  inode_dec_file_ref(handle->inode);
  free(handle);
}

size_t
file_read(FileHandle *handle, void *dst, size_t num) {
  size_t read = inode_read(handle->inode, dst, handle->seek, num);
  handle->seek += read;
  return read;
}

size_t
file_write(FileHandle *handle, const void *src, size_t num) {
  size_t written = inode_write(handle->inode, src, handle->seek, num);
  handle->seek += written;
  return written;
}

off_t
file_seek(FileHandle *handle, off_t offset, int whence) {
  switch (whence) {
    case SEEK_SET:
      handle->seek = offset;
      break;
    case SEEK_CUR:
      handle->seek += offset;
      break;
    case SEEK_END:
      handle->seek = inode_size(handle->inode) - offset;
      break;
  }
  return handle->seek;
}
