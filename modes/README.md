# Modes

gactar can be used in several different modes. This directory contains the code to handle each of them.

More details may be found in the main [README file](../README.md).

## Default Mode

If you use the command line without the `web` or `cli` command, gactar lets the user process a file (and optionally run it with `-r`).

```sh
$ ./gactar {amod file}
```

Run `./gactar help` for a list of options.

## CLI (Interactive)

This allow the user to use commands to work with gactar interactively - loading & running amod files.

```sh
$ ./gactar cli
```

Run `./gactar help cli` for a list of options.

## Web

This runs a web server with an HTTP interface to load & run amod models. Details of the endpoints may be found in the [doc directory](../doc/Web%20API.md).

```sh
$ ./gactar web
```

Run `./gactar help web` for a list of options.
