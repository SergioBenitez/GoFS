#ifndef _SB_FILE_H
#define _SB_FILE_H

#include <stddef.h>
#include <sys/types.h>
#include "fs_types.h"

// Access Flags
#define O_RDONLY    (1 << 0)
#define O_WRONLY    (1 << 1)
#define O_RDWR      (1 << 2)
#define O_NONBLOCK  (1 << 3)
#define O_APPEND    (1 << 4)
#define O_CREAT     (1 << 5)
#define O_TRUNC     (1 << 6)
#define O_EXCL      (1 << 7)
#define O_SHLOCK    (1 << 8)
#define O_EXLOCK    (1 << 9)
#define O_NOFOLLOW  (1 << 10)
#define O_SYMLINK   (1 << 11)
#define O_EVTONLY   (1 << 12)
#define O_CLOEXEC   (1 << 13)

size_t file_read(FileHandle *, void *dst, size_t);
size_t file_write(FileHandle *, const void *src, size_t);
off_t file_seek(FileHandle *, off_t, int whence);

Inode *new_inode();
FileHandle *new_handle(Inode *);

void delete_inode(Inode *);
void delete_handle(FileHandle *);

#endif // _SB_FILE_H
