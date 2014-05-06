#include <stdlib.h>
#include <string.h>
#include "inc/directory.h"

DirectoryEntry new_directory_entry(const char *, void *);
DirectoryEntry *directory_get_entry(Directory *, const char *);
void directory_print_off(Directory *, off_t);
void repeat_char(const char c, int times);

DirectoryEntry
new_directory_entry(const char *name, void *inode) {
  DirectoryEntry entry;

  size_t name_len = strlen(name) + 1; // strlen doesn't include '\0'
  entry.name = (char *)malloc(name_len);
  strncpy(entry.name, name, name_len);
  
  entry.inode = (Inode *)inode;
  return entry;
}

/**
 * Finds the entry with name 'name'. If 'name' == NULL, then returns the first
 * entry where 'name' == NULL indicating a free slot.
 */
DirectoryEntry *
directory_get_entry(Directory *dir, const char *name) {
  for (int i = 0; i < MAX_DIR_ENTRIES; ++i) {
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

int
directory_insert(Directory *dir, const char *name, void *item) {
  // Checking to see if entry with 'name' already exists
  DirectoryEntry *old_entry = directory_get_entry(dir, name);
  if (old_entry != NULL) return -1;

  // Finding a free entry
  DirectoryEntry *free_entry = directory_get_entry(dir, NULL);
  if (free_entry == NULL) panic("No free entries: directory is full!");

  DirectoryEntry new_entry = new_directory_entry(name, item);
  *free_entry = new_entry;
  return 0;
}

int 
directory_remove(Directory *dir, const char *name) {
  DirectoryEntry *entry = directory_get_entry(dir, name);
  if (entry == NULL) return -1;

  free(entry->name);
  entry->name = NULL;
  entry->inode = NULL;

  return 0;
}

Inode *
directory_get(Directory *dir, const char *name) {
  DirectoryEntry *entry = directory_get_entry(dir, name);
  return (entry == NULL) ? NULL : entry->inode;
}

Directory *
new_directory(Directory *parent) {
  Directory *dir = (Directory *)malloc(sizeof(Directory));
  memset(dir->entries, 0, MAX_DIR_ENTRIES * sizeof(DirectoryEntry));

  Directory *parentDir = (parent == NULL) ? dir : parent;
  dir->entries[0] = new_directory_entry("..", parentDir);
  dir->entries[1] = new_directory_entry(".", dir);
  dir->type = F_DIRECTORY;
  return dir;
}

void
delete_directory(Directory *dir) {
  if (directory_entry_count(dir) > 2)
    panic("Cannot delete directory: directory is not empty.");

  free(dir->entries[0].name);
  free(dir->entries[1].name);
  free(dir);
}

size_t
directory_entry_count(Directory *dir) {
  size_t size = 0;
  for (int i = 0; i < MAX_DIR_ENTRIES; ++i) {
    DirectoryEntry *entry = &dir->entries[i];
    if (entry->name != NULL) size++;
  }
  return size;
}

void
repeat_char(const char c, int times) {
  for (int i = 0; i < times; ++i) putchar(c);
}

/** Just for debugging. Prints the structure of the directory adding 'offset'
 * number of '\t' characters at the beginning of each line of input so that
 * directories printed recursively appear nested in the output.
 *
 * 2 directories, 1 files
 * [1] ..
 * [1] .
 * [0] myfile
*/
void
directory_print_off(Directory *dir, off_t offset) {
  int num_files = 0;
  int num_dirs = 0;
  for (int i = 0; i < MAX_DIR_ENTRIES; ++i) {
    DirectoryEntry entry = dir->entries[i];
    if (entry.name != NULL) {
      if (entry.inode->type == F_DATA) num_files++;
      if (entry.inode->type == F_DIRECTORY) num_dirs++;
    }
  }

  repeat_char('\t', offset);
  printf("%d directories, %d files\n", num_dirs, num_files);
  for (int i = 0; i < MAX_DIR_ENTRIES; ++i) {
    DirectoryEntry entry = dir->entries[i];
    if (entry.name != NULL) {
      repeat_char('\t', offset);
      printf("[%d] %s\n", entry.inode->type, entry.name);

      // Print directories besides '..' and '.' recursively.
      if (entry.inode->type == F_DIRECTORY && i > 1)
        directory_print_off((Directory *)entry.inode, offset + 1);
    }
  }
}

void
directory_print(Directory *dir) {
  directory_print_off(dir, 0);
}
