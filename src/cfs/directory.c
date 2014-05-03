#include <stdlib.h>
#include <string.h>
#include "inc/directory.h"

#define UNUSED(x) (void)(x) 

DirectoryEntry new_directory_entry(const char *, void *);
DirectoryEntry *directory_get_entry(Directory *, const char *);

DirectoryEntry
new_directory_entry(const char *name, void *inode) {
  DirectoryEntry entry;

  size_t name_len = strlen(name) + 1; // strlen doesn't include '\0'
  entry.name = (char *)malloc(name_len);
  strncpy(entry.name, name, name_len);
  
  entry.inode = (Inode *)inode;
  return entry;
}

DirectoryEntry *
directory_get_entry(Directory *dir, const char *name) {
  for (int i = 0; i < MAX_ENTRIES; ++i) {
    DirectoryEntry *entry = &dir->entries[i];
    if (name == NULL) {
      if (entry->name == NULL) return entry;
    } else {
      if (entry->name == NULL) continue;
      if (!strcmp(entry->name, name)) return entry;
    }
  }
  
  return NULL;
}

void
directory_insert(Directory *dir, const char *name, void *item) {
  DirectoryEntry new_entry = new_directory_entry(name, item);
  DirectoryEntry *entry = directory_get_entry(dir, NULL);

  if (entry == NULL) panic("Directory is full!");
  *entry = new_entry;
}

int 
directory_remove(Directory *dir, const char *name) {
  DirectoryEntry *entry = directory_get_entry(dir, name);
  if (entry == NULL) return -1;

  free(entry->name);
  entry->name = NULL;
  entry->inode = NULL;

  return 0;
  // Decrement the inode link count here?
  // Who's responsibility should it be? Probably not the directories.
}

Inode *
directory_get(Directory *dir, const char *name) {
  DirectoryEntry *entry = directory_get_entry(dir, name);
  return (entry == NULL) ? NULL : entry->inode;
}

Directory *
new_directory(Directory *parent) {
  Directory *dir = (Directory *)malloc(sizeof(Directory));
  Directory *parentDir = (parent == NULL) ? dir : parent;
  dir->entries[0] = new_directory_entry("..", parentDir);
  dir->entries[1] = new_directory_entry(".", dir);
  dir->type = F_DIRECTORY;
  return dir;
}

void
directory_print(Directory *dir) {
  int num_files = 0;
  int num_dirs = 0;
  for (int i = 0; i < MAX_ENTRIES; ++i) {
    DirectoryEntry entry = dir->entries[i];
    if (entry.name != NULL) {
       if (entry.inode->type == F_DATA) num_files++;
       if (entry.inode->type == F_DIRECTORY) num_dirs++;
    }
  }

  printf("\nDir Contents: %d directories, %d files\n", num_dirs, num_files);
  for (int i = 0; i < MAX_ENTRIES; ++i) {
    DirectoryEntry entry = dir->entries[i];
    if (entry.name != NULL) {
      printf("%s (type: %d)\n", entry.name, entry.inode->type);
    }
  }
}
