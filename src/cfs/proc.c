#include <stdlib.h>
#include <stdio.h>
#include "inc/proc.h"
#include "inc/file.h"

#define UNUSED(x) (void)(x) 

FileDescriptor getNewFD(Process *);
void returnFD(Process *, FileDescriptor);

void
panic(const char *message) {
  fputs(message, stderr);
  exit(1);
}

FileDescriptor
getNewFD(Process *p) {
  if (p->next_fd > MAX_FDS) panic("Out of FDs!");

  FileDescriptor next = p->free_fds[p->next_fd];
  FileDescriptor fd = (next == 0) ? p->next_fd : next;
  p->next_fd++;

  return fd;
}

void
returnFD(Process *p, FileDescriptor fd) {
  if (p->next_fd <= START_FD) panic("Over-freeing FDS!");
  p->free_fds[--p->next_fd] = fd;
}

/*
 * Does all of the open operations with the exception of assigning
 * a file descriptor to the file object.
 */
FileHandle *
openFile(Process *p, const char *path, uint32_t flags) {
  UNUSED(p);
  UNUSED(path);

  /*
   * Something like this needs to go here:
   *
   * Inode *inode = resolvePath(p, path);
   * if (inode != NULL) return newFileHandle(inode);
   *
   * Otherwise, we ALWAYS allocate a new inode.
   */

  if (flags & O_CREAT) {
    Inode *inode = newInode();
    return newFileHandle(inode);
  } else {
    panic("Directories not yet implemented.");
    return NULL;
  }
}

FileDescriptor
open(Process *p, const char *path, uint32_t flags) {
  UNUSED(p);
  UNUSED(path);
  UNUSED(flags);

  FileDescriptor fd = getNewFD(p);
  p->fd_table[fd] = openFile(p, path, flags);
  return fd;
}

int
close(Process *p, FileDescriptor fd) {
  FileHandle *handle = p->fd_table[fd];
  returnFD(p, fd);
  deleteFileHandle(handle);
  return 0;
}

int
unlink(Process *p, const char *path) {
  UNUSED(p);
  UNUSED(path);
  return 0;
}

size_t
read(Process *p, FileDescriptor fd, void *dst, size_t num) {
  FileHandle *handle = p->fd_table[fd];
  return file_read(handle, dst, num);
}

size_t
write(Process *p, FileDescriptor fd, const void *src, size_t num) {
  FileHandle *handle = p->fd_table[fd];
  return file_write(handle, src, num);
}

Process *
newProcess() {
  Process *proc = (Process *)malloc(sizeof(Process));
  proc->next_fd = START_FD; // 0, 1, 2 are taken
  return proc;
}

int
main() {
  Process *p = newProcess();
  for (int i = 0; i < 1e3; i++) {
    FileDescriptor fd = open(p, "myfile", O_CREAT);
    close(p, fd);
  }
  unlink(p, "myfile");
}
