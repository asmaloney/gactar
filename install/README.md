# Install Files

This directory contains the requirements files for installing python packages using `pip`.

Installation is handled using gactar itself:

```
./gactar env setup
```

## Developer Packages

If you are going to use the gactar virtual environment in an IDE to edit python files, it might be useful to install the dev requirements. This will install a linter (pylint) and a formatter (autopep8).

From the main gactar directory you need to activate the virtual environment then run `pip` to install the packages:

```
$ . ./env/bin/activate
$ pip install -r install/requirements-dev.txt
```

If you want to deactivate the virtual environment:

```
$ deactivate
```
