# golangci-lint integration

Testolint uses golangci-lint v2's [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/).
It requires a custom binary; adding YAML settings to a stock binary is not enough.
You need Go, Git, network access, curl and sh. See the
[verified compatibility](../README.md#verified-compatibility) before building.

## Build in your project

Create `.custom-gcl.yml` in your project directory:

```yaml
version: v2.14.0
name: golangci-lint
destination: ./bin
plugins:
  - module: github.com/nashabanov/testolint
    import: github.com/nashabanov/testolint/plugin/golangci
    version: latest
```

`latest` selects the latest published testolint release. Pin a published version
instead when you need reproducible builds.

Install the bootstrap tool and build the custom binary:

```sh
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./bin/bootstrap v2.14.0
./bin/bootstrap/golangci-lint custom -v
```

The resulting `bin/golangci-lint` contains testolint. Rebuild it after updating
the plugin or golangci-lint.

## Configure and run

Create `.golangci.yml`, or merge this into your v2 configuration:

```yaml
version: "2"
linters:
  default: standard
  enable:
    - testolint
  settings:
    custom:
      testolint:
        type: module
        description: Check Testo suite usage.
```

Run the custom binary from your project directory:

```sh
./bin/golangci-lint run
```

See [rule configuration](configuration.md#golangci-lint) for enable/disable
settings and the [driver contract](contract.md#driver-results) for exit codes.
For builds from a testolint checkout and integration verification, see
[Contributing](../CONTRIBUTING.md#plugin-integration).
