# runzip

A command line utility for unpacking .rar files.

```text
USAGE
	runzip <archive.rar> [./dst/]

EXAMPLES
	runzip ./archive.rar                 # ./inner-dir/
	runzip ./archive.rar ./existing-dir/ # ./existing-dir/inner-dir/
	runzip ./archive.rar ./new-dir/      # ./new-dir/
```

For archives with a single file or folder, this will extract that to the given directory.

For archives with multiple files or folders, it will create a directory of the same name as the archive.
