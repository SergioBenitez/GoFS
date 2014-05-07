#include <stdlib.h>
#include <check.h>
#include "../inc/proc.h"
#include "../inc/file.h"

Process *p;

void
setup() {
  srand(time(NULL));
  p = new_process();
}

void
teardown() {
  delete_process(p);
}

unsigned char *
rand_bytes(size_t n) {
  unsigned char *bytes = (unsigned char *)malloc(n);
  for (size_t i = 0; i < n; ++i) {
    bytes[i] = rand() & 0xFF; 
  }

  return bytes;
}

unsigned int
rand_interval(unsigned int min, unsigned int max) {
  unsigned int r;
  const unsigned int range = max - min;
  const unsigned int buckets = RAND_MAX / range;
  const unsigned int limit = buckets * range;

  /* Create equal size buckets all in a row, then fire randomly towards
   * the buckets until you land in one of them. All buckets are equally
   * likely. If you land off the end of the line of buckets, try again. */
  do {
    r = rand();
  } while (r >= limit);

  return min + (r / buckets);
}


START_TEST(write_large_read) {
  size_t size = 4096 * 5 + 92;
  uint8_t *data = rand_bytes(size);
  uint8_t *buf = (uint8_t *)malloc(size);

  FileDescriptor fd = open(p, "file", O_CREAT);
  size_t written = write(p, fd, data, size);
  ck_assert_uint_eq(written, size);

  seek(p, fd, 0, SEEK_SET);
  read(p, fd, buf, size);
  ck_assert(!memcmp(data, buf, size));

  off_t seek_point = 323; // random position to make sure seek + read/write ok
  seek(p, fd, seek_point, SEEK_SET);
  read(p, fd, buf, size - seek_point);
  ck_assert(!memcmp(data + seek_point, buf, size - seek_point));

  close(p, fd);
  unlink(p, "file");
  free(data);
  free(buf);
} END_TEST

START_TEST(write_small_read) {
  char *data = "Hello, world!";
  size_t data_size = strlen(data) + 1;
  char *buf = (char *)malloc(data_size);

  // Creates 'file', writes 'data' to it, reads it back and asserts equality
  FileDescriptor fd = open(p, "file", O_CREAT);
  write(p, fd, data, data_size);
  seek(p, fd, 0, SEEK_SET);
  read(p, fd, buf, data_size);
  ck_assert_str_eq(data, buf);
  close(p, fd);

  // Clearing buf.
  memset(buf, '\0', data_size);

  // Reopenning 'file' to make sure 'data' is still there.
  fd = open(p, "file", O_CREAT);
  read(p, fd, buf, data_size);
  ck_assert_str_eq(data, buf);
  close(p, fd);
  unlink(p, "file");
  
  // Clearing buf.
  memset(buf, '\0', data_size);

  // Making sure memory for file was cleared on unlink
  fd = open(p, "file", O_CREAT);
  size_t num = read(p, fd, buf, data_size);
  ck_assert_uint_eq(num, 0);
  ck_assert_str_eq("", buf);
  close(p, fd);
  unlink(p, "file");

  // Creating new file, closing before ever reading, then opening and reading
  fd = open(p, "file2", O_CREAT);
  write(p, fd, data, data_size);
  close(p, fd);
  
  fd = open(p, "file2", O_CREAT);
  read(p, fd, buf, data_size);
  ck_assert_str_eq(data, buf);
  close(p, fd);
  unlink(p, "file2");

  free(buf);
} END_TEST

START_TEST(link_unlink) {
  char *filename = "file";
  int size = 5000;
  uint8_t *data = rand_bytes(size);
  uint8_t *buf = (uint8_t *)malloc(size);

  FileDescriptor fd = open(p, filename, O_CREAT | O_RDWR);
  size_t bytes_written = write(p, fd, data, size);
  seek(p, fd, 0, SEEK_SET);
  size_t bytes_read = read(p, fd, buf, size);
  ck_assert_uint_eq(bytes_written, size);
  ck_assert_uint_eq(bytes_written, bytes_read);
  ck_assert(!memcmp(data, buf, size));
  close(p, fd);

  // Clearing buf.
  memset(buf, '\0', size);

  // Linking to another filename, making sure contents still there
  char *other_filename = "me2please";
  link(p, filename, other_filename);
  fd = open(p, other_filename, O_CREAT | O_RDWR);
  bytes_read = read(p, fd, buf, size);
  ck_assert_uint_eq(bytes_read, size);
  ck_assert(!memcmp(data, buf, size));
  close(p, fd);
  
  // Clearing buf.
  memset(buf, '\0', size);

  // Unlinking previous one, making sure contents still there
  unlink(p, filename);
  fd = open(p, other_filename, O_CREAT | O_RDWR);
  bytes_read = read(p, fd, buf, size);
  ck_assert_uint_eq(bytes_read, size);
  ck_assert(!memcmp(data, buf, size));
  close(p, fd);

  // All done
  unlink(p, other_filename);
  free(data);
  free(buf);
} END_TEST

