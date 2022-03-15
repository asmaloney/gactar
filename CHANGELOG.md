# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

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
