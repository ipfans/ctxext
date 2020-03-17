# ctxext

`ctxext` is extensions of stdlib `context.Context` with compatibility.

## Status

Pre-Alpha stage.

API may changes in the future version.

## HOWTO

`ctxext.Context` is compatibility with `context.Context` interface.
You can use this struct as `context.Context`. But It supplies more
powerful storage for contextual data.

```go
    ctx := ctxext.New(nil)
    ctx.Set("data", data)
    val := ctx.Value("data")
```

## Install

```
go get -u -v github.com/ipfans/ctxext
```
