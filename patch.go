// Package jsonpatch compiles and applies JSON Patch, predicate, and extended
// operations to Go values while preserving the caller's document type.
package jsonpatch

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/deepclone"
	"github.com/kaptinlin/jsonpointer"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
	oppkg "github.com/kaptinlin/jsonpatch/op"
)

var errNilOperation = errors.New("nil operation")

// Capability identifies an operation vocabulary that may be compiled.
type Capability uint64

const (
	// RFC6902 enables the core JSON Patch operations.
	RFC6902 Capability = 1 << iota
	// Predicate enables non-regex predicate operations.
	Predicate
	// RegexPredicate enables the matches predicate operation.
	RegexPredicate
	// Extended enables JSON Patch Extended operations.
	Extended
)

// AllCapabilities enables every operation vocabulary implemented by the package.
const AllCapabilities = RFC6902 | Predicate | RegexPredicate | Extended

// CompileOption configures patch compilation.
type CompileOption func(*compileOptions)

type compileOptions struct {
	capabilities  Capability
	createMatcher internal.CreateRegexMatcher
	codec         string
}

func defaultCompileOptions() compileOptions {
	return compileOptions{capabilities: RFC6902}
}

// WithCapabilities sets the operation vocabularies accepted during compilation.
func WithCapabilities(capabilities ...Capability) CompileOption {
	return func(o *compileOptions) {
		var enabled Capability
		for _, capability := range capabilities {
			enabled |= capability
		}
		o.capabilities = enabled
	}
}

// WithCompileMatcher sets the regex matcher factory used while decoding matches operations.
func WithCompileMatcher(createMatcher CreateRegexMatcher) CompileOption {
	return func(o *compileOptions) {
		o.createMatcher = createMatcher
	}
}

