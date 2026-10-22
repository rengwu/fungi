package main

// atomicBytes already flushes the writable file before replacing its name.
// Windows FlushFileBuffers requires a writable handle; os.Open(directory)
// cannot provide a portable equivalent of Unix directory fsync. Do not report
// a failed save after a successful rename by flushing a read-only directory.
// Directory-entry durability across power loss is not promised on Windows.
func syncDir(string) error { return nil }
