# gactar

`gactar` is a tool for creating [ACT-R](https://en.wikipedia.org/wiki/ACT-R) models using a new declarative file format (called _amod_).

## Proof-of-Concept

**This is just a proof-of-concept.**

Currently, `gactar` will take an [_amod_ file](#amod-file-format) and generate the python code to run it with the [CCM Suite](https://github.com/CarletonCognitiveModelingLab/CCMSuite3).

gactar will work with the small tutorial models included in the _examples_ directory. It doesn't handle a lot beyond what's in there.

The format still feels a little heavy, so if I continue with this project I would expect to iterate on it. One goal would be to remove python from the "do blocks" by defining a parsable language to manipulate the model. This would have the advantage of allow other "backends" besides CCMSuite and would also formalize the writing of ACT-R models by defining a proper language to do so.

## Requirements

gactar requires two things:

1.  **python3** needs to be somewhere in your environment's `PATH`
2.  The [CCM Suite](https://github.com/CarletonCognitiveModelingLab/CCMSuite3) (for python3) needs to be available in `PYTHONPATH`.

    On Linux/macOS, you can do this in the terminal before running `gactar`:

    ```
    export PYTHONPATH=/path/to/CCMSuite3/
    ```

## Build

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
    "countFrom 1 3 starting"
    "countFrom 2 5 starting"
}

==config==

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
    { "count 0 1" }
    { "count 1 2" }
    { "count 2 3" }
    { "count 3 4" }
    { "count 4 5" }
}

==productions==

// Name of the production
start {
    // Buffers to match
    match {
        goal: 'countFrom ?start ?end starting'
    }
    // Steps to execute (currently python code)
    do #<
        memory.request('count ?start ?next')
        goal.set('countFrom ?start ?end counting')
    >#
}

increment {
    match {
        goal: 'countFrom ?x !?x counting'
        retrieve: 'count ?x ?next'
    }
    do #<
        print(x)
        memory.request('count ?next ?nextNext')
        goal.modify(_1=next)
    >#
}

stop {
    match {
        goal: 'countFrom ?x ?x counting'
    }
    do #<
        print(x)
        goal.set('countFrom ?x ?x stop')
    >#
}
```

## Examples

### Generate a python file

```
gactar examples/count.amod
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
$ ./gactar -i
gactar version v0.0.1
Type 'help' for a list of commands.
To exit, type 'exit' or 'quit'.
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
gactar -w
Serving gactar on http://localhost:8181
```

Open `http://localhost:8181` in your browser, modify the amod description &amp; the goal, and click **Run**.
