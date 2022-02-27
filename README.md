# crt

crt is a cli tool to get [Certificate Transparency](https://en.wikipedia.org/wiki/Certificate_Transparency) logs of a given domain name.

## Usage

```
Usage: crt [options...] <domain name>

Options:
  -o <path> Output file path. Write to file instead of stdout.
  -e        Exclude expired certificates.
  -l <int>  Limit the number of results. (default: 1000) 
  -json     Turn results to JSON.
  -csv      Turn results to CSV.

Examples:
  crt example.com
  crt -o logs.json -json example.com
  crt -csv -o logs.csv -l 15 example.com
```