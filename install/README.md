# Install Files

This directory contains the requirements files for installing python packages using `pip`.

Installation is handled using gactar itself:

```
./gactar env setup
```

## Developer Packages

If you are going to use the gactar virtual environment in an IDE to edit python files, it might be useful to install the dev requirements. This will install a linter (pylint) and a formatter (autopep8).

To include the developer packages:

```
./gactar env setup -dev
```
