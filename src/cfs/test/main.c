#include <stdlib.h>
#include <check.h>
#include "../inc/proc.h"
#include "../inc/file.h"

Process *p;

void
setup() {
  p = new_process();
}

void
teardown() {
  // process_delete(p)?
}

unsigned char *
rand_bytes(size_t n) {
  srand(time(NULL));

  unsigned char *bytes = (unsigned char *)malloc(n);
  for (size_t i = 0; i < n; ++i) {
    bytes[i] = rand() & 0xFF; 
  }

  return bytes;
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

  off_t seek_point = 323;
  seek(p, fd, seek_point, SEEK_SET);
  read(p, fd, buf, size - seek_point);
  ck_assert(!memcmp(data + seek_point, buf, size - seek_point));

  unlink(p, "file");
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

Suite *
test_suite() {
  Suite *s = suite_create("CFS");

  // The 'core' case with setup/teardown fixture
  TCase *tc_core = tcase_create("Core");
  tcase_add_checked_fixture(tc_core, setup, teardown); // checked = once/test

  // Adding tests to case 'tc_core'
  tcase_add_test(tc_core, write_small_read);
  tcase_add_test(tc_core, write_large_read);

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
