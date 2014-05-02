#ifndef _SB_FILE_H
#define _SB_FILE_H

#include <stdint.h>
#include <stddef.h>

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

// Array Size Constants
#define MAX_BLOCKS  256
#define MAX_FDS     512

typedef int FileDescriptor;

typedef enum FILE_STATUS_E {
  F_OPEN,
  F_CLOSED,
} FILE_STATUS;

typedef struct Inode_T {
  char *blocks[MAX_BLOCKS];

  double mod_time;
  double access_time;
  double create_time;

  int link_count;
  int file_count;
} Inode;

typedef struct Process_T {
  Inode *file_descriptor_table[MAX_FDS];
} Process;

typedef struct FileHandle_T {
  Inode *inode;
  int seek;
  FILE_STATUS status;
} FileHandle;

FileDescriptor open(Process *, const char *path, uint32_t flags);
size_t read(Process *, FileDescriptor, void *dst, size_t);
size_t write(Process *, FileDescriptor, const void *src, size_t);
int close(FileDescriptor);
int unlink(const char *path);

#endif // _SB_FILE_H
