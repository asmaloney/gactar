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
6. Parses buffer patterns to catch and report errors.

   **Example #1 (invalid variable name)**

   ```
    match {
        goal: `isMember ?obj ?cat None`
    }
    do {
        recall `property ?ojb category ?`
    }
   ```

   The CCM Suite implementation _fails silently_ when given invalid variables which makes it difficult to catch errors & can result in incorrect output. Instead of ignoring the incorrect variable, gactar outputs a nice error message so it's obvious what the problem is:

   ```
   recall statement variable '?ojb' not found in matches for production 'initialRetrieve' (line 58)
   ```

   **Example #2 (invalid slot name)**

   ```
    match {
        goal: `isMember ?obj ?cat None`
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
   File "<production-initialRetrieve>", line 2, in <module>
   File "/path/CCMSuite3/ccm/model.py", line 22, in __call__
    val = self.func(self.obj, *args, **keys)
   File "/path/CCMSuite3/ccm/lib/actr/buffer.py", line 60, in modify
    raise Exception('No slot "%s" to modify to "%s"' % (k, v))
   Exception: No slot "resutl" to modify to "pending"
   end...
   ```

   Instead, by adding validation, gactar produces a much better message:

   ```
   slot 'resutl' does not exist in match buffer 'goal' in production 'initialRetrieve' (line 57)
   ```

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
gactar [OPTIONS] [FILES...]
```

### OPTIONS

**--debug, -d**: turn on debugging output (mainly output tokens from lexer)

**--ebnf**: output amod EBNF to stdout and quit

**--interactive, -i**: run an interactive shell

**--port, -p** [number]: port to run the webserver on (default: 8181)

**--web, -w**: start a webserver to run in a browser

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
actr { log: false }

// Declare chunks and their layouts
chunks {
    count( first second )
    countFrom( start end status )
}

==init==

// Initialize the memory
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
        goal: `countFrom ?start ?end starting`
    }
    // Steps to execute
    do {
        recall `count ?start ?`
        set goal to `countFrom ?start ?end counting`
    }
}

increment {
    match {
        goal: `countFrom ?x !?x counting`
        retrieve: `count ?x ?next`
    }
    do {
        print x
        recall `count ?next ?`
        set start of goal to next
    }
}

stop {
    match {
        goal: `countFrom ?x ?x counting`
    }
    do {
        print x
        clear goal
    }
}
```

You can find other examples of amod files in the [examples folder](examples).

### Special Chunks

User-defined chunks must not begin with '\_' - these are reserved for internal use. Currently there is one internal chunk - _\_status_ - which is used to check the status of buffers and memory.

It is used in a `match` as follows:

```
match {
    goal: `_status full`
    memory: `_status error`
}
```

For buffers, the valid statuses are `full` and `empty`.

For memory, valid statuses are `busy`, `free`, `error`.

### Syntax

The _match_ section matches _patterns_ to buffers. Patterns are delineated by backticks - e.g. `` `property ?obj category ?cat` ``. The first item is the chunk name and the others are the slots. These are parsed to ensure their format is consistent with chunks which are declared in the _config_ section.

The _do_ section in the productions uses a small language which currently understands the following commands:

| command                                                         | example                      |
| --------------------------------------------------------------- | ---------------------------- |
| clear _(buffer name)+_                                          | clear buff1, buff2           |
| print _(string or ident or number)+_                            | print foo, 'text', 42        |
| recall _(pattern)_                                              | recall \`car ?colour\`       |
| set _name_ of _(buffer name)_ to _(string or ident or number)_  | set sum of goal to 6         |
| set _(buffer name)_ to _(string or ident or number or pattern)_ | set goal to \`start 6 None\` |
| write _(string or ident or number)+_ to _(text output name)_    | write 'foo' to text          |
