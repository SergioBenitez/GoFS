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

// Whence Values
#define SEEK_SET    0
#define SEEK_CUR    1
#define SEEK_END    2

#define PAGE_SIZE   4096

FileHandle *new_handle(Inode *);
void delete_handle(FileHandle *);
size_t file_read(FileHandle *, void *dst, size_t);
size_t file_write(FileHandle *, const void *src, size_t);
off_t file_seek(FileHandle *, off_t, int whence);

Inode *new_inode();
void delete_inode(Inode *);
size_t inode_read(Inode *, void *dst, off_t, size_t);
size_t inode_write(Inode *, const void *src, off_t, size_t);
size_t inode_size(Inode *);

#endif // _SB_FILE_H
