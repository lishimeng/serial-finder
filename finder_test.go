package sfinder

import (
	"context"
	"testing"
	"time"
)

var stub Finder
func TestFinder_Find(t *testing.T) {
	stub = New(FT2x)
	ports, err := stub.Find()
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, p := range ports {
		t.Logf("serial port: %s", p)
	}
}

func TestFinder_Monitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stub = New(FT2x)
	(stub).Monitor(MonitorConfig{
		Ctx:   ctx,
		Delay: time.Second*2,
		Build: func(name string) (i interface{}, err error) {
			i = name + "_ppp"
			t.Logf("create element:%s", i)
			return
		},
	})

	time.AfterFunc(time.Second*10, func() {
		t.Log("cancel")
		cancel()
	})

	t.Log("wait 10 seconds...")
	<-ctx.Done()
}

func TestFinder_Remove(t *testing.T) {
	var lastName string
	ctx, cancel := context.WithCancel(context.Background())
	stub = New(FT2x)
	(stub).Monitor(MonitorConfig{
		Ctx:   ctx,
		Delay: time.Second*2,
		Build: func(name string) (i interface{}, err error) {
			lastName = name
			i = name + "_ppp"
			t.Logf("create element:%s", i)
			return
		},
	})

	time.AfterFunc(time.Second*10, func() {
		t.Log("cancel")
		cancel()
	})

	time.AfterFunc(time.Second*5, func() {
		t.Logf("remove %s", lastName)
		stub.Remove(lastName)
	})

	t.Log("wait 10 seconds...")
	<-ctx.Done()
}