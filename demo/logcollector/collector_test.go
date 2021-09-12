package logcollector

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logIdCtxKey{}, "2010")
	//output := os.Stdout
	output, err := os.OpenFile("run.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	logCollector := NewCollector(output)
	logCollector.Start()
	fmt.Println(1)

	go func() {
		count := 0
		for range time.Tick(time.Second) {
			logCollector.CtxLog(ctx, LogLevelError, "count:%d\n", count)
			count++
		}
	}()
	time.Sleep(time.Second * 100)
}
