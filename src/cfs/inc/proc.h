#ifndef _SB_PROC_H
#define _SB_PROC_H

#include <stdint.h>
#include <sys/types.h>
#include "fs_types.h"

// The file descriptor to start with
#define START_FD 3

Process *new_process();

FileDescriptor open(Process *, const char *path, uint32_t flags);
size_t read(Process *, FileDescriptor, void *dst, size_t);
size_t write(Process *, FileDescriptor, const void *src, size_t);
off_t seek(Process *, FileDescriptor, off_t, int whence);
int close(Process *, FileDescriptor);
int unlink(Process *, const char *path);

#endif // _SB_PROC_H
