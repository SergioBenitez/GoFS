#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <time.h>
#include "benchmark.h"

#define UNUSED(x) (void)(x) 

const int NUM = 100;

int
ceil_div(int x, int y) {
  return (x + y - 1) / y;
}

unsigned char *
rand_bytes(Benchmark *b, size_t n) {
  bench_pause(b);
  srand(time(NULL));

  unsigned char *bytes = (unsigned char *)malloc(n);
  for (size_t i = 0; i < n; ++i) {
    bytes[i] = rand() & 0xFF; 
  }

  bench_resume(b);
  return bytes;
}

char *
init_filename(int n, int pfx_len, int end_len, char *pfix, char *end) {
  int len = ceil_div(n, 26);
  char *filename = (char *)malloc(pfx_len + len + end_len + 1);

  // We choose '@' as a placeholder since '@' + 1 = 'A'
  for (int i = 0; i < pfx_len; ++i) filename[i] = pfix[i];
  for (int i = 0; i < len; ++i) filename[pfx_len + i] = '@';
  for (int i = 0; i < end_len; ++i) filename[pfx_len + len + i] = end[i];
  filename[pfx_len + len + end_len] = '\0';
  
  return filename;
}

/**
 * A new FILE*[] is allocated and returned. Caller must free it.
 */
int *
open_many_c(Benchmark *b, int n, int (f) (int, char *)) {
  // pausing for filename and fd array allocation
  bench_pause(b);

  // creating initial filename and fd array
  int pfx_len = 9; // length of /dev/shm/ without null char
  char *filename = init_filename(n, pfx_len, 4, "/dev/shm/", ".out");
  int *fds = (int *)malloc(n * sizeof(int));
  
  // Done with allocations
  bench_resume(b);

  for (int i = 0; i < n; ++i) {
    filename[pfx_len + (i / 26)] += 1;
    fds[i] = open(filename, O_CREAT | O_RDWR, S_IRWXU | S_IRGRP | S_IROTH);
    if (f) f(fds[i], filename);
  }

  // deallocating filenames array
  bench_pause(b);
  free(filename);
  bench_resume(b);

  return fds;
}

void
unlink_all(Benchmark *b, int n) {
  // pausing for filename array allocation
  bench_pause(b);
  char *filename = init_filename(n, 9, 4, "/dev/shm/", ".out");
  bench_resume(b);

  // unlinking
  for (int i = 0; i < n; ++i) {
    filename[9 + (i / 26)] += 1;
    unlink(filename);
  }
  
  // deallocating filenames array
  bench_pause(b);
  free(filename);
  bench_resume(b);
}

int *
open_many(Benchmark *b, int n) {
  return open_many_c(b, n, NULL);
}

void
close_all(int *fds, int n) {
  for (int i = 0; i < n; ++i) close(fds[i]);
}

int
help_close(int fd, char *name) {
  UNUSED(name);
  return close(fd);
}

int
help_unlink(int fd, char *name) {
  UNUSED(fd);
  return unlink(name);
}

int
help_close_unlink(int fd, char *name) {
  help_close(fd, name);
  return help_unlink(fd, name);
}

void
OCSingle(Benchmark *b) {
  UNUSED(b);
  int fd = open("/dev/shm/test.out", O_CREAT | O_RDWR,
      S_IRWXU | S_IRGRP | S_IROTH);
  close(fd);
}

void
OtC(Benchmark *b) {
  int *fds = open_many(b, NUM);
  close_all(fds, NUM);
  free(fds);
}

void
OC(Benchmark *b) {
  int *fds = open_many_c(b, NUM, help_close);
  free(fds);
}

void
OtCtU(Benchmark *b) {
  int *fds = open_many(b, NUM);
  close_all(fds, NUM);
  unlink_all(b, NUM);
  free(fds);
}

void
OCU(Benchmark *b) {
  int *fds = open_many_c(b, NUM, help_close_unlink);
  free(fds);
}

