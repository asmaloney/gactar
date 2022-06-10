# Modes

gactar can be used in several different modes. This directory contains the code to handle each of them.

More details may be found in the main [README file](../README.md).

## Default Mode (CLI)

If you use the command line without the `-i` or `-w` options, gactar lets the user process a file (and optionally run it with `-r`).

```sh
$ ./gactar {amod file}
```

## Shell (Interactive CLI)

This allow the user to use commands to work with gactar interactively - loading & running amod files.

```sh
$ ./gactar -i
```

## Web

This runs a web server with an HTTP interface to load & run amod models. Details of the endpoints may be found in the [doc directory](../doc/Web%20API.md).

```sh
$ ./gactar -w
```
