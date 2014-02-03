package dstore

type DataStore interface {
  // Read/Write beginning at offset o from/to p
  Read(o int, p []byte) (int, error)
  Write(o int, p []byte) (int, error)

  // Returns the number of bytes stored
  Size() int
}
