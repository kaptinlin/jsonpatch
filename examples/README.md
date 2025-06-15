# JSON Patch Examples

This directory contains practical examples demonstrating different JSON Patch usage patterns.

## Examples

### 1. Basic Operations (`basic-operations/`)
Demonstrates fundamental operations: `add`, `replace`, `remove`, and `test`.

```bash
cd basic-operations && go run main.go
```

### 2. Array Operations (`array-operations/`)
Shows how to work with arrays: adding elements, removing items, and updating array contents.

```bash
cd array-operations && go run main.go
```

### 3. Conditional Operations (`conditional-operations/`)
Illustrates safe updates using `test` operations for validation and optimistic locking.

```bash
cd conditional-operations && go run main.go
```

### 4. Batch Update (`batch-update/`)
Demonstrates efficient batch operations for updating multiple items at once.

```bash
cd batch-update && go run main.go
```

### 5. Copy and Move Operations (`copy-move-operations/`)
Shows data restructuring and migration using `copy` and `move` operations.

```bash
cd copy-move-operations && go run main.go
```

### 6. String Operations (`string-operations/`)
Illustrates text editing capabilities with string insertion operations.

```bash
cd string-operations && go run main.go
```

### 7. Error Handling (`error-handling/`)
Demonstrates error handling patterns, validation techniques, and recovery strategies.

```bash
cd error-handling && go run main.go
```

### 8. Mutate Option (`mutate-option/`)
Shows the Mutate option for performance optimization and its limitations in Go.

⚠️ **Note**: This example demonstrates a known issue where the Mutate option is not working as expected.

```bash
cd mutate-option && go run main.go
```
