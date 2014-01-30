package gofs

import (
  "testing"
)

func TestSplitPath(t *testing.T) {
  paths := []string{
    "hello/world",
    "hello.world/wordofme",
    "meyou/me.txt",
    "meyou/me.txt/",
    "hello",
    "/hello/there",
    "/",
    "",
    "../test",
    "./../test/me.txt",
  }

  splits := []string{
    "hello", "world",
    "hello.world", "wordofme",
    "meyou", "me.txt",
    "meyou/me.txt", "",
    "", "hello",
    "/hello", "there",
    "", "",
    "", "",
    "..", "test",
    "./../test", "me.txt",
  }

  for i, path := range paths {
    exp0, exp1 := splits[i * 2], splits[i * 2 + 1]
    actual0, actual1 := splitPath(path)
    if exp0 != actual0 { t.Error("Expected:", exp0, "Got:", actual0) }
    if exp1 != actual1 { t.Error("Expected:", exp1, "Got:", actual1) }
  }
}
