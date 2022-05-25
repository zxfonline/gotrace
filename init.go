package trace

import (
	"net/http"

	"github.com/zxfonline/gotrace/expvar"
	"github.com/zxfonline/gotrace/golangtrace"
	"github.com/zxfonline/gotrace/pprof"
)

func Init(handler *http.ServeMux) {
	golangtrace.AuthRequest = func(req *http.Request) (any bool) {
		//TODO iptable init

		// RemoteAddr is commonly in the form "IP" or "IP:port".
		// If it is in the form "IP:port", split off the port.
		//host, _, err := net.SplitHostPort(req.RemoteAddr)
		//if err != nil {
		//	host = req.RemoteAddr
		//}
		//switch host {
		//case "localhost", "127.0.0.1", "::1":
		//	return true
		//default:
		//	return false
		//}
		return true
	}

	handler.HandleFunc("/debug/pprof/", func(w http.ResponseWriter, req *http.Request) {
		if !golangtrace.AuthRequest(req) {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		pprof.Index(w, req)
	})
	handler.HandleFunc("/debug/pprof/cmdline", func(w http.ResponseWriter, req *http.Request) {
		if !golangtrace.AuthRequest(req) {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		pprof.Cmdline(w, req)
	})
	handler.HandleFunc("/debug/pprof/profile", func(w http.ResponseWriter, req *http.Request) {
		if !golangtrace.AuthRequest(req) {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		pprof.Profile(w, req)
	})
	handler.HandleFunc("/debug/pprof/symbol", func(w http.ResponseWriter, req *http.Request) {
		if !golangtrace.AuthRequest(req) {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		pprof.Symbol(w, req)
	})
	handler.HandleFunc("/debug/pprof/trace", func(w http.ResponseWriter, req *http.Request) {
		if !golangtrace.AuthRequest(req) {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		pprof.Trace(w, req)
	})

	handler.HandleFunc("/debug/pprof/requests", func(w http.ResponseWriter, req *http.Request) {
		any := golangtrace.AuthRequest(req)
		if !any {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		golangtrace.Render(w, req, any)
	})

	handler.HandleFunc("/debug/pprof/events", func(w http.ResponseWriter, req *http.Request) {
		any := golangtrace.AuthRequest(req)
		if !any {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		golangtrace.RenderEvents(w, req, any)
	})
	//==============================

	handler.HandleFunc("/debug/pprof/vars", func(w http.ResponseWriter, req *http.Request) {
		if !golangtrace.AuthRequest(req) {
			http.Error(w, "not allowed", http.StatusUnauthorized)
			return
		}
		expvar.ExpvarHandler(w, req)
	})
	//==============================

	expvar.Publish("cmdline", expvar.Func(expvar.Cmdline))
	expvar.Publish("memstats", expvar.Func(expvar.Memstats))

	expvar.Publish("Goroutines", expvar.Func(goroutines))
	expvar.Publish("Uptime", expvar.Func(uptime))

	expvar.Publish("tracetotal", expvar.Func(traceTotal))
	expvar.Publish("tracehour", expvar.Func(traceHour))
	expvar.Publish("traceminute", expvar.Func(traceMinute))

	expvar.Publish("chanstats", expvar.Func(chanStats))
	//==============================
}
