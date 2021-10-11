# Changelog

All notable changes to this project will be documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## 0.2.0 - (in progress)

### Added

- Added optional _authors_ field to the _model_ section. ([#54](https://github.com/asmaloney/gactar/pull/54)) It is a list of strings.

  Example:

  ```
  authors {
   	'Andy Maloney <andy@example.com>'
   	'Hiro Protagonist <hiro@example.com>'
  }
  ```

- Added check for unused variables. Output an info message to suggest changing them to anonymous variables. ([#58](https://github.com/asmaloney/gactar/pull/58))
  ```
  INFO: variable ?blat is not used - should be simplified to '?' (line 9)
  ```

### Changed

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

### Fixed

- pyactr generated code now handles printing of numbers and variables. ([#65](https://github.com/asmaloney/gactar/pull/65))

  It is still limited to one `print` per production ([#66](https://github.com/asmaloney/gactar/issues/66))

## [0.1.0](https://github.com/asmaloney/gactar/releases/tag/v0.1.0) - 2021-09-22

Initial release
