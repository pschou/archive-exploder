# Archive Exploder

The goal of this tool is to have the ability to expand archives within archives.  This extraction, at every level, is done by use of a stream reader.  This means the final files, no matter how many layer deep of archives, can be all expanded at once.  Even more important is the reader of the archive goes in the forward only direction and the writer does not need to read in any part that is written.  This makes the best use of every drive input-output operation (IOPS).  Speedy, fast, and accurate.

This tool takes an archive file and will extract/expand it.  If an archive is contained within the current archive, it will also be expanded into a folder named with the name of the holding archive.  This archive-inside-archive ability is called recursion. By default, the recursion limit is set to 1, but increasing this enables deeper levels of extraction.

Take for example, a compressed (zip) iso file (iso9660) that contains Debian package files.  Once this archive-exploder is called, one will be left with directories and files with all the text files, media, and binaries which are contained within the original iso file.

Space is a major consideration driving the need for this tool.  If one were extracting a large DVD image (4GB) filled with tar files, and inside each they have war/jar files.  In this case each layer would add N x Size, so, with this particular case, 4 levels of recursion, starting with a 4GB DVD, one could easily see 10GB in a final write out directory and upwards 40 GB of disk usage in holding the intermediary files.

*This tool gives you the end file structure without the transient bloat.*

A thing to note is the file extension is not used to determine file type, but instead the underlying bytes.  This way, a file with an incorrect extension (such as tar when one meant tgz) will be adequately treated.


TODO:
- Handle symlinks in ISO9660 files
- Handle non-regular files in tar (hard, sym, character, block, etc...)
- Add timestamps for internal files
- When an archive matches multiple, try them one at a time to see which works

Improvements:
- Rework EOF handling of individual archives, size vs read error

## Example ISO
```bash
./archive-exploder -output out -input debian-11.3.0-amd64-netinst.iso -r 3
```

## Example APK and DEB
One must note that for deb and apk files, the package file is a compilation of two or more archives smashed together, so you need to use recursion twice to extract the header files from the header archive(s) and then content archive.
```bash
./archive-exploder --input tests/aaudit-0.7.2-r2.apk -output apk/ -r 2
```

```bash
./archive-exploder -input tests/vagrant_2.2.6_i686.deb -output testva/ -r 2
```

## Usage
```
$ ./archive-exploder -h
Archive Exploder,  Version: 0.1.20220517.1523

Usage: ./archive-exploder [options...]

  -debug
        Turn on debug, more verbose
  -input string
        Path to input archive file
  -output string
        Path to put the extracted file(s) (default ".")
  -r int
        Levels of recusion (archives-inside-archives) to expand (default 1)

Formats supported:
  - 7zip
  - bzip2
  - cab
  - debian
  - gzip / bgzf / apk
  - iso9660
  - lzma
  - rar
  - rpm
  - tar
  - xz
  - zip
  - zstd
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
