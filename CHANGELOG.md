# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [0.7.0] - (in progress)

### Added

- The `-env` option will let you use different virtual environments (the default is `./env` which is next to the gactar executable).
- (_Linux_) Setup will now try to automatically download and install the [SBCL Lisp compiler](http://www.sbcl.org).
- Added _max_spread_strength_ config option to declarative **memory**. This turns on the spreading activation calculation & sets the maximum associative strength. ([#141](https://github.com/asmaloney/gactar/pull/141))
- Added _spreading_activation_ config option to **goal**. This only takes effect if spreading activation is turned on via _max_spread_strength_ (see above). ([#148](https://github.com/asmaloney/gactar/pull/148))
- Added _trace_activations_ config option to **gactar**. This turns on detailed info about activations if available (currently _pyactr_ and _vanilla_ support it). ([#160](https://github.com/asmaloney/gactar/pull/160))

### Changed

- No longer need to run "source ./env/bin/activate" to activate the Python virtual environment. gactar will set the variables itself. ([#130](https://github.com/asmaloney/gactar/pull/130))
- Don't create md5 files with the releases.
- Rename "darwin" to "macOS" in releases.
- Don't try to install _pyactr_ if running Python 3.10+. It is currently not supported. ([see issue #137](https://github.com/asmaloney/gactar/issues/137))
- (_Windows_) Improved setup script: create symlink for _python3_ and fix _activate_ script path. (@ren-oz) ([#149](https://github.com/asmaloney/gactar/pull/149))

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
