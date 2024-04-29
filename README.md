# gn-coval-index

## Avro schema update
In case you want to modify the existing avro schema, after you finish your changes, you need to re-generate the corresponding code, by:

1. Get code generator binary
```bash
go install github.com/elodina/go-avro/codegen@v0.1.0.0
```

2. Run `go generate` from `schema/codegen.go`
