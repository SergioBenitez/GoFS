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

START_TEST(write_small_read) {
  char *data = "Hello, world!";
  char buf[13];

  // Clearing buf.
  memset(buf, '\0', 13);

  // Creates 'file', writes 'data' to it, reads it back and asserts equality
  FileDescriptor fd = open(p, "file", O_CREAT);
  write(p, fd, data, 13);
  read(p, fd, buf, 13);
  ck_assert_str_eq(data, buf);
  close(p, fd);

  // Clearing buf.
  memset(buf, '\0', 13);

  // Reopenning 'file' to make sure 'data' is still there.
  fd = open(p, "file", O_CREAT);
  read(p, fd, buf, 13);
  ck_assert_str_eq(data, buf);
  close(p, fd);
  unlink(p, "file");
  
  // Clearing buf.
  memset(buf, '\0', 13);

  // Making sure memory for file was cleared on unlink
  fd = open(p, "file", O_CREAT);
  size_t num = read(p, fd, buf, 13);
  ck_assert_uint_eq(num, 0);
  ck_assert_str_eq("", buf);
  close(p, fd);
  unlink(p, "file");

  // Creating new file, closing before ever reading, then opening and reading
  fd = open(p, "file2", O_CREAT);
  write(p, fd, data, 13);
  close(p, fd);
  
  fd = open(p, "file2", O_CREAT);
  read(p, fd, buf, 13);
  ck_assert_str_eq(data, buf);
  close(p, fd);
  unlink(p, "file2");
} END_TEST

Suite *
test_suite() {
  Suite *s = suite_create("CFS");

  // The 'core' case with setup/teardown fixture
  TCase *tc_core = tcase_create("Core");
  tcase_add_checked_fixture(tc_core, setup, teardown); // checked = once/test

  // Adding tests to case 'tc_core'
  tcase_add_test(tc_core, write_small_read);

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
