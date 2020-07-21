package tinylog

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"testing"
	"time"
)

type concurrentWriter struct {
	b  *bytes.Buffer
	mu sync.Mutex
}

func (cw *concurrentWriter) Write(p []byte) (n int, err error) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.b.Write(p)
}

func (cw *concurrentWriter) String() string {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	return cw.b.String()
}

func checkOutput(t *testing.T, b *bytes.Buffer, wantLevel string) {
	r := new(Record)
	err := json.Unmarshal(b.Bytes(), r)
	if err != nil {
		t.Error(err)
	}
	if r.Level != wantLevel {
		t.Errorf(`l.Info(test) = %s, want %s`, r.Level, wantLevel)
	}
}

func testTestTinyLoggerFactoryRace(t *testing.T, got, pattern, want string) {
	ok, err := regexp.MatchString(pattern, got)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf(`%s not found in resulting output`, want)
	}
}

func TestTemplate(t *testing.T) {
	ctx := context.TODO()
	b := new(bytes.Buffer)
	b.Write([]byte("\n"))
	lf := NewTinyLoggerFactory(b, String, time.Stamp)
	lf.SetLogLevel(Debug)
	l1 := lf.GetLogger(ctx, NilModule)
	l2 := lf.GetLogger(ctx, "TestMedium")
	l2.AddTag(ctx, "tool", "tinylog")
	l2.AddTag(ctx, "user", "me", "cat")
	l3 := lf.GetLogger(ctx, "TestLooooooooong")
	l1.Info("Hello World!")
	for i := 5; i > 0; i-- {
		l2.Debug("Hello World!")
		l2.Info("Hello World!")
		l2.Warn("Hello World!")
		l2.Error("Hello World!")
		l1.Error("Hello World!")
		l1.Info("Hello World!")
		l3.Warn("Hello World!")
		l3.Error("Hello World!")
		l3.Info("Hello World!")
	}
	t.Log(b.String())
}

func TestDefaultLogLevel(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.Debug("test")
	got := b.String()
	want := ""
	if got != want {
		t.Errorf("l.Debug(test) = %q, want %q; LogLevel: default level", got, want)
	}
}

func TestWarnLogLevel(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.SetLogLevel(Warn)
	l.Debug("test")
	l.Info("test")
	got := b.String()
	want := ""
	if got != want {
		t.Errorf("got %q, want %q; LogLevel: %d", got, want, Warn)
	}
}

func TestErrLogLevel(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.SetLogLevel(Error)
	l.Debug("test")
	l.Info("test")
	l.Warn("test")
	got := b.String()
	want := ""
	if got != want {
		t.Errorf("got %q, want %q; LogLevel: %d", got, want, Error)
	}
}

func TestFatalLogLevel(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.SetLogLevel(Fatal)
	l.Debug("test")
	l.Info("test")
	l.Warn("test")
	l.Error("test")
	got := b.String()
	want := ""
	if got != want {
		t.Errorf("got %q, want %q; LogLevel: %d", got, want, Fatal)
	}
}

func TestNoneLogLevel(t *testing.T) {
	b := new(bytes.Buffer)
	if os.Getenv("BE_CRASHER") == "1" {
		l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
		l.SetLogLevel(None)
		l.Debug("test")
		l.Info("test")
		l.Warn("test")
		l.Error("test")
		l.Fatal("test")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestNoneLogLevel")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		got := b.String()
		want := ""
		if got != want {
			t.Errorf("got %q, want %q; LogLevel: %d", got, want, None)
		}
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)

}

func TestDebug(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.SetLogLevel(Debug)
	l.Debug("test")
	checkOutput(t, b, "DEBUG")
}

func TestInfo(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.Info("test")
	checkOutput(t, b, "INFO")
}

func TestWarn(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.Warn("test")
	checkOutput(t, b, "WARN")
}

func TestError(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.Error("test")
	checkOutput(t, b, "ERROR")
}

func TestTinyLoggerFactoryRace(t *testing.T) {
	ctx := context.TODO()
	b := &concurrentWriter{b: new(bytes.Buffer)}
	b.Write([]byte("\n"))
	lf := NewTinyLoggerFactory(b, String, time.RubyDate)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger(ctx, "Test")
		l.Info("Test")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger(ctx, "Test1")
		l.Info("Test")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger(ctx, "Test")
		l.Info("Test")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger(ctx, "Test2")
		l.Info("Test")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		lf.SetLogLevel(Debug)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger(ctx, "Test3")
		l.Info("Test")
	}()
	wg.Wait()
	got := b.String()

	testTestTinyLoggerFactoryRace(t, got, `Test`, "Test")
	testTestTinyLoggerFactoryRace(t, got, `Test1`, "Test1")
	testTestTinyLoggerFactoryRace(t, got, `Test2`, "Test2")
	testTestTinyLoggerFactoryRace(t, got, `Test3`, "Test3")
}

func TestContextCancelTags(t *testing.T) {
	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	l.AddTag(ctx, "user", "me", "cat")
	cancel()
	time.Sleep(time.Millisecond)
	l.Info("test")
	r := new(Record)
	err := json.Unmarshal(b.Bytes(), r)
	if err != nil {
		t.Error(err)
	}
	if tags := r.Tags["user"]; r.Message != "test" || len(tags) > 0 {
		t.Errorf(`l.Info(test) without Tags = {Message: "%s", Tags: %v}, want {Message: "test", Tags: []}`, r.Message, tags)
	}
}

func TestTags(t *testing.T) {
	ctx := context.TODO()
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, JSON, NilModule, time.RubyDate)
	testTags(l, ctx, b, "user", "me", t)
}

func TestIndependentTagsByContext(t *testing.T) {
	go func() {
		ctx := context.TODO()
		b := new(bytes.Buffer)
		lf := NewTinyLoggerFactory(b, JSON, time.RubyDate)
		l := lf.GetLogger(ctx, NilModule)
		testTags(l, ctx, b, "user", "me", t)
	}()
	go func() {
		ctx := context.TODO()
		b := new(bytes.Buffer)
		lf := NewTinyLoggerFactory(b, JSON, time.RubyDate)
		l := lf.GetLogger(ctx, NilModule)
		testTags(l, ctx, b, "user", "cat", t)
	}()
}

func testTags(l Logger, ctx context.Context, b *bytes.Buffer, key, value string, t *testing.T) {
	l.AddTag(ctx, "user", "me", "cat")
	l.Info("test")

	r := new(Record)

	if err := json.Unmarshal(b.Bytes(), r); err != nil {
		t.Error(err)
	}

	if tags := r.Tags[key]; tags[0] != value {
		t.Errorf(`l.Info(test) with Tags "%s": ["%s"] = {Tags: map[%s]}, want {Tags: map[%s][%s]}`, key, value, tags, key, value)
	}
}