void
read_write_large_seek_action(int size, int limit, int seeks) {
  char *filename = "file";
  uint8_t *data = rand_bytes(size);
  uint8_t *buf = (uint8_t *)malloc(size);
  memset(buf, '\0', size);

  FileDescriptor fd = open(p, filename, O_CREAT | O_RDWR);
  size_t bytes_written = write(p, fd, data, size);
  seek(p, fd, 0, SEEK_SET);
  size_t bytes_read = read(p, fd, buf, size);
  ck_assert_uint_eq(bytes_written, size);
  ck_assert_uint_eq(bytes_written, bytes_read);
  ck_assert(!memcmp(data, buf, size));

  // Verify a random number of bytes from a random position 5000 times.
  for (int i = 0; i < seeks; ++i) {
    // Choosing how many bytes to read then zeroing that many bytes in buf
    size_t bytes = rand_interval(0, limit);
    memset(buf, '\0', bytes);

    // Choosing a valid seek position so as to not read past the end, reading
    off_t position = rand_interval(0, size - bytes);
    seek(p, fd, position, SEEK_SET);
    bytes_read = read(p, fd, buf, bytes);
    
    ck_assert_uint_eq(bytes_read, bytes);
    ck_assert(!memcmp(data + position, buf, bytes));
  }

  close(p, fd);
  unlink(p, filename);
  free(data);
  free(buf);
}

START_TEST(read_write_large_seek) {
  read_write_large_seek_action(4096 * 256, 4096 * 128, 5000);
} END_TEST

START_TEST(read_write_really_large_seek) {
  // Writes 257MB, verifies random 50MB chunk 750x
  read_write_large_seek_action(4096 * 256 * 256 + 4096 * 256,
      4096 * 256 * 50, 750);
} END_TEST

// Makes sure a file remains existing until close is called
START_TEST(unlink_before_close) {
  char *filename = "file";
  int size = 5000;
  uint8_t *data = rand_bytes(size);
  uint8_t *buf = (uint8_t *)malloc(size);
  memset(buf, '\0', size);

  FileDescriptor fd = open(p, filename, O_CREAT | O_RDWR);
  unlink(p, filename);

  size_t bytes_written = write(p, fd, data, size);
  seek(p, fd, 0, SEEK_SET);
  size_t bytes_read = read(p, fd, buf, size);

  ck_assert_uint_eq(bytes_written, size);
  ck_assert_uint_eq(bytes_written, bytes_read);
  ck_assert(!memcmp(data, buf, size));

  close(p, fd);
  free(data);
  free(buf);
} END_TEST

Suite *
test_suite() {
  Suite *s = suite_create("CFS");

  // The 'core' case with setup/teardown fixture
  TCase *tc_core = tcase_create("Core");
  tcase_add_checked_fixture(tc_core, setup, teardown); // checked = once/test
  tcase_set_timeout(tc_core, 20);

  // Adding tests to case 'tc_core'
  tcase_add_test(tc_core, write_small_read);
  tcase_add_test(tc_core, write_large_read);
  tcase_add_test(tc_core, link_unlink);
  tcase_add_test(tc_core, unlink_before_close);
  tcase_add_test(tc_core, read_write_large_seek);
  tcase_add_test(tc_core, read_write_really_large_seek);

  // Adding case to suite
  suite_add_tcase(s, tc_core);
  return s;
}

int
main() {
  Suite *s = test_suite();
  SRunner *sr = srunner_create(s);

  srunner_run_all(sr, CK_NORMAL);
  int number_failed = srunner_ntests_failed(sr);
  srunner_free(sr);

  return (number_failed == 0) ? EXIT_SUCCESS : EXIT_FAILURE;
}