func buildCompileOptions(opts []CompileOption) compileOptions {
	options := defaultCompileOptions()
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

// JSONText marks a string as JSON text instead of a scalar string document.
type JSONText string

// Patch is a compiled, reusable operation sequence.
type Patch struct {
	ops []Op
}

type applyOptions struct {
	mutate bool
}

// Compile compiles Go-built operations with the default RFC 6902 capability.
func Compile(ops ...Op) (*Patch, error) {
	return CompileOps(ops)
}

// CompileOps compiles Go-built operations.
func CompileOps(ops []Op, opts ...CompileOption) (*Patch, error) {
	options := buildCompileOptions(opts)
	return compileOps(ops, options)
}

// CompileOperations compiles JSON-shaped Operation values.
func CompileOperations(operations []jsoncodec.Operation, opts ...CompileOption) (*Patch, error) {
	options := buildCompileOptions(opts)
	options.codec = "json"

	ops := make([]Op, len(operations))
	for i := range operations {
		decoded, err := jsoncodec.DecodeOperations([]internal.Operation{operations[i]}, internal.JSONPatchOptions{
			CreateMatcher: options.createMatcher,
		})
		if err != nil {
			return nil, newFieldError(
				ErrPayloadInvalid,
				i,
				operations[i].Op,
				operations[i].Path,
				operations[i].From,
				options.codec,
				err,
			)
		}
		ops[i] = decoded[0]
	}
	return compileOps(ops, options)
}

// CompileJSON compiles a JSON patch document.
func CompileJSON(data []byte, opts ...CompileOption) (*Patch, error) {
	options := buildCompileOptions(opts)
	options.codec = "json"

	var operations []map[string]any
	if err := json.Unmarshal(data, &operations); err != nil {
		return nil, newPayloadError(options.codec, err)
	}

	ops := make([]Op, len(operations))
	for i := range operations {
		decoded, err := jsoncodec.Decode([]map[string]any{operations[i]}, internal.JSONPatchOptions{
			CreateMatcher: options.createMatcher,
		})
		if err != nil {
			return nil, newFieldError(
				ErrPayloadInvalid,
				i,
				stringMapValue(operations[i], "op"),
				stringMapValue(operations[i], "path"),
				stringMapValue(operations[i], "from"),
				options.codec,
				err,
			)
		}
		ops[i] = decoded[0]
	}
	return compileOps(ops, options)
}

func stringMapValue(values map[string]any, key string) string {
	value, _ := values[key].(string)
	return value
}

func compileOps(ops []Op, options compileOptions) (*Patch, error) {
	compiled := make([]Op, len(ops))
	for i, operation := range ops {
		if operation == nil {
			return nil, newError(ErrPayloadInvalid, i, nil, options.codec, errNilOperation)
		}
		if err := operation.Validate(); err != nil {
			return nil, newError(ErrPayloadInvalid, i, operation, options.codec, err)
		}
		if !operationAllowed(operation, options.capabilities) {
			return nil, newError(ErrUnsupportedCapability, i, operation, options.codec, nil)
		}
		cloned, err := cloneCompiledOperation(operation, options)
		if err != nil {
			return nil, newError(ErrPayloadInvalid, i, operation, options.codec, err)
		}
		compiled[i] = cloned
	}
	return &Patch{ops: compiled}, nil
}

func cloneCompiledOperation(operation Op, options compileOptions) (Op, error) {
	jsonOp, ok := operation.(internal.JSONOp)
	if !ok {
		return nil, fmt.Errorf("operation %T cannot encode to JSON", operation)
	}
	jsonOperation, err := jsonOp.ToJSON()
	if err != nil {
		return nil, err
	}

	decoded, err := jsoncodec.DecodeOperations(
		[]internal.Operation{deepclone.Clone(jsonOperation)},
		internal.JSONPatchOptions{CreateMatcher: options.createMatcher},
	)
	if err != nil {
		return nil, err
	}
	return decoded[0], nil
}

func operationAllowed(operation Op, capabilities Capability) bool {
	switch name := string(operation.Op()); {
	case name == string(internal.OpMatchesType):
		return capabilities&RegexPredicate != 0
	case internal.IsJSONPatchOperation(name):
		return capabilities&RFC6902 != 0
	case internal.IsPredicateOperation(name):
		return capabilities&Predicate != 0
	case internal.IsJSONPatchExtendedOperation(name):
		return capabilities&Extended != 0
	default:
		return false
	}
}

// Len returns the number of compiled operations.
func (p *Patch) Len() int {
	if p == nil {
		return 0
	}
	return len(p.ops)
}

// Result is the typed result of applying a compiled patch.
type Result[T internal.Document] struct {
	Doc   T
	Steps []Step
}

// Step describes one applied operation.
type Step struct {
	index   int
	op      string
	path    string
	from    string
	old     any
	applied bool
}

// Index returns the operation index.
func (s *Step) Index() int {
	return s.index
}

// Op returns the operation name.
func (s *Step) Op() string {
	return s.op
}

// Path returns the operation target path.
func (s *Step) Path() string {
	return s.path
}

// From returns the operation source path when present.
func (s *Step) From() string {
	return s.from
}

// Old returns the previous value reported by the operation.
func (s *Step) Old() any {
	return s.old
}

// Applied reports whether the operation completed.
func (s *Step) Applied() bool {
	return s.applied
}

// Apply applies patch immutably to doc.
func Apply[T internal.Document](patch *Patch, doc T) (*Result[T], error) {
	if patch == nil {
		return nil, newPayloadError("", errors.New("nil patch"))
	}
	return applyCompiledByDocumentType(patch, doc, &applyOptions{})
}

// ApplyInPlace applies patch and stores the result back in doc.
func ApplyInPlace[T internal.Document](patch *Patch, doc *T) error {
	if patch == nil {
		return newPayloadError("", errors.New("nil patch"))
	}
	if doc == nil {
		return newPayloadError("", errors.New("nil document pointer"))
	}
	result, err := applyCompiledByDocumentType(patch, *doc, &applyOptions{mutate: true})
	if err != nil {
		return err
	}
	*doc = result.Doc
	return nil
}

func applyCompiledByDocumentType[T internal.Document](patch *Patch, doc T, options *applyOptions) (*Result[T], error) {
	switch value := any(doc).(type) {
	case JSONText:
		return applyJSONTextDocument(patch, value, doc, options)
	case []byte:
		return applyJSONBytesDocument(patch, value, doc, options)
	case string:
		return applyDirectDocument(patch, value, doc, options)
	case map[string]any:
		return applyDirectDocument(patch, value, doc, options)
	case nil:
		return applyStructLikeDocument(patch, doc, options)
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return applyDirectDocument(patch, value, doc, options)
	case []any:
		return applyDirectDocument(patch, value, doc, options)
	default:
		return applyCompiledByReflection(patch, doc, options)
	}
}

func applyCompiledByReflection[T internal.Document](patch *Patch, doc T, options *applyOptions) (*Result[T], error) {
	docValue := reflect.ValueOf(any(doc))
	if !docValue.IsValid() || (docValue.Kind() == reflect.Pointer && docValue.IsNil()) {
		return applyStructLikeDocument(patch, doc, options)
	}

	switch docValue.Type().Kind() {
	case reflect.String:
		return applyDirectDocument(patch, any(doc), doc, options)
	case reflect.Slice:
		if docValue.Type().Elem().Kind() == reflect.Uint8 {
			bytes, ok := docValue.Convert(reflect.TypeFor[[]byte]()).Interface().([]byte)
			if !ok {
				return nil, conversionError(doc, nil)
			}
			return applyJSONBytesDocument(patch, bytes, doc, options)
		}
		return applyDirectDocument(patch, any(doc), doc, options)
	case reflect.Interface:
		return applyDirectDocument(patch, any(doc), doc, options)
	default:
		return applyStructLikeDocument(patch, doc, options)
	}
}

func applyJSONTextDocument[T internal.Document](patch *Patch, doc JSONText, original T, options *applyOptions) (*Result[T], error) {
	var parsed any
	if err := json.Unmarshal([]byte(doc), &parsed); err != nil {
		return nil, newPayloadError("json", err)
	}

	resultDoc, opResults, err := patch.apply(parsed, options)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(resultDoc)
	if err != nil {
		return nil, conversionError(original, err)
	}
	return resultFromRaw(patch, string(resultBytes), opResults, original)
}

func applyJSONBytesDocument[T internal.Document](patch *Patch, doc []byte, original T, options *applyOptions) (*Result[T], error) {
	var parsed any
	if err := json.Unmarshal(doc, &parsed); err != nil {
		return nil, newPayloadError("json", err)
	}

	resultDoc, opResults, err := patch.apply(parsed, options)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(resultDoc)
	if err != nil {
		return nil, conversionError(original, err)
	}
	return resultFromRaw(patch, resultBytes, opResults, original)
}

func applyStructLikeDocument[T internal.Document](patch *Patch, doc T, options *applyOptions) (*Result[T], error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, conversionError(doc, err)
	}

	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, conversionError(doc, err)
	}

	resultDoc, opResults, err := patch.apply(parsed, options)
	if err != nil {
		return nil, err
	}

	resultData, err := json.Marshal(resultDoc)
	if err != nil {
		return nil, conversionError(doc, err)
	}

	var result T
	if err := json.Unmarshal(resultData, &result); err != nil {
		return nil, conversionError(doc, err)
	}
	return &Result[T]{Doc: result, Steps: stepsFromResults(patch.ops, opResults)}, nil
}

