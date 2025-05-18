# Dockerfile parser

Basic dockerfile parser

Note: This is more of a PoC than a production ready parser

## Known Issues

- [x] Multiline commands
- [x] Optional array parsing where supported
- [ ] Comments (inline do not get detected and comment lines just get thrown out before the lexer)
- [ ] Stage reference detection
- [ ] Arg behaviour parsing when not actively setting a value
- [ ] Several edge cases and special scenarios for a few instructions
- [ ] Dockerfiles not starting with FROM (apparently they can start with ARG)
- [ ] Support for dockerfiles with EOF blocks  (this is probably the worst offender with this: https://github.com/apache/airflow/blob/main/Dockerfile)
- [x] Support for expose statement exposing multiple ports

## Benchamarking

```sh
go test -bench=. ./...
```
