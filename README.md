# Artie CLI

A command-line interface for interacting with the Artie API.

## Installation

You can install it via go install or download the binaries under [releases](https://https://github.com/artie-labs/artiecli/releases)

```bash
go install github.com/artie-labs/artiecli@latest
```

## Usage

> The Artie CLI requires an API key to authenticate with the Artie API. If you don't have this, you can generate one from our dashboard.

```bash
export ARTIE_API_KEY=your_api_key
```

### Commands

```bash
artiecli list-deployments

artiecli get-deployment --deployment-uuid UUID

artiecli cancel-deployment-backfill --deployment-uuid UUID --table-uuids UUID1,UUID2
```

## Development

This project is written in Go. To build from source:

```bash
git clone https://github.com/artie-labs/artiecli.git
cd artiecli
go build
```

## License

[MIT](LICENSE) 