func applyDirectDocument[T internal.Document](patch *Patch, working any, original T, options *applyOptions) (*Result[T], error) {
	resultDoc, opResults, err := patch.apply(working, options)
	if err != nil {
		return nil, err
	}
	return resultFromRaw(patch, resultDoc, opResults, original)
}

func resultFromRaw[T internal.Document](patch *Patch, resultDoc any, opResults []internal.OpResult[any], original T) (*Result[T], error) {
	result, err := convertResult(resultDoc, original)
	if err != nil {
		return nil, err
	}
	return &Result[T]{Doc: result, Steps: stepsFromResults(patch.ops, opResults)}, nil
}

func convertResult[T internal.Document](resultDoc any, original T) (T, error) {
	var zero T
	if resultDoc == nil {
		if reflect.TypeFor[T]().Kind() == reflect.Interface {
			return zero, nil
		}
		return zero, conversionError(original, fmt.Errorf("operation resulted in null value"))
	}

	if result, ok := resultDoc.(T); ok {
		return result, nil
	}

	resultValue := reflect.ValueOf(resultDoc)
	targetType := reflect.TypeFor[T]()
	if resultValue.IsValid() && resultValue.Type().ConvertibleTo(targetType) {
		return resultValue.Convert(targetType).Interface().(T), nil
	}

	return zero, conversionError(original, nil)
}

