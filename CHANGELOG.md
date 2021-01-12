# Changelog

## 0.9.3: Allow `\n` as an initial character

This release allows `\n` (ASCII 10) instead of the `\0` as an initial character. This is needed because in TTY mode Kubernetes "eats" non-printable characters.

## 0.9.2: Package signing

This release introduces package signing for `deb`, `rpm`, and `apk` packages.

## 0.9.1: Bugfixing exec

The previous release of the agent failed to pass the executable path as the first parameter of the executable being run. This caused the execution to fail.

## 0.9.0: Initial release

This is the initial release.