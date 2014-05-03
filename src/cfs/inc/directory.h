#ifndef _SB_DIRECTORY_H
#define _SB_DIRECTORY_H

#include "fs_types.h"

Directory *new_directory(Directory *parent);
void directory_insert(Directory *, char *name, void *);
void directory_remove(Directory *, char *name);
Inode *directory_get(Directory *, char *name);

#endif
