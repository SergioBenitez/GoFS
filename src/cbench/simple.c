#include <stdio.h>

int main() {
  FILE *file = fopen("test", "wb");
  fseek(file, 100, SEEK_SET);
  fputs("hello, world!", file);
  fclose(file);
}
