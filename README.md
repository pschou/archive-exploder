# Archive Exploder

This tool takes an archive files in the list of types supported and will expand
them.  If any archive is within the current archive, it will also be expanded
into a folder named with the archives name which is expanded.  By default,
recursion is set to 1, but this can be set as large a desired.

The goal of this tool is give one the ability to expand archives without having
to extract every level one at a time which they can be all expanded at once.
Take for example, a compressed (gzip), iso file (iso9660), with package files
inside.  Each layer would add N x Size, so with 5 levels of recusion, a DVD
(with 4GB in size) could become on the order of 40 GB on disk if compression is
used and each layer is kept.  This tool, would give you the end file structure
without the transident bloat.

A thing to note is the file extension is not used, but instead the underlying
bytes.  This way a file with an incorrect extension (such as tar when one meant
tgz) will be treated properly.

TODO:
- Handle symlinks in ISO9660 files
- Handle non-regular files in tar (hard, sym, character, block, etc...)

Improvements:
- Rework EOF handling of individual archives, size vs read error

## Example Commandline
```bash
./archive-exploder -output out -input debian-11.3.0-amd64-netinst.iso -r 3
```

## Usage
```
$ ./archive-exploder
Archive Exploder,  Version: 0.1.20220501.2231

Usage: ./archive-exploder [options...]

  -debug
        Turn on debug, more verbose
  -input string
        Path to put the extracted files
  -output string
        Path to put the extracted files (default ".")
  -r int
        Levels of recusion (archives-inside-archives) to expand (default 1)

Formats supported:
  - gzip
  - iso9660
  - rpm
  - tar
  - zip
```

## Example single level recusion:

```
$ ./archive-exploder -input tests/createrepo_c-0.10.0-18.el7.x86_64.rpm -output out
$ find out/
out/
out/usr
out/usr/bin
out/usr/bin/createrepo_c
out/usr/bin/mergerepo_c
out/usr/bin/modifyrepo_c
out/usr/bin/sqliterepo_c
out/usr/share
out/usr/share/bash-completion
out/usr/share/bash-completion/completions
out/usr/share/bash-completion/completions/createrepo_c
out/usr/share/doc
out/usr/share/doc/createrepo_c-0.10.0
out/usr/share/doc/createrepo_c-0.10.0/README.md
out/usr/share/man
out/usr/share/man/man8
out/usr/share/man/man8/createrepo_c.8.gz
out/usr/share/man/man8/mergerepo_c.8.gz
out/usr/share/man/man8/modifyrepo_c.8.gz
out/usr/share/man/man8/sqliterepo_c.8.gz
```
