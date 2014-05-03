#ifndef _SB_FILE_TYPES_H
#define _SB_FILE_TYPES_H

#include <time.h>
#include <sys/types.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

/**
 * Directory/Inode/Entry Hackery Note
 *
 * Directory entries point to inodes, but inodes can be various things. Although
 * we don't really deal with the different types of inodes here, we do deal with
 * two types:
 *
 * 1) Inodes that point to data.
 * 2) Inodes that aren't inodes but directories.
 *
 * It's the second type, the directory, that causes us to use a bit of hackery.
 * Instead of actually using an Inode structure to store the contents of a
 * directory, which would require serializing/deserializing the stored data on
 * each directory access, we have a new structure, Directory, which is actually
 * a directory. However, since the directory entries point to Inodes, how do we
 * put a directory inside of a directory?
 *
 * This is where the Inode->type field comes in. We simply put the type field in
 * the same location of both the Inode and Directory structure (in our case, the
 * beginning). Then, when we look at a directory entry, we can always look at
 * the Inode->type field and determine if the structure it's pointing to is
 * actually and Inode or a Directory. If it's a Directory, we can simply cast it
 * as such and be on our way.
 */

// Array Size Maximums
#define MAX_BLOCKS  256
#define MAX_FDS     512
#define MAX_ENTRIES 128

typedef int FileDescriptor;

typedef enum FILE_STATUS_E {
  F_OPEN,
  F_CLOSED,
} FILE_STATUS;

typedef enum FILE_TYPE_E {
  F_DATA,
  F_DIRECTORY,
} FILE_TYPE;

typedef struct Inode_T {
  FILE_TYPE type; // Doing some hackery here. See directory/inode note above.
  char *blocks[MAX_BLOCKS];

  time_t mod_time;
  time_t access_time;
  time_t create_time;

  int link_count;
  int file_count;
} Inode;

typedef struct DirectoryEntry_T {
  char *name;
  Inode *inode; // See directory/inode note above.
} DirectoryEntry;

typedef struct Directory_T {
  FILE_TYPE type; // See directory/inode note above.
  DirectoryEntry entries[MAX_ENTRIES];
} Directory;

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

void inline
panic(const char *message) {
  fputs(message, stderr);
  exit(1);
}

#endif
