# QueryXSS

QueryXSS is a tool to test for reflected inputs in the response.

**Beware:** This tool is still in development, so you can expect bugs.

## Scanners

The list is not final, but it's a good start. We may add more scanners in the future.

| ID | Description | Implemented |
| --- | --- | --- |
| simple-query | SimpleQuery scanner requests the URL as it is, with no modifications. It looks for query keys and values being reflected. | YES |
| simple-headers | SimpleHeaders scanner requests the URL as it is, with no modifications. It looks for headers being reflected. | NO |

## Usage

```
Usage:
  queryxss [flags]

Flags:
  -d, --debug                Enable debug mode
  -f, --file string          File with URLs to scan
  -H, --header stringArray   Headers to send with the request (specify multiple times)
                             Example: -H 'X-Forwarded-For: 127.0.0.1' -H 'X-Random: 1234'
  -h, --help                 help for queryxss
  -n, --no-color             Disable color output
  -r, --rate-limit uint      Number of requests per second (default 25)
  -s, --silent               Outputs only errors and the results
```

## Install

### Using go install

Make sure you have [Go installed and configured](https://go.dev/doc/install).

```bash
go install github.com/vitorfhc/hacks/queryxss
```

### Manual install

```bash
git clone github.com/vitorfhc/hacks
cd hacks/queryxss
go install
```