# Dockerfile parser

Basic dockerfile parser

Note: This is more of a PoC than a production ready parser

## Known Issues

- [x] Multiline commands
- [x] Optional array parsing where supported
- [x] Support for expose statement exposing multiple ports
- [x] Stage reference detection
- [x] Dockerfiles not starting with FROM (apparently they can start with ARG)
- [x] Comments (inline do not get detected)
- [x] Arg behaviour parsing when not actively setting a value
- [x] Parser directives -> Are recognized and parsed into the ast...but as of now I dont actually do anything with them
- [ ] The full extend of heredoc (basics are implemented)
- [ ] Bash like variabe magic

## Benchmarking

```sh
go test -bench=. ./...
```