func conversionError(doc any, cause error) error {
	if cause == nil {
		cause = fmt.Errorf("failed to convert result back to type %T", doc)
	}
	return &Error{
		kind:  ErrConversionFailed,
		index: -1,
		cause: cause,
	}
}

func (p *Patch) apply(doc any, options *applyOptions) (any, []internal.OpResult[any], error) {
	workingDoc := doc
	if !options.mutate {
		workingDoc = deepclone.Clone(doc)
	}

	results := make([]internal.OpResult[any], 0, len(p.ops))
	for i, operation := range p.ops {
		if operation == nil {
			return nil, nil, newError(ErrPayloadInvalid, i, nil, "", errNilOperation)
		}
		opResult, err := operation.Apply(workingDoc)
		if err != nil {
			return nil, nil, newError(kindForApplyError(err), i, operation, "", err)
		}
		workingDoc = opResult.Doc
		results = append(results, opResult)
	}
	return workingDoc, results, nil
}

func kindForApplyError(err error) error {
	switch {
	case errors.Is(err, oppkg.ErrTestFailed),
		errors.Is(err, oppkg.ErrTestOperationFailed),
		errors.Is(err, oppkg.ErrTestOperationNumberStringMismatch),
		errors.Is(err, oppkg.ErrTestOperationStringNotEquivalent),
		errors.Is(err, oppkg.ErrStringMismatch),
		errors.Is(err, oppkg.ErrSubstringMismatch),
		errors.Is(err, oppkg.ErrStringLengthMismatch),
		errors.Is(err, oppkg.ErrContainsMismatch),
		errors.Is(err, oppkg.ErrComparisonFailed),
		errors.Is(err, oppkg.ErrOperationFailed),
		errors.Is(err, oppkg.ErrDefinedTestFailed),
		errors.Is(err, oppkg.ErrUndefinedTestFailed),
		errors.Is(err, oppkg.ErrAndTestFailed),
		errors.Is(err, oppkg.ErrOrTestFailed),
		errors.Is(err, oppkg.ErrNotTestFailed):
		return ErrTestFailed
	case errors.Is(err, oppkg.ErrTypeMismatch):
		return ErrTypeMismatch
	default:
		return ErrRuntimeConflict
	}
}

func stepsFromResults(ops []Op, results []internal.OpResult[any]) []Step {
	steps := make([]Step, len(results))
	for i := range results {
		operation := ops[i]
		step := Step{
			index:   i,
			op:      string(operation.Op()),
			path:    jsonpointer.Format(operation.Path()...),
			old:     results[i].Old,
			applied: true,
		}
		if from, ok := operation.(interface{ From() []string }); ok {
			step.from = jsonpointer.Format(from.From()...)
		}
		steps[i] = step
	}
	return steps
}
