#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "inc/proc.h"
#include "inc/file.h"
#include "inc/directory.h"

#define UNUSED(x) (void)(x) 

FileDescriptor get_fd(Process *);
void return_fd(Process *, FileDescriptor);

FileDescriptor
get_fd(Process *p) {
  if (p->next_fd >= MAX_FDS) panic("Out of FDs!");

  FileDescriptor next = p->free_fds[p->next_fd];
  FileDescriptor fd = (next == 0) ? p->next_fd : next;
  p->next_fd++;

  return fd;
}

void
return_fd(Process *p, FileDescriptor fd) {
  if (p->next_fd <= START_FD) panic("Over-freeing FDS!");
  p->free_fds[--p->next_fd] = fd;
}

/*
 * Does all of the open operation with the exception of assigning a file
 * descriptor to the file object.
 */
FileHandle *
open_file(Process *p, const char *path, uint32_t flags) {
  /*
   * Something like this needs to go here:
   *
   * Directory *dir = resolvePath(p, path);
   * char *filename = basename(path);
   * Inode *inode = directory_get(dir, filename);
   *
   * Otherwise, we only have 1-level directories.
   * FIXME: Need to check what type of inode was returned when we add multilevel
   * directory support.
   */

  Inode *inode = directory_get(p->cwd, path);
  if (inode == NULL && flags & O_CREAT) {
    inode = new_inode();
    directory_insert(p->cwd, path, inode);
  }

  if (inode == NULL)
    panic("No O_CREAT flag and file not found.");

  return new_handle(inode);
}

FileDescriptor
open(Process *p, const char *path, uint32_t flags) {
  FileDescriptor fd = get_fd(p);
  p->fd_table[fd] = open_file(p, path, flags);
  return fd;
}

int
close(Process *p, FileDescriptor fd) {
  FileHandle *handle = p->fd_table[fd];
  p->fd_table[fd] = NULL;
  return_fd(p, fd);
  delete_handle(handle);
  return 0;
}

int
unlink(Process *p, const char *path) {
  // Again, as in open_file, need to resolve path for multi-level directories
  // Also need to deal with inode reference counts.
  return directory_remove(p->cwd, path);
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

off_t
seek(Process *p, FileDescriptor fd, off_t offset, int whence) {
  FileHandle *handle = p->fd_table[fd];
  return file_seek(handle, offset, whence);
}

Process *
new_process() {
  Process *proc = (Process *)malloc(sizeof(Process));
  memset(proc->fd_table, 0, MAX_FDS * sizeof(FileHandle *));
  memset(proc->free_fds, 0, MAX_FDS * sizeof(uint16_t));

  proc->next_fd = START_FD; // 0, 1, 2 are taken
  proc->cwd = new_directory(NULL); // FIXME: Need global dir.
  return proc;
}

/* int */
/* main() { */
/*   Process *p = new_process(); */
/*   for (int i = 0; i < 1e6; i++) { */
/*     FileDescriptor fd = open(p, "myfile", O_CREAT); */
/*     close(p, fd); */
/*   } */
/*   directory_print(p->cwd); */
/*   unlink(p, "myfile"); */
/* } */
