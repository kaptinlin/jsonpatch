package jsonpatch_test

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/op"
)

const compileCycleScenarioEnv = "JSONPATCH_COMPILE_CYCLE_SCENARIO"

func TestCompileOpsRejectsSelfReferentialComposite(t *testing.T) {
	if os.Getenv(compileCycleScenarioEnv) == "self" {
		and := op.NewAnd(nil)
		and.Operations = []op.PredicateOp{and}
		assertCompileCycleRejected(t, and)
		return
	}
	runCompileCycleSubprocess(t, "self")
}

func TestCompileOpsRejectsMutuallyReferentialComposites(t *testing.T) {
	if os.Getenv(compileCycleScenarioEnv) == "mutual" {
		and := op.NewAnd(nil)
		or := op.NewOr(nil)
		and.Operations = []op.PredicateOp{or}
		or.Operations = []op.PredicateOp{and}
		assertCompileCycleRejected(t, and)
		return
	}
	runCompileCycleSubprocess(t, "mutual")
}

func TestCompileOpsAcceptsSharedAcyclicCompositeAndClonesIt(t *testing.T) {
	operand := op.NewDefined([]string{"name"})
	shared := op.NewOr(nil, operand)
	root := op.NewAnd(nil, shared, shared)

	patch, err := jsonpatch.CompileOps(
		[]jsonpatch.Op{root},
		jsonpatch.WithCapabilities(jsonpatch.Predicate),
	)
	if err != nil {
		t.Fatalf("CompileOps() error = %v", err)
	}

	shared.Operations[0] = op.NewUndefined([]string{"name"})
	result, err := jsonpatch.Apply(patch, map[string]any{"name": "Ada"})
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if result == nil {
		t.Fatal("Apply() result is nil")
	}
}

func assertCompileCycleRejected(t *testing.T, operation jsonpatch.Op) {
	t.Helper()
	previousMaxStack := debug.SetMaxStack(256 << 10)
	defer debug.SetMaxStack(previousMaxStack)

	patch, err := jsonpatch.CompileOps(
		[]jsonpatch.Op{operation},
		jsonpatch.WithCapabilities(jsonpatch.Predicate),
	)
	if patch != nil {
		t.Fatalf("CompileOps() patch = %v, want nil", patch)
	}
	if !errors.Is(err, jsonpatch.ErrPayloadInvalid) {
		t.Fatalf("CompileOps() error = %v, want ErrPayloadInvalid", err)
	}
	if !errors.Is(err, op.ErrPredicateCycle) {
		t.Fatalf("CompileOps() error = %v, want ErrPredicateCycle cause", err)
	}
}

func runCompileCycleSubprocess(t *testing.T, scenario string) {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.count=1")
	command.Env = append(os.Environ(), compileCycleScenarioEnv+"="+scenario)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("cycle compile subprocess failed: %v\n%s", err, output)
	}
}
