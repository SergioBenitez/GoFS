#ifndef _SB_DIRECTORY_H
#define _SB_DIRECTORY_H

#include "fs_types.h"

Directory *new_directory(Directory *parent);
void directory_insert(Directory *, const char *name, void *);
int directory_remove(Directory *, const char *name);
Inode *directory_get(Directory *, const char *name);
void directory_print(Directory *);

#endif
