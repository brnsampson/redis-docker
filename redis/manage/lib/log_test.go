package lib

import (
	"testing"
	"log"
	"bytes"
)

func TestEnableDebug(t *testing.T) {
	var b bytes.Buffer

	rl := RedisLog{ debugOn: false,
					Logger: log.New(&b, "", 0)}

	expected := "test"
	rl.Debug(expected)
	got := b.String()

	if got != "" {
		t.Errorf("Debug message printed with DebugOn == false")
	}

	expected = "Debug logging enabled\n"
	rl.EnableDebug()
	got = b.String()

	if got != expected {
		t.Errorf("expected EnableDebug() to output %q, got %q", expected, got)
	}
}

func TestDisableDebug(t *testing.T) {
	var b bytes.Buffer

	rl := RedisLog{ debugOn: true,
					Logger: log.New(&b, "", 0)}

	expected := "Debug logging disabled\n"
	rl.DisableDebug()
	got := b.String()

	if got != expected {
		t.Errorf("expected DisableDebug() to output %q, got %q", expected, got)
	}
}

func TestDebug(t *testing.T) {
	var b bytes.Buffer

	rl := RedisLog{ debugOn: false,
					Logger: log.New(&b, "", 0)}

	unexpected := "test"
	rl.Debug(unexpected)
	got := b.String()

	if got != "" {
		t.Errorf("Debug message printed with DebugOn == false")
	}

	rl.EnableDebug()
	_ = b.String()

	expected := "test"
	rl.Debug(expected)
	got = b.String()

	if got != expected {
		t.Errorf("expected Debug() to output %q, got %q", expected, got)
	}
}
