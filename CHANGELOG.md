# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Use environment variable to set configuration file path, with `--config` flag taking precedence if both are provided. ([#14](https://github.com/nantli/goodcommit/pull/14))

## [1.0.0]

### Added

- Scaffold the initial CLI program with `huh` library. ([#2](https://github.com/nantli/goodcommit/pull/2))
- Implement `type` and `scope` modules. ([#3](https://github.com/nantli/goodcommit/pull/3))
- Implement pages and checkpoints functionality. ([#4](https://github.com/nantli/goodcommit/pull/4))
- Add `greetings` module and enhance UX of existing modules. ([#5](https://github.com/nantli/goodcommit/pull/5))
- Implement `co-authors`, `description`, `body`, and `why` modules. ([#6](https://github.com/nantli/goodcommit/pull/6))
- Bring your own commiter functionality. ([#6](https://github.com/nantli/goodcommit/pull/6))
- Pin modules functionality. ([#8](https://github.com/nantli/goodcommit/pull/8))
- Add `logo` module. ([#8](https://github.com/nantli/goodcommit/pull/8))
- Add `signoff` module. ([#9](https://github.com/nantli/goodcommit/pull/9))
- Implement dependencies in modules to ensure module functionality is only available if its dependencies are met. ([#10](https://github.com/nantli/goodcommit/pull/10))
- Implement the `breakingmsg` module for detailed input on breaking changes, dependent on the `breaking` module being active. ([#10](https://github.com/nantli/goodcommit/pull/10))
- Add author and coauthors emojis to the body of the commit message. ([#12](https://github.com/nantli/goodcommit/pull/12))
- Add `--config` flag to specify a custom configuration file. ([#12](https://github.com/nantli/goodcommit/pull/12))
- Commit staged files using the formatted commit message. ([#12](https://github.com/nantli/goodcommit/pull/12))

[unreleased]: https://github.com/nantli/goodcommit/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/nantli/goodcommit/compare/v0.0.0...v1.0.0

