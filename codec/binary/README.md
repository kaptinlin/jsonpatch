# Binary Codec for JSON Patch+

The **binary** codec serializes JSON Patch+ operations directly to **MessagePack** using the same array shape as the **compact** codec. It skips any intermediate JSON/compact conversion, resulting in smaller payloads and higher performance.

## Features

- Direct MessagePack output – no temporary `compact` representation
- Full support for first-order JSON Patch, extended and predicate operations (composite predicates not yet implemented)
- Familiar API: `Encode` / `Decode`, aligning with the other codecs

## Quick Example

```go
import (
    "github.com/kaptinlin/jsonpatch/codec/binary"
    "github.com/kaptinlin/jsonpatch/op"
    "github.com/kaptinlin/jsonpatch/internal"
)

func main() {
    // Build operations
    ops := []internal.Op{
        op.NewAdd([]string{"/", "foo"}, "bar"),
    }

    codec := binary.Codec{}

    // Encode to MessagePack
    data, err := codec.Encode(ops)
    if err != nil {
        panic(err)
    }

    // Decode back to operations
    decoded, err := codec.Decode(data)
    if err != nil {
        panic(err)
    }

    _ = decoded // use decoded operations
}
```

> **TIP**: Fully round-trip compatible with the `json` and `compact` codecs – mix & match freely. 