void
OWsC(Benchmark *b) {
  const size_t size = 1024;
  unsigned char *content = rand_bytes(b, size);

  int do_it(int fd, char *name) {
    write(fd, content, size);
    return help_close(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

void
OWsCU(Benchmark *b) {
  const size_t size = 1024;
  unsigned char *content = rand_bytes(b, size);

  int do_it(int fd, char *name) {
    write(fd, content, size);
    help_close(fd, name);
    return help_unlink(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

void
OWbC(Benchmark *b) {
  const size_t size = 40960;
  unsigned char *content = rand_bytes(b, size);

  int do_it(int fd, char *name) {
    write(fd, content, size);
    return help_close(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

void
OWbCU(Benchmark *b) {
  const size_t size = 40960;
  unsigned char *content = rand_bytes(b, size);

  int do_it(int fd, char *name) {
    write(fd, content, size);
    help_close(fd, name);
    return help_unlink(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

void
OWMsC(Benchmark *b) {
  const size_t size = 1024;
  const size_t many = 4096;
  unsigned char *content = rand_bytes(b, size);

  int do_it(int fd, char *name) {
    for (size_t i = 0; i < many; ++i) {
      write(fd, content, size);
    }
    return help_close(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

void
OWMsCU(Benchmark *b) {
  const size_t size = 1024;
  const size_t many = 4096;
  unsigned char *content = rand_bytes(b, size);

  int do_it(int fd, char *name) {
    for (size_t i = 0; i < many; ++i) {
      write(fd, content, size);
    }
    help_close(fd, name);
    return help_unlink(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

void
OWMbC(Benchmark *b) {
  const size_t size = 1048576;
  const size_t many = 32;
  unsigned char *content = rand_bytes(b, size);

  int do_it(int fd, char *name) {
    for (size_t i = 0; i < many; ++i) {
      write(fd, content, size);
    }
    return help_close(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

void
OWMbCU(Benchmark *b) {
  const size_t size = 1048576;
  const size_t many = 32;
  unsigned char *content = rand_bytes(b, size);

  int do_it(int fd, char *name) {
    for (size_t i = 0; i < many; ++i) {
      write(fd, content, size);
    }
    help_close(fd, name);
    return help_unlink(fd, name);
  }

  int *fd = open_many_c(b, NUM, do_it);
  free(fd);
}

void
OWMbbC(Benchmark *b) {
  const size_t start_size = 2;
  const size_t many = 4096;
  unsigned char *content = rand_bytes(b, start_size * many);

  int do_it(int fd, char *name) {
    for (size_t i = 1; i <= many; ++i) {
      write(fd, content, i * start_size);
    }
    return help_close(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

void
OWMbbCU(Benchmark *b) {
  const size_t start_size = 2;
  const size_t many = 4096;
  unsigned char *content = rand_bytes(b, start_size * many);

  int do_it(int fd, char *name) {
    for (size_t i = 1; i <= many; ++i) {
      write(fd, content, i * start_size);
    }
    help_close(fd, name);
    return help_unlink(fd, name);
  }

  int *fds = open_many_c(b, NUM, do_it);
  free(fds);
}

// Clears the /dev/shm/ directory
void
cleanup() {
  system("exec rm -r /dev/shm/*.out");
}

int main() {
  benchmark("Open-Close-Single", OCSingle, cleanup, 2);
  benchmark("Open-Close", OtC, cleanup, 2);
  benchmark("OpenAndClose", OC, cleanup, 2);
  benchmark("Open-Close-Unlink", OtCtU, NULL, 4);
  benchmark("OpenAndCloseAndUnlink", OCU, NULL, 4);
  benchmark("OpenWriteSmallClose", OWsC, cleanup, 2);
  benchmark("OpenWriteSmallCloseUnlink", OWsCU, NULL, 4);
  benchmark("OpenWriteBigClose", OWbC, cleanup, 4);
  benchmark("OpenWriteBigCloseUnlink", OWbCU, NULL, 5);
  benchmark("OpenWriteManySmallClose", OWMsC, cleanup, 3);
  benchmark("OpenWriteManySmallCloseUnlink", OWMsCU, NULL, 5);
  benchmark("OpenWriteManyBigClose", OWMbC, cleanup, 5);
  benchmark("OpenWriteManyBigCloseUnlink", OWMbCU, NULL, 6);
  benchmark("OpenWriteManyBiggerClose", OWMbbC, cleanup, 5);
  benchmark("OpenWriteManyBiggerCloseUnlink", OWMbbCU, NULL, 6);
}
