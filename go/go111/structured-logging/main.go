package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		opeID := strconv.FormatInt(time.Now().UnixNano(), 10)

		info(logEntry{
			traceID: getAppEngineTraceID(r),
			message: "example!!",
			operation: logEntryOperation{
				id:       opeID,
				isFirst:  true,
				producer: "github.com/ryutah/appengine-sample/go/go111/structured-logging/example-first",
			},
			payload: logPayload{
				"foo": "bar",
			},
		})
		info(logEntry{
			traceID: getAppEngineTraceID(r),
			message: "example2!!",
			operation: logEntryOperation{
				id:       opeID,
				producer: "github.com/ryutah/appengine-sample/go/go111/structured-logging/example-second",
			},
			payload: logPayload{
				"foo": "bar",
			},
		})
		info(logEntry{
			traceID: getAppEngineTraceID(r),
			message: "example3!!",
			operation: logEntryOperation{
				id:       opeID,
				isLast:   true,
				producer: "github.com/ryutah/appengine-sample/go/go111/structured-logging/example-third",
			},
			payload: logPayload{
				"foo": "bar",
			},
		})

		w.Write([]byte("sucecss!!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

type (
	logEntry struct {
		traceID   string
		message   string
		payload   logPayload
		operation logEntryOperation
	}

	logEntryOperation struct {
		id       string
		isFirst  bool
		isLast   bool
		producer string
	}

	logPayload map[string]interface{}
)

func info(entry logEntry) {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	logger := log.New(os.Stdout, "", 0)
	payload := map[string]interface{}{
		"message":                      entry.message,
		"severity":                     "INFO",
		"logging.googleapis.com/trace": fmt.Sprintf("projects/%s/traces/%s", projectID, entry.traceID),
		"logging.googleapis.com/sourceLocation": map[string]interface{}{
			"file":     file,
			"line":     line,
			"function": f.Name(),
		},
		// optional, the id could be used to filtering log.
		"logging.googleapis.com/operation": map[string]interface{}{
			"id":       entry.operation.id,
			"producer": entry.operation.producer,
			"first":    entry.operation.isFirst,
			"last":     entry.operation.isLast,
		},
	}
	for k, v := range entry.payload {
		payload[k] = v
	}
	msg, _ := json.Marshal(payload)
	logger.Println(string(msg))
}

func getAppEngineTraceID(r *http.Request) string {
	val := r.Header.Get("X-Cloud-Trace-Context")
	return strings.SplitN(val, "/", 2)[0]
}
