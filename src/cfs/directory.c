#include <stdlib.h>
#include <string.h>
#include "inc/directory.h"

#define UNUSED(x) (void)(x) 

void
directory_insert(Directory *dir, char *name, void *item) {
  UNUSED(dir);
  UNUSED(name);
  UNUSED(item);
}

void 
directory_remove(Directory *dir, char *name) {
  UNUSED(dir);
  UNUSED(name);
}

Inode *
directory_get(Directory *dir, char *name) {
  UNUSED(dir);
  UNUSED(name);
  return NULL;
}

DirectoryEntry
new_directory_entry(const char *name, void *inode) {
  DirectoryEntry entry;

  size_t name_len = strlen(name) + 1; // strlen doesn't include '\0'
  entry.name = (char *)malloc(name_len);
  strncpy(entry.name, name, name_len);
  
  entry.inode = (Inode *)inode;
  return entry;
}

Directory *
new_directory(Directory *parent) {
  Directory *dir = (Directory *)malloc(sizeof(Directory));
  dir->entries[0] = new_directory_entry("..", parent);
  dir->entries[1] = new_directory_entry(".", dir);
  return dir;
}
