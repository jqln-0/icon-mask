# icon-mask

A simple implementation of Android-style icon masking for GTK (and possibly
KDE) icon themes.

## Installation

The program is written in Go, so installation should be simple:

```bash
$ go get github.com/jqln-0/icon-mask
```

## Usage

You will need

 - An icon theme
 - A base image to 'mask' with.

Running the program with the `--help` flag should guide you the rest of the
way. An example execution for the 'Moka' theme is:

```bash
$ icon-mask --extra-themes=Faba Moka base.png out
```

You can then add the generated icons to the chosen theme by simply merging the
output and theme folders;

```bash
$ cp -r out/* /usr/share/icons/Moka/
```

## TODO

 - Make SVG output work,
 - Perform compositing concurrently,
 - Auto-generate symlinks?

