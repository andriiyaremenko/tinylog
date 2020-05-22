package tinylog

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
)

func TestTemplateJSON(t *testing.T) {
	b := new(bytes.Buffer)
	b.Write([]byte("\n"))
	b.Write([]byte("["))
	lf := NewTinyLoggerFactory(NewJSONWriter(b))
	l1 := lf.GetLogger(NilModule)
	l2 := lf.GetLogger("TestMedium")
	l3 := lf.GetLogger("TestLooooooooong")
	l1.Info("Hello World!")
	for i := 5; i > 0; i-- {
		b.Write([]byte(","))
		b.Write([]byte("\n"))
		l2.Info("Hello World!")
		b.Write([]byte(","))
		b.Write([]byte("\n"))
		l2.Warn("Hello World!")
		b.Write([]byte(","))
		b.Write([]byte("\n"))
		l1.Err("Hello World!")
		b.Write([]byte(","))
		b.Write([]byte("\n"))
		l1.Info("Hello World!")
		b.Write([]byte(","))
		b.Write([]byte("\n"))
		l3.Warn("Hello World!")
		b.Write([]byte(","))
		b.Write([]byte("\n"))
		l3.Err("Hello World!")
		b.Write([]byte(","))
		b.Write([]byte("\n"))
		l3.Info("Hello World!")
	}
	b.Write([]byte("]"))
	t.Log(b.String())
	err := json.Unmarshal(b.Bytes(), new([]record))
	if err != nil {
		t.Error(err)
	}
}

func TestTinyLoggerFactoryRaceJSON(t *testing.T) {
	b := &concurrentWriter{b: new(bytes.Buffer)}
	b.Write([]byte("\n"))
	lf := NewTinyLoggerFactory(NewJSONWriter(b))
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
