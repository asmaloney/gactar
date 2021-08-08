![Build](https://github.com/asmaloney/gactar/actions/workflows/build.yaml/badge.svg)

# gactar

`gactar` is a tool for creating [ACT-R](https://en.wikipedia.org/wiki/ACT-R) models using a new declarative file format (called _amod_).

## Proof-of-Concept

**This is just a proof-of-concept.**

Currently, `gactar` will take an [_amod_ file](#amod-file-format) and generate the python code to run it with the [CCM Suite](https://github.com/CarletonCognitiveModelingLab/CCMSuite3).

gactar will work with the small tutorial models included in the _examples_ directory. It doesn't handle a lot beyond what's in there - it only works with memory modules, not perceptual-motor ones - so _it's limited at the moment_.

The format still feels a little heavy, so if I continue with this project I would expect to iterate on it.

## Why?

1. Provides a human-readable, easy-to-understand, standard format to define basic ACT-R models.
2. Allows the easy exchange of models with other researchers.
3. Abstracts away the "programming" to focus on writing and understanding models.
4. Restricts the model to a small language to prevent programming "outside the model".
5. Provides a very simple setup for teaching environments.

## Setup

1. `gactar` requires **python3** which needs to be somewhere in your environment's `PATH` environment variable.

2. `gactar` requires the [CCM Suite](https://github.com/CarletonCognitiveModelingLab/CCMSuite3) (for python3) - see the following two options for how to set that up.

### Setup with virtual python environment (easiest)

A python virtual environment keeps all of your python packages local to your project so it does not affect the rest of your system. For more information see the [python docs](https://docs.python.org/3/library/venv.html).

1. Run `./scripts/setupPython.sh`

   This will create a virtual environment for the project, download the [CCM Suite](https://github.com/CarletonCognitiveModelingLab/CCMSuite3), and put its files in the right place.

2. You need to activate the virtual environment by running this in the terminal before you run gactar:

   ```sh
   source ./pyenv/bin/activate
   ```

   If it activated properly, your command line prompt will start with `(pyenv)`. To deactivate it, run `deactivate`.

### Setup by cloning CCMSuite

2.  Clone the [CCM Suite](https://github.com/CarletonCognitiveModelingLab/CCMSuite3) (for python3):

    ```sh
    git clone https://github.com/CarletonCognitiveModelingLab/CCMSuite3
    ```

3.  The ccm package from there needs to be available in your `PYTHONPATH`.

    You can do this in the terminal each time you want to run `gactar` (or you can set it in your environment variables):

    ```
    export PYTHONPATH=/path/to/CCMSuite3/
    ```

    Note that setting PYTHONPATH affects your entire system, so it may interfere with other python projects.

## Build

If you want to build `gactar`, you will need the [go compiler](https://golang.org/) installed.

Then you just need to clone this repo:

```sh
git clone https://github.com/asmaloney/gactar
cd gactar
```

...and run the build command:

```
go build
```

This will create the `gactar` executable.

## Usage

```
gactar [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

### GLOBAL OPTIONS

**--debug, -d**: turn on debugging output (mainly output tokens from lexer)

**--ebnf**: output amod EBNF to stdout and quit

**--interactive, -i**: run an interactive shell

**--port, -p** [number]: port to run the webserver on (default: 8181)

**--web, -w**: start a webserver to run in a browser

## amod File Format

Here is an example of the file format:

```
==model==

// The name of the model (used when generating code and for error messages)
name: count

// Description of the model (currently output as a comment in the python code)
description: "This is a model which adds numbers. Based on the u1_count.py tutorial."

// Examples of starting goals to use when running the model
examples {
    "countFrom 2 5 starting"
    "countFrom 1 7 starting"
}

==config==

// Turn on logging by setting 'log' to 'true' or 1
actr { log: false }

// List of buffers to create by name
buffers { goal, retrieve }

// Memories to create
memories {
    memory {
        // Attach this buffer by name
        buffer: retrieve
    }
}

==init==

// Initialize the memory named "memory"
memory {
    'count 0 1'
    'count 1 2'
    'count 2 3'
    'count 3 4'
    'count 4 5'
}

==productions==

// Name of the production
start {
    // Buffers to match
    match {
        goal: 'countFrom ?start ?end starting'
    }
    // Steps to execute
    do {
        recall 'count ?start ?next' from memory
        set goal to 'countFrom ?start ?end counting'
    }
}

increment {
    match {
        goal: 'countFrom ?x !?x counting'
        retrieve: 'count ?x ?next'
    }
    do {
        print x
        recall 'count ?next ?nextNext' from memory
        set field 1 of goal to next
    }
}

stop {
    match {
        goal: 'countFrom ?x ?x counting'
    }
    do {
        print x
        set goal to 'countFrom ?x ?x stop'
    }
}
```

You can find other examples of amod files in the [examples folder](examples).

### Syntax

The _do_ section in the productions uses a small language which currently understands the following commands:

| command                                                                          | example                          |
| -------------------------------------------------------------------------------- | -------------------------------- |
| print _(!keyword)+_                                                              | print foo, 'text', 42            |
| recall _(string)_ from _(memory name)_                                           | recall 'car ?colour' from memory |
| set field _(number or name)_ of _(buffer name)_ to _(string or ident or number)_ | set field 1 of goal to 6         |
| set _(buffer name)_ to _(string or ident or number)_                             | set goal to 'start 6 None'       |
| write _(!keyword)+_ to _(text output name)_                                      | write 'foo' to text              |

## Example Usage

These examples assume you have set up your environment properly - either using python's virtual environment or by setting up your PYTHONPATH. See [setup](#setup) above.

### Generate a python file

```
(pyenv)$ ./gactar examples/count.amod
gactar version v0.0.2
Using Python 3.9.6 from /path/to/gactar/pyenv/bin/python3
-- Generating code for examples/count.amod
   Written to gactar_Count.py
```

This will generate a python file called `gactar_Count.py` in the directory you are running from. It doesn't contain the run command, so in order to use it you would need to create another python file like this:

```py
from gactar_Count import gactar_Count


model = gactar_Count()
model.goal.set('countFrom 2 5 starting')
model.run()
```

### Run interactively

```
(pyenv)$ ./gactar -i
gactar version v0.0.2
Type 'help' for a list of commands.
To exit, type 'exit' or 'quit'.
Using Python 3.9.6 from /path/to/gactar/pyenv/bin/python3
> help
  exit:     exits the program
  history:  outputs your command history
  load:     loads a model: load [FILENAME]
  quit:     exits the program
  reset:    resets the current model
  run:      runs the current model: run [INITIAL STATE]
  version:  outputs version info
> load examples/count.amod
 model loaded
 examples:
           run countFrom 2 5 starting
           run countFrom 1 7 starting
> run countFrom 2 5 starting
2
3
4
5
end...
> quit
```

### Run as a web server

```
(pyenv)$ ./gactar -w
Using Python 3.9.6 from /path/to/gactar/pyenv/bin/python3
Serving gactar on http://localhost:8181
```

Open `http://localhost:8181` in your browser, modify the amod description &amp; the goal, and click **Run**.
