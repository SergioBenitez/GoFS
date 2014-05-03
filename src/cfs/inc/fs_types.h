#ifndef _SB_FILE_TYPES_H
#define _SB_FILE_TYPES_H

#include <time.h>
#include <sys/types.h>
#include <stdint.h>

// Array Size Maximums
#define MAX_BLOCKS  256
#define MAX_FDS     512

typedef int FileDescriptor;
typedef void Directory; // not sure what it is yet

typedef enum FILE_STATUS_E {
  F_OPEN,
  F_CLOSED,
} FILE_STATUS;

typedef struct Inode_T {
  char *blocks[MAX_BLOCKS];

  time_t mod_time;
  time_t access_time;
  time_t create_time;

  int link_count;
  int file_count;
} Inode;

typedef struct FileHandle_T {
  Inode *inode;
  int seek;
  FILE_STATUS status;
} FileHandle;

typedef struct Process_T {
  Directory *cwd;
  FileHandle *fd_table[MAX_FDS];
  uint16_t free_fds[MAX_FDS];
  int next_fd;
} Process;

#endif
