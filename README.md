# Archive Expander

This tool takes an archive files in the list of types supported and will expand them.  If any archive is within the current archive, it will also be expanded into a folder named with the archives name which is expanded.  By default, recursion is set to 1, but this can be set as large a desired.

The goal of this tool is give one the ability to expand archives without having to extract every level one at a time which they can be all expanded at once.  Take for example, a compressed (gzip), iso file (iso9660), with package files inside.  Each layer would add N x Size, so with 5 levels of recusion, a DVD (with 4GB in size) could become on the order of 40 GB on disk if compression is used and each layer is kept.  This tool, would give you the end file structure without the transident bloat.

```bash
./archive-exploder -output out -input debian-11.3.0-amd64-netinst.iso -r 3
```

```
$ ./archive-exploder -h
Archive Exploder,  Version: 0.1.20220501.1847

Usage: ./archive-exploder [options...]

  -debug
        Turn on debug, more verbose
  -input string
        Path to put the extracted files (default "test.iso")
  -output string
        Path to put the extracted files (default ".")
  -r int
        Levels of recusion (archives-inside-archives) to expand (default 1)

Formats supported:
  - gzip
  - iso9660
```
