# rwgo
A command-line utility that uses the ~~srslang~~ anything-gos interpreter with
a bunch of built-in functions for common operations at Rosewood Media.

## Interpreter Specification

| Key | Value |
| --- | ----- |
| Interpreter     | anything-gos/interp_a |
| Syntax Frontend | gottagofast/ParseListSimple (no options) |
| Syntax Options: Quotes | single-quotes only (subject to change) |
| Syntax Options: Brackets | parenthesis only (subject to change) |
| Imports         | anything-gos/builtins_a |
| Defined Here    | rpc |

## Commands

### `froute host <public address> <directory>`

Hosts static files.

The following example will serve files in the current directory over HTTP on
port 8001:
```
rgo froute :8001 .
```

### `froute proxy <public address> <origin address>`

Reverse proxy another server, such as a host started using `froute host`.

The following example will serve `myvpn.example.com`'s files over a different
port.
```
rgo froute :8002 http://myvpn.example.com:8001
```
