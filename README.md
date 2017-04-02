# spirali
**spirali** is a go based database migration tool.

[![GoDoc](https://godoc.org/github.com/kudohamu/spirali?status.svg)](https://godoc.org/github.com/kudohamu/spirali)
[![wercker status](https://app.wercker.com/status/bc0b4e372f5f0c80f0a88c0f91410273/s/master "wercker status")](https://app.wercker.com/project/byKey/bc0b4e372f5f0c80f0a88c0f91410273)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)



## Overview

* Have `create`, `up`, `down` command
* Be able to use in golang code (this is flexible more than the CLI)
* Customizable versioning
* Taking `bindata` into consideration

## Installation

```sh
$ go get github.com/kudohamu/spirali/cmd/spirali
```

You can check with the below command whether spirali install

```sh
$ spirali -h
```

## Usage

### config file

Prepare config file in accordance with below (file type is `.toml` ).

```toml
[dev]
dsn=""
driver="mysql" # only mysql now...
directory="" # relative path for directory of migration files from working directory.

[test]
dsn=""
driver="mysql"
directory=""
```

### Create

Create new migration files.

```sh
$ spirali --path=path/to/your/config/file --env=dev create create_user_table
```

### Up

Apply remaining migrations.

```sh
$ spirali --path=path/to/your/config/file --env=dev up
```

### Down

Roll back the latest migration.

```sh
$ spirali --path=path/to/your/config/file --env=dev down
```

## Use in golang code

```go
spirali.Create(vg VersionG, name string, config *Config, metadata *MetaData) (*MetaData, error)

spirali.Up(metadata *MetaData, config *Config, driver Driver, readable Readable) error

spirali.Down(metadata *MetaData, config *Config, driver Driver, readable Readable) error
```

(TODO: More friendly code sample)

### Driver?

You can generate driver with below code.

```go
driver, err := spirali.NewDriver(config)
```

### VersionG?

You can customize versioning of migration files.

```go
type VersionG interface {
  GenerateNextVersion() (uint64, error)                        // generates next version.
  IsSmall(targetVersion uint64, comparisonVersion uint64) bool // returns whether targetVersion is smaller than comparisonVersion.
}
```

Now spirali prepares two generators.

* TimestampBasedVersionG (versioning with `YYYYMMDDhhmmss` format)
* IncrementalVersionG (versioning with `1`, `2`, `3`, ... format)

(Be able to use only `TimestampBasedVersionG` in CLI now.)

### Readable?

Spirali abstracts directory to read migration files.

```go
type Readable interface {
  Read(path string) ([]byte, error)
}
```

In CLI, spirali uses `spirali.Dir`. If you want to read migration files from `bindata`, you can use `spirali.Bindata` in your golang code.

```go
readable := spirali.NewReadableFromBindata(bindata.Asset)
```

## TODO

* More drivers (`postgresql`, `SQLite3`)
* More useful commands

## More details

See [godoc](https://godoc.org/github.com/kudohamu/spirali).

## Contributing
Please feel free to submit issues, and send pull requests.

## License

Released under the [MIT License](https://github.com/kudohamu/spirali/blob/master/LICENSE).
