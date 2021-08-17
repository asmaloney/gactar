![Build](https://github.com/asmaloney/gactar/actions/workflows/build.yaml/badge.svg)

# gactar

`gactar` is a tool for creating [ACT-R](https://en.wikipedia.org/wiki/ACT-R) models using a new declarative file format (called _amod_).

## Proof-of-Concept

**This is just a proof-of-concept.**

Currently, `gactar` will take an [_amod_ file](#amod-file-format) and generate code to run it on three different ACT-R implementations:

- [CCM Suite](https://github.com/CarletonCognitiveModelingLab/CCMSuite3) (python)
- [pyactr](https://github.com/jakdot/pyactr) (python)
- ["vanilla" ACT-R](https://act-r.psy.cmu.edu/) (lisp)

`gactar` will work with the short tutorial models included in the _examples_ directory. It doesn't handle a lot beyond what's in there - it only works with memory modules, not perceptual-motor ones - so _it's limited at the moment_.

The format still feels a little heavy, so if I continue with this project I would expect to iterate on it.

## Why?

1. Provides a human-readable, easy-to-understand, standard format to define basic ACT-R models.
2. Allows the easy exchange of models with other researchers.
3. Abstracts away the "programming" to focus on writing and understanding models.
4. Restricts the model to a small language to prevent programming "outside the model".
5. Provides a very simple setup for teaching environments.
6. Runs the same model on multiple ACT-R implementations.
7. Generates human-readable code (for now!) which is useful for learning the implementations and comparing them.
8. Parses chunks to catch and report errors in a user-friendly manner.

   **Example #1 (invalid variable name)**

   ```
    match {
        goal `isMember( ?obj ?cat None )`
    }
    do {
        recall `property( ?ojb category ? )`
    }
   ```

   The CCM Suite implementation _fails silently_ when given invalid variables which makes it difficult to catch errors & can result in incorrect output. Instead of ignoring the incorrect variable, gactar outputs a nice error message so it's obvious what the problem is:

   ```
   recall statement variable '?ojb' not found in matches for production 'initialRetrieval' (line 58)
   ```

   **Example #2 (invalid slot name)**

   ```
    match {
        goal `isMember( ?obj ?cat None )`
    }
    do {
        set resutl of goal to 'pending'
    }
   ```

   The CCM Suite implementation produces the following error:

   ```
   Traceback (most recent call last):
   File "/path/gactar_Semantic_Run.py", line 8, in <module>
    model.run()
   File "/path/CCMSuite3/ccm/model.py", line 254, in run
    self.sch.run()
   File "/path/CCMSuite3/ccm/scheduler.py", line 116, in run
    self.do_event(heapq.heappop(self.queue))
   File "/path/CCMSuite3/ccm/scheduler.py", line 161, in do_event
    result=event.func(*event.args,**event.keys)
   File "/path/CCMSuite3/ccm/lib/actr/core.py", line 64, in _process_productions
    choice.fire(self._context)
   File "/path/CCMSuite3/ccm/production.py", line 51, in fire
    exec(self.func, context, self.bound)
   File "<production-initialRetrieval>", line 2, in <module>
   File "/path/CCMSuite3/ccm/model.py", line 22, in __call__
    val = self.func(self.obj, *args, **keys)
   File "/path/CCMSuite3/ccm/lib/actr/buffer.py", line 60, in modify
    raise Exception('No slot "%s" to modify to "%s"' % (k, v))
   Exception: No slot "resutl" to modify to "pending"
   end...
   ```

   Instead, by adding validation, gactar produces a much better message:

   ```
   slot 'resutl' does not exist in match buffer 'goal' in production 'initialRetrieval' (line 57)
   ```

## Setup

1. Although the `gactar` executable itself is compiled for each platform, it requires **python3** to run the setup and to run the _ccm_ and _pyactr_ implementations. **python3** needs to be somewhere in your `PATH` environment variable.

2. `gactar` requires one or more of the three implementations (_ccm_, _pyactr_, _vanilla_) be installed.

`gactar` uses a python virtual environment to keep all the required python packages, lisp files, and other implementation files in one place so it does not affect the rest of your system. For more information about the virtual environment see the [python docs](https://docs.python.org/3/library/venv.html).

### Setup Virtual Environment

1. Run `./scripts/setup.sh`

   This will create a virtual environment for the project in a directory called `env`, download the [CCM Suite](https://github.com/CarletonCognitiveModelingLab/CCMSuite3) & put its files in the right place, and install [pyactr](https://github.com/jakdot/pyactr) using pip.

2. You will need to activate the virtual environment by running this in the terminal before you run `gactar`:

   ```sh
   source ./env/bin/activate
   ```

   If it activated properly, your command line prompt will start with `(env)`. If you want to deactivate it, run `deactivate`.

### Install SBCL Lisp Compiler

For now this is not automated because the required files are not easy to determine programmatically. I may be able to improve this in the future by adding it to the auto-setup process.

1. We are using the [Steel Bank Common Lisp](https://www.sbcl.org/index.html) (sbcl) compiler. Download the correct version [from here](https://www.sbcl.org/platform-table.html) by finding your platform (OS and architecture) in the table and clicking its box. Put the file in the `env` directory and unpack it there.

2. To install it in our environment, change to the new directory it created (e.g. `sbcl-1.2.11-x86-64-darwin`) and run this command (setting the path to wherever the env directory is):
   ```sh
   INSTALL_ROOT=/path/to/gactar/env/ ./install.sh
   ```

### Install Vanilla ACT-R

For now this is not automated because the required files are not easy to determine programmatically. I may be able to improve this in the future by adding it to the auto-setup process.

1. Download the zip file for your OS from [here](https://act-r.psy.cmu.edu/software). Put the zip file in the `env` directory and unpack it there. This should create a directory named `actr7.x`

2. Back in the `env` directory, run the following command to compile the main actr files using the lisp compiler (setting the path to wherever the env directory is):
   ```sh
   export SBCL_HOME=/path/to/env/lib/sbcl; sbcl --script actr7.x/load-single-threaded-act-r.lisp
   ```
   This will take a few moments to compile all the ACT-R files so it is ready to use.

## Build

If you want to build `gactar`, you will need [git](https://git-scm.com/) and the [go compiler](https://golang.org/) installed.

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

## Test

To run the built-in tests, from the top-level of the repo run:

```
go test ./...
```

## Usage

```
gactar [OPTIONS] [FILES...]
```

### Command Line Options

**--debug, -d**: turn on debugging output (mainly output tokens from lexer)

**--ebnf**: output amod EBNF to stdout and quit

**--interactive, -i**: run an interactive shell

**--port, -p** [number]: port to run the webserver on (default: 8181)

**--web, -w**: start a webserver to run in a browser

## Example Usage

These examples assume you have set up your virtual environment properly. See [setup](#setup) above.

### Generate a Python File

```
(env)$ ./gactar examples/count.amod
gactar version v0.0.2
ccm: Using Python 3.9.6 from /path/to/gactar/env/bin/python3
-- Generating code for examples/count.amod
   Written to gactar_ccm_Count.py.py
```

This will generate a python file called `gactar_ccm_Count.py.py` in the directory you are running from. It doesn't contain the run command, so in order to use it you would need to create another python file like this:

```py
from gactar_Count import gactar_Count


model = gactar_Count()
model.goal.set('countFrom 2 5 starting')
model.run()
```

Currently this form only generates the `ccm` version. This will be [fixed in the future](https://github.com/asmaloney/gactar/issues/15).

### Run Interactively

```
(env)$ ./gactar -i
gactar version v0.0.2
Type 'help' for a list of commands.
To exit, type 'exit' or 'quit'.
ccm: Using Python 3.9.6 from /path/to/gactar/env/bin/python3
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

Currently only runs the `ccm` version. This will be [fixed in the future](https://github.com/asmaloney/gactar/issues/15).

### Run As Web Server

```
(env)$ ./gactar -w
ccm: Using Python 3.9.6 from /path/to/gactar/env/bin/python3
pyactr: Using Python 3.9.6 from /path/to/gactar/env/bin/python3
vanilla: Using SBCL 1.2.11 from /path/to/gactar/env/bin/sbcl
Serving gactar on http://localhost:8181
```

Open `http://localhost:8181` in your browser, select an example from the menu, modify the amod description &amp; the goal, and click **Run**.

## amod File Format

Here is an example of the file format:

```
==model==

// The name of the model (used when generating code and for error messages)
name: count

// Description of the model (currently output as a comment in the generated code)
description: 'This is a model which adds numbers. Based on the u1_count.py tutorial.'

// Examples of starting goals to use when running the model
examples {
    'countFrom 2 5 starting'
    'countFrom 1 3 starting'
}

==config==

// Turn on logging by setting 'log' to 'true' or 1
actr { log: true }

// Declare chunks and their layouts
chunks {
    count( first second )
    countFrom( start end status )
}

==init==

// Initialize the memory
memory {
    `count( 0 1 )`
    `count( 1 2 )`
    `count( 2 3 )`
    `count( 3 4 )`
    `count( 4 5 )`
}

==productions==

// Name of the production
start {
    // Buffers to match
    match {
        goal `countFrom( ?start ?end starting )`
    }
    // Steps to execute
    do {
        recall `count( ?start ?)`
        set goal to `countFrom( ?start ?end counting )`
    }
}

increment {
    match {
        goal `countFrom( ?x !?x counting )`
        retrieval `count( ?x ?next )`
    }
    do {
        print x
        recall `count( ?next ? )`
        set start of goal to next
    }
}

stop {
    match {
        goal `countFrom( ?x ?x counting )`
    }
    do {
        print x
        clear goal
    }
}
```

You can find other examples of `amod` files in the [examples folder](examples).

### Special Chunks

User-defined chunks must not begin with '\_' or be named `goal`, `retrieval`, or `memory` - these are reserved for internal use. Currently there is one internal chunk - _\_status_ - which is used to check the status of buffers and memory.

It is used in a `match` as follows:

```
match {
    goal `_status( full )`
    memory `_status( error )`
}
```

For buffers, the valid statuses are `full` and `empty`.

For memory, valid statuses are `busy`, `free`, `error`.

### Pattern Syntax

The _match_ section matches _patterns_ to buffers. Patterns are delineated by backticks - e.g. `` `property( ?obj category ?cat )` ``. The first item is the chunk name and the items between the parentheses are the slots. These are parsed to ensure their format is consistent with _chunks_ which are declared in the _config_ section.

The _do_ section in the productions uses a small language which currently understands the following commands:

| command                                                               | example                         |
| --------------------------------------------------------------------- | ------------------------------- |
| clear _(buffer name)+_                                                | clear goal, retrieval           |
| print _(string or ident or number)+_                                  | print foo, 'text', 42           |
| recall _(pattern)_                                                    | recall \`car( ?colour )\`       |
| set _(slot name)_ of _(buffer name)_ to _(string or ident or number)_ | set sum of goal to 6            |
| set _(buffer name)_ to _(pattern)_                                    | set goal to \`start( 6 None )\` |
| write _(string or ident or number)+_ to _(text output name)_          | write 'foo' to text             |
