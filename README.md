# codesee-deps-go

[![Version](https://img.shields.io/badge/version-v0.0.0-green.svg)](https://github.com/Codesee-io/codesee-deps-go/releases)

A command line tool that gives you a list of all usages between files within a
project, while ignoring external dependencies.

This is used by the [`codesee` CLI](https://www.npmjs.com/package/codesee) to
generate accurate data for Golang projects.

Currently, this utility only works on projects that use Go modules.

## Usage

```sh
codesee-deps-go <directory>
```

This will output a JSON array of objects with `from` and `to` keys. For example:

```json
[
  {
    "from": "cmd/api/main.go",
    "to": "pkg/server/server.go"
  },
  {
    "from": "cmd/api/main.go",
    "to": "pkg/signals/signals.go"
  },
  {
    "from": "pkg/server/server.go",
    "to": "pkg/handlers/handlers.go"
  },
  {
    "from": "pkg/signals/signals_test.go",
    "to": "pkg/signals/signals.go"
  }
]
```

## Development

### Building

```sh
make build
```

### Testing

```sh
make test
```

### Running Locally

```sh
go run ./cmd/deps <directory>
```

You can either use this project, the test `simple-repo` project, or any other Go
project that you have for testing purposes.

```sh
go run ./cmd/deps .
go run ./cmd/deps ./pkg/testdata/simple-repo
go run /path/to/other/project
```

### Releasing

To make a new release, you need to have a GitHub token with `repo` permissions.
You can generate one [here](https://github.com/settings/tokens/new). Once you
have it, you just need to run the following command:

```sh
GITHUB_TOKEN=<github_token> TAG=<new_tag> make release
# e.g.
GITHUB_TOKEN=TOKEN TAG=v1.0.0 make release
```

This does the following:

- Changes the version number in the `README.md`
- Commits the change with the `TAG` as the commit message
- Creates a tag with the specified `TAG`
- Pushes the commit and tag to GitHub
- Runs `goreleaser` to generate the binaries and uploads them to GitHub as a new
  release
