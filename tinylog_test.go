package tinylog

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"testing"
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
	b := new(bytes.Buffer)
	b.Write([]byte("\n"))
	lf := NewTinyLoggerFactory(b)
	l1 := lf.GetLogger(NilModule)
	l2 := lf.GetLogger("TestMedium")
	l3 := lf.GetLogger("TestLooooooooong")
	l1.Info("Hello World!")
	for i := 5; i > 0; i-- {
		l2.Info("Hello World!")
		l2.Warn("Hello World!")
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
	l := NewTinyLogger(b, NilModule)
	l.Debug("test")
	got := b.String()
	want := ""
	if got != want {
		t.Errorf("l.Debug(test) = %q, want %q; LogLevel: default level", got, want)
	}
}

func TestWarnLogLevel(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, NilModule)
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
	l := NewTinyLogger(b, NilModule)
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
	l := NewTinyLogger(b, NilModule)
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
		l := NewTinyLogger(b, NilModule)
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
	l := NewTinyLogger(b, NilModule)
	l.SetLogLevel(Debug)
	l.Debug("test")
	got := b.String()
	want := `DEBUG`
	ok, err := regexp.MatchString(want, got)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf(`l.Debug(test) = %q, want "[debug]..."`, got)
	}
}

func TestInfo(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, NilModule)
	l.Info("test")
	got := b.String()
	want := `INFO`
	ok, err := regexp.MatchString(want, got)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf(`l.Info(test) = %q, want "[info]..."`, got)
	}
}

func TestWarn(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, NilModule)
	l.Warn("test")
	got := b.String()
	want := `WARN`
	ok, err := regexp.MatchString(want, got)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf(`l.Warn(test) = %q, want "[warn]..."`, got)
	}
}

func TestError(t *testing.T) {
	b := new(bytes.Buffer)
	l := NewTinyLogger(b, NilModule)
	l.Error("test")
	got := b.String()
	want := `ERROR`
	ok, err := regexp.MatchString(want, got)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Errorf(`l.Error(test) = %q, want "[error]..."`, got)
	}
}

func TestTinyLoggerFactoryRace(t *testing.T) {
	b := &concurrentWriter{b: new(bytes.Buffer)}
	b.Write([]byte("\n"))
	lf := NewTinyLoggerFactory(b)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger("Test")
		l.Info("Test")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger("Test1")
		l.Info("Test")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger("Test")
		l.Info("Test")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		l := lf.GetLogger("Test2")
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
		l := lf.GetLogger("Test3")
		l.Info("Test")
	}()
	wg.Wait()
	got := b.String()

	testTestTinyLoggerFactoryRace(t, got, `\|Test\|`, "|Test|")
	testTestTinyLoggerFactoryRace(t, got, `\|Test1\|`, "|Test1|")
	testTestTinyLoggerFactoryRace(t, got, `\|Test2\|`, "|Test2|")
	testTestTinyLoggerFactoryRace(t, got, `\|Test3\|`, "|Test3|")
}
