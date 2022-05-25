package trace

import (
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/zxfonline/gotrace/golangtrace"
)

var (
	startTime      = time.Now().UTC()
	monitorChanMap = make(map[string]interface{}, 16)
	monitorLock    sync.RWMutex
)

func goroutines() interface{} {
	return runtime.NumGoroutine()
}

func uptime() interface{} {
	uptime := time.Since(startTime)
	return int64(uptime)
}

func traceTotal() interface{} {
	return golangtrace.GetAllExpvarFamily(2)
}
func traceHour() interface{} {
	return golangtrace.GetAllExpvarFamily(1)
}
func traceMinute() interface{} {
	return golangtrace.GetAllExpvarFamily(0)
}

func isChan(a interface{}) bool {
	if a == nil {
		return false
	}

	v := reflect.ValueOf(a)
	if v.Kind() != reflect.Chan {
		return false
	}
	if v.IsNil() {
		return false
	}
	return true
}

func chanInfo(a interface{}) (bool, int, int) {
	if a == nil {
		return false, 0, 0
	}

	v := reflect.ValueOf(a)
	if v.Kind() != reflect.Chan {
		return false, 0, 0
	}
	if v.IsNil() {
		return false, 0, 0
	}
	return true, v.Cap(), v.Len()
}

type ChanInfo struct {
	Cap  int
	Len  int
	Rate float64
}

func chanStats() interface{} {
	monitorLock.RLock()
	defer monitorLock.RUnlock()
	mp := make(map[string]ChanInfo)
	for k, v := range monitorChanMap {
		_, m, i := chanInfo(v)
		mp[k] = ChanInfo{Cap: m, Len: i, Rate: float64(int64((float64(i) / float64(m) * 10000.0))) / 10000.0}
	}
	return mp
}
