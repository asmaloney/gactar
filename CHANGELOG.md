# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## 0.12.0 - (in progress)

### Added

- New command-line command `module` outputs information about modules. Currently includes two subcommands:
  - `info [name]` outputs detailed info about a module - name, version, description, any buffers, and any parameters. `name` can be a space-separated list of modules or `all`.
  - `list` outputs the list of modules - name, version, and description

### Changed

- ACT-R (vanilla) was updated to version 7.27.7.

### Fixed

- Improved some error messages.

## [0.11.0](https://github.com/asmaloney/gactar/releases/tag/v0.11.0) - 2023-01-03

### Added

- {shell} Interactive mode now recognizes up a & down arrow keys to navigate history. ([#287](https://github.com/asmaloney/gactar/pull/287))

- Command line output now uses colour. ([#284](https://github.com/asmaloney/gactar/pull/284))
  - May be turned off using a command-line option (`--no-colour` or `--no-color`) or by setting the `NO_COLOR` environment variable.

### Changed

- Command line options changed to commands ([#298](https://github.com/asmaloney/gactar/pull/298))
  - Instead of `gactar -w`, it is now called using `gactar web`.
  - Instead of `gactar -i`, it is now called using `gactar cli`.
  - Run `gactar help` for a list of commands and options.

### Fixed

- {cli} Fixes the `version` command. ([#286](https://github.com/asmaloney/gactar/pull/286))

## [0.10.0](https://github.com/asmaloney/gactar/releases/tag/v0.10.0) - 2022-07-07

### Added

- Now tracks and outputs declarations for implicit chunk names. This avoids warnings on some frameworks. ([#241](https://github.com/asmaloney/gactar/pull/241), [#247](https://github.com/asmaloney/gactar/pull/247), [#249](https://github.com/asmaloney/gactar/pull/249))

- Allow strings in chunk patterns. ([#243](https://github.com/asmaloney/gactar/pull/243))

- Allow naming of initialized chunks. ([#250](https://github.com/asmaloney/gactar/pull/250))

  e.g.

  ```
  ~~ init ~~

  memory {
    castle  [meaning: 'castle']
    earl    [meaning: 'earl']
  }
  ```

  In pyactr and vanilla, these names are used in the chunk creation. In ccm, the names are added as comments as it doesn't seem to use the "chunk name" concept.

- Allow setting of similarities in the _init_ section. ([#257](https://github.com/asmaloney/gactar/pull/257))

  They are specified like this:

  ```
  ~~ init ~~

  similar {
    ( first second -0.5 )
    ( second third -0.5 )
  }
  ```

- Added _random_seed_ option to the **gactar** section. This sets the seed to use for generating pseudo-random numbers (allows for reproducible runs). ([#265](https://github.com/asmaloney/gactar/pull/265))

- Added tabs to the web UI output section to split out each framework's results. ([#269](https://github.com/asmaloney/gactar/pull/269))

### Changed

- Replaced _partial_matching_ option from the **procedural** module with the _mismatch_penalty_ option for the **memory** module. Setting this turns on partial matching and sets the penalty in the activation equation to this value. ([#261](https://github.com/asmaloney/gactar/pull/261))

- Updated web UI to vue 2.7.x. ([#272](https://github.com/asmaloney/gactar/pull/272))

### Fixed

- Give proper error when trying to use an invalid type with `_status`. ([#242](https://github.com/asmaloney/gactar/pull/242))

- Check for reserved chunk names. ([#271](https://github.com/asmaloney/gactar/pull/271))

  The following names are reserved according to ACT-R: `busy`, `clear`, `empty`, `error`, `failure`, `free`, `full`, `requested`, and `unrequested`

## [0.9.0](https://github.com/asmaloney/gactar/releases/tag/v0.9.0) - 2022-06-20

### Added

- gactar now handles installation of python packages, ACT-R code, and the Lisp compiler itself instead of using external scripts. ([#212](https://github.com/asmaloney/gactar/pull/212))

  There is a new command to run setup:

  ```
  $ ./gactar env setup
  ```

  Use the `-dev` flag to also install optional developer packages for linting & formatting Python code.

  ```
  $ ./gactar env setup -dev
  ```

- gactar's new setup capability should work on **Windows** with a couple of caveats:

  - It has only been tried with the 3.10.5 release from [python.org](https://www.python.org/downloads/windows/) on Windows 10.
  - gactar uses the PATH environment variable to find the Python interpreter. The easiest way to do this is to check the **_Add Python 3.10 to PATH_** checkbox when installing Python.
  - The Clozure Common Lisp compiler is currently broken on Windows (waiting on a new build). It will download, but will fail to run.

- Added a command to check the health of your virtual environment. ([#220](https://github.com/asmaloney/gactar/pull/220))

  ```
  $ ./gactar env doctor
  ```

  This will check paths, ensure that Python packages are installed properly, and check for the lisp compiler.

- Added an `extra_buffers` module to allow declaration of... extra buffers. ([#217](https://github.com/asmaloney/gactar/pull/217))

  Declare them in the module config section like this (they currently don't have any configuration options):

  ```
  modules {
      extra_buffers {
          foo {}
          bar {}
      }
  }
  ```

- Added _partial_matching_ option to the **procedural** module to turn on partial matching. ([#223](https://github.com/asmaloney/gactar/pull/223))

  **Note:** while this can be turned on, specifying similarity of chunks isn't handled yet. (See [#234](https://github.com/asmaloney/gactar/issues/234))

- Added _decay_ option to the declarative **memory** module for the base-level learning calculation. ([#226](https://github.com/asmaloney/gactar/pull/226))

### Changed

- Allow ID in `set` statements. ([#200](https://github.com/asmaloney/gactar/pull/200))

  Instead of:

  ```
  set goal.state to 'harvest_location'
  ```

  You can use it without quotes:

  ```
  set goal.state to harvest_location
  ```

- Web assets are now compressed using brotli compression. ([#218](https://github.com/asmaloney/gactar/pull/218))

- Moved the default temp directory (`gactar-temp`) into the environment directory. ([#229](https://github.com/asmaloney/gactar/pull/229))

### Fixed

- When not running `setup`, restrict the PATH environment variable to paths within the virtual environment directory. ([#230](https://github.com/asmaloney/gactar/pull/230))

## [0.8.0](https://github.com/asmaloney/gactar/releases/tag/v0.8.0) - 2022-06-13

### Added

- Added a new statement to the amod language: **stop**. ([#170](https://github.com/asmaloney/gactar/pull/170))
- Web UI now highlights any errors in the code editor. ([#197](https://github.com/asmaloney/gactar/pull/197))

### Changed

- Removed developer packages (autopep8 & pylint) from general installation of pip packages. These may be installed by running these commands in the gactar directory:
  ```sh
  $ . ./env/bin/activate
  (env) $ pip install -r ./scripts/requirements-dev.txt
  ```
- Reduce binary size by turning off some cli documentation tools.
- Replace [Steel Bank Common Lisp compiler](http://www.sbcl.org) (sbcl) with the [Clozure Common Lisp compiler](https://ccl.clozure.com/) (ccl). ([#191](https://github.com/asmaloney/gactar/pull/191))
- Grammar changes:
  - Replace `==` in sections headers with `~~` and allow spaces. ([#192](https://github.com/asmaloney/gactar/pull/192))
  - Add a `when` clause to replace the complicated internal format. ([#193](https://github.com/asmaloney/gactar/pull/193))
    ```
    match {
        goal [add: * ?num2 ?count!?num2 ?sum]
    }
    ```
    becomes:
    ```
    match {
        goal [add: * ?num2 ?count ?sum] when (?count != ?num2)
    }
    ```

### Fixed

- Update [pyactr to 0.3.1](https://pypi.org/project/pyactr/) to fix compatibility problems with Python 3.10.
- Several amod lexing issues were fixed:
  - Invalid section names would hang gactar. ([#181](https://github.com/asmaloney/gactar/pull/181))
  - An amod file ending in a comment without a newline would hang gactar. ([#184](https://github.com/asmaloney/gactar/pull/184))
  - A malformed comment like "/ Some comment" would hang gactar.
  - Fixed handling of numbers (and errors with numbers) such as: `+., +.9, 0., .42.5`.
- Fixed the "Load Example" icon on Safari. ([#189](https://github.com/asmaloney/gactar/pull/189))

## [0.7.0](https://github.com/asmaloney/gactar/releases/tag/v0.7.0) - 2022-06-06

### Added

- The `-env` option will let you use different virtual environments (the default is `./env` which is next to the gactar executable).
- (_Linux_) Setup will now try to automatically download and install the [SBCL Lisp compiler](http://www.sbcl.org).
- Added _max_spread_strength_ config option to declarative **memory**. This turns on the spreading activation calculation & sets the maximum associative strength. ([#141](https://github.com/asmaloney/gactar/pull/141))
- Added _instantaneous_noise_ config option to declarative **memory**. This turns on the activation noise calculation & sets instantaneous noise. ([#162](https://github.com/asmaloney/gactar/pull/162))
- Added _spreading_activation_ config option to **goal**. This only takes effect if spreading activation is turned on via _max_spread_strength_ (see above). ([#148](https://github.com/asmaloney/gactar/pull/148))
- Added _trace_activations_ config option to **gactar**. This turns on detailed info about activations if available (currently _pyactr_ and _vanilla_ support it). ([#160](https://github.com/asmaloney/gactar/pull/160))
- Added documentation for which modules are available and their configuration options. (See [amod Config](./doc/amod%20Config.md).)

### Changed

- No longer need to run `source ./env/bin/activate` to activate the Python virtual environment. gactar will set the variables itself. ([#130](https://github.com/asmaloney/gactar/pull/130))
- Don't create md5 files with the releases.
- Rename `darwin` to `macOS` in releases.
- Don't try to install _pyactr_ if running Python 3.10+. It is currently not supported. ([see issue #137](https://github.com/asmaloney/gactar/issues/137))
- (_Windows_) Improved setup script: create symlink for _python3_ and fix _activate_ script path. ([@ren-oz](https://github.com/ren-oz)) ([#149](https://github.com/asmaloney/gactar/pull/149))

### Fixed

- Use "." instead of "source" in `setup.sh` since we are using "sh". This was breaking on Linux. ([#135](https://github.com/asmaloney/gactar/pull/135))
- Clarify some documentation.
- Generated Python code now conforms to [PEP8](https://peps.python.org/pep-0008/) style. ([#157](https://github.com/asmaloney/gactar/pull/157))

## [0.6.0](https://github.com/asmaloney/gactar/releases/tag/v0.6.0) - 2022-05-31

### Added

- Added a _warning_ level for issues. ([#108](https://github.com/asmaloney/gactar/pull/108))
- Frameworks can now validate the parsed code before running. This lets us return issues on a per-framework basis. ([#112](https://github.com/asmaloney/gactar/pull/112))
- Output a warning when no **goal** is available - either directly or in the initializers. ([#116](https://github.com/asmaloney/gactar/pull/116))
- Output the initial **goal** as info before running. ([#117](https://github.com/asmaloney/gactar/pull/117))
- Added config option for **procedural** module: _default_action_time_ is the time it takes to fire a production (seconds). ([#122](https://github.com/asmaloney/gactar/pull/122))

### Changed

- Cleaned up the `/api/run` return structure. ([#109](https://github.com/asmaloney/gactar/pull/109))
- `/api/run` now returns the generated code even if the run failed.
- Adjusted the **memory** module config options:
  - Added a _latency_exponent_ option.
  - Rename _latency_ to _latency_factor_.
  - Rename _threshold_ to _retrieval_threshold_.
  - Remove _max_time_ (may be able to add it back later).
  - Added some range checks.
  - Warn per-framework about any unsupported options.
  - Turn on **pyactr**'s _subsymbolic_ option to be in line with **vanilla**
- Replaced "anonymous variable" (`?`) with a wildcard character (`*`). ([#123](https://github.com/asmaloney/gactar/pull/123))

### Fixed

- In the web UI, only use syntax highlighting on variables if they are within square brackets.

## [0.5.0](https://github.com/asmaloney/gactar/releases/tag/v0.5.0) - 2022-05-26

### Added

- New command line option `temp` to specify where to generate the intermediate code files. If not specified, it defaults to `./gactar-temp` in the directory gactar was run from. The directory will be created if it does not exist. ([#94](https://github.com/asmaloney/gactar/pull/94))
- Web UI now allows the user to select which frameworks to run from the ones available on the server. ([#100](https://github.com/asmaloney/gactar/pull/100))
- Added TypeScript interfaces for all endpoints (in `api.ts`).
- Added new `/api/frameworks` endpoint to get info on frameworks available on the server. ([#99](https://github.com/asmaloney/gactar/pull/99))
- The `/api/run` endpoint now accepts an optional list of frameworks to run. If not specified, it will run on all available frameworks. ([#97](https://github.com/asmaloney/gactar/pull/97))
- The return data for `/api/run` now includes the full path to the intermediate code file in the property `filePath`.
- Added column numbers to error output. ([#102](https://github.com/asmaloney/gactar/pull/102))
- Added extra checks on patterns for valid chunk names and number of slots.

### Changed

- Use [camelCase](https://en.wikipedia.org/wiki/Camel_case) for all returned properties in the API.
- Clean up API TypeScript interfaces.

### Fixed

- When running as a web server, always create temp folder before a run in case it was removed. ([#107](https://github.com/asmaloney/gactar/pull/107))

## [0.4.0](https://github.com/asmaloney/gactar/releases/tag/v0.4.0) - 2022-05-20

### Added

- New command line options to support the new [gactar VS Code extension](https://marketplace.visualstudio.com/items?itemName=asmaloney.gactar) ([source here](https://github.com/asmaloney/gactar-vscode)).
  - `--output` (or `-o`) specifies where to put the intermediate source files. Defaults to "./".
  - `--run` (or `-r`) tells gactar to run the models after generating the code.

## [0.3.0](https://github.com/asmaloney/gactar/releases/tag/v0.3.0) - 2022-03-15

### Changed

- All endpoints are now prefixed by `/api`. This allows us to control the routes better in the web interface. ([#83](https://github.com/asmaloney/gactar/pull/83))

- Web development & build environment migrated from [vue-cli](https://cli.vuejs.org/) to [vite](https://vitejs.dev/). It is faster & reduces our dependencies. ([#84](https://github.com/asmaloney/gactar/pull/84))

- Convert web interface to use TypeScript. ([#86](https://github.com/asmaloney/gactar/pull/86))

- Update vanilla ACT-R to version 7.27.0 (from 15 Sep 2021).

- Change underlying `ccm` code from [CCMSuite3](https://github.com/CarletonCognitiveModelingLab/CCMSuite3) to [python_actr](https://github.com/CarletonCognitiveModelingLab/python_actr). The python_actr code was extracted from CCMSuite3 and now has a pip package to make installation easier.

  - **Naming note:** When gactar was written, it used [CCMSuite3](https://github.com/CarletonCognitiveModelingLab/CCMSuite3) and it was referred to throughout gactar as `ccm`. Instead of changing everything to refer to `python_actr` I've decided to leave it as `ccm`. This helps avoid confusion between `python_actr` and `pyactr`.

- Update all underlying dependencies (both go and npm).

## [0.2.0](https://github.com/asmaloney/gactar/releases/tag/v0.2.0) - 2021-11-17

### Added

- Added optional _authors_ field to the _model_ section. ([#54](https://github.com/asmaloney/gactar/pull/54)) It is a list of strings.

  Example:

  ```
  authors {
   	'Andy Maloney <andy@example.com>'
   	'Hiro Protagonist <hiro@example.com>'
  }
  ```

- Generated source files now include the gactar version which was used to generate them in the comments at the top. ([#78](https://github.com/asmaloney/gactar/pull/78))

- Added new web API endpoints for creating sessions, and compiling &amp; running models. These are intended to be used by other software to compile and run amod models using gactar running as a server. See the [Web API documentation](<doc/Web API.md>) for details.

- Added [documentation](<doc/Web API.md>) for existing web endpoints.

### Changed

- Unused variables now produce an error. ([#58](https://github.com/asmaloney/gactar/pull/58))

  ```
  ERROR: variable ?blat is not used - should be simplified to '?' (line 9)
  ```

- Anonymous variables ("?") in set statements now produce an error. ([#59](https://github.com/asmaloney/gactar/pull/59))

  ```
  do {
    set goal.thing to ?
    set goal to [foo: ?]
  }
  ```

  This will result in:

  ```
  ERROR: cannot set 'goal.thing' to anonymous var ('?') in production 'start' (line 10)
  ```

- Anonymous variables ("?") in print statements now produce an error. ([#60](https://github.com/asmaloney/gactar/pull/60))

  ```
  do {
    print ?
  }
  ```

  This will result in:

  ```
  ERROR: cannot print anonymous var ('?') in production 'start' (line 9)
  ```

- Compound variables ("?foo!?bar") in set statements now produce an error. ([#63](https://github.com/asmaloney/gactar/pull/63))

  ```
  do {
    set goal to [foo: ?foo!?bar]
  }
  ```

  This will result in:

  ```
  ERROR: cannot set 'goal.thing' to compound var in production 'start' (line 10)
  ```

- Multiple _recall_ statements in a production now produce an error. ([#69](https://github.com/asmaloney/gactar/pull/69))

  ```
  do {
    recall [foo: ?next ?]
    recall [foo: ? ?next]
  }
  ```

  This will result in:

  ```
  ERROR: only one recall statement per production is allowed in production 'start' (line 12)
  ```

- **pyactr**
  - Turn off _subsymbolic_ on the model as it is not necessary for what we are doing at the moment. ([#68](https://github.com/asmaloney/gactar/pull/68))
  - Clear the retrieval buffer before trying to fill it with a recall statement. This forces the pyactr productions to work like the vanilla ACT-R ones. ([#68](https://github.com/asmaloney/gactar/pull/68))

### Fixed

- **all frameworks**

  - Only output a description comment in the generated code if the _description_ field is present in the amod file.

- **pyactr**

  - Generated code now handles printing of numbers and variables. ([#65](https://github.com/asmaloney/gactar/pull/65))

    It is still limited to one `print` per production ([#66](https://github.com/asmaloney/gactar/issues/66))

  - Fix _addition2_ example. ([#39](https://github.com/asmaloney/gactar/pull/39))

## [0.1.0](https://github.com/asmaloney/gactar/releases/tag/v0.1.0) - 2021-09-22

Initial release
