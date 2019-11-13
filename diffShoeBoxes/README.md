# diff shoe boxes

## setup
This tool require libmagickwand. PLEASE READ [LIBRARY README](https://github.com/gographics/imagick/blob/master/README.md). This section is only shortly summary. 

### Debian / Ubuntu
1. install libmagickwand `sudo apt install libmagickwand-dev`
2. check those lib version `pkg-config --cflags --libs MagickWand`
3. read [GoImagick README](https://github.com/gographics/imagick/blob/master/README.md) and check lib version range(note: now i use v3lib, this is compatible **7.x <= ImageMagick <= 7.x** )
4. if no issue, `go get -u`. you perceive issue, please fix library version.

### ArchLinux
Rough flow is same as Debian, Ubuntu(ex: apt -> pacman). But ArchLinuxRepository is Rolling release. Maybe no compatible library in ArchLinuxRepository. I recommend no to use Arch, and use DebianFamily. Otherwise to need [build ImageMagick](https://imagemagick.org/script/install-source.php#unix).

## Troubleshooting
### Close library repository
Please rewrite diff function. replace exec. 