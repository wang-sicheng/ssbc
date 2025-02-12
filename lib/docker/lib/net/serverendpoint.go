package net

import (
	"encoding/json"
	"github.com/cloudflare/cfssl/api"
	"github.com/cloudflare/cfssl/log"
	"github.com/ssbc/lib/caerrors"
	"net/http"
)

type serverEndpoint struct {
	// The HTTP methods ("GET", "POST", etc) which the function will handle
	Methods []string
	// The HTTP status code for a successful response
	successRC int
	// Handler is the handler function for this endpoint
	Handler func(ctx *serverRequestContextImpl) (interface{}, error)
	// Server which hosts this endpoint
	Server *Server
}

func (se *serverEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp interface{}
	url := r.URL.String()
	log.Debugf("Received request for %s", url)
	w = newHTTPResponseWriter(r, w, se)
	err := se.validateMethod(r)
	if err == nil {
		// Call the endpoint handler to handle the request.  The handler may
		// a) return the response in the 'resp' variable below, or
		// b) write the response one chunk at a time, which is appropriate if the response may be large
		//    and we don't want the server to buffer the entire response in memory.
		resp, err = se.Handler(newServerRequestContext(r, w, se))
	}
	he := getHTTPErr(err)
	if he != nil {
		// An error occurred
		w.WriteHeader(he.GetStatusCode())
		log.Infof(`%s %s %s %d %d "%s"`, r.RemoteAddr, r.Method, r.URL, he.GetStatusCode(), he.GetLocalCode(), he.GetLocalMsg())
	} else {
		// No error occurred
		scode := se.getSuccessRC()
		w.WriteHeader(scode)
		log.Infof(`%s %s %s %d 0 "OK"`, r.RemoteAddr, r.Method, r.URL, scode)
	}
	// If a response was returned by the handler, write it now.
	if resp != nil {
		writeJSON(resp, w)
	}
	// If nothing has been written, write an empty string for the response
	if !w.(*httpResponseWriter).writeCalled {
		w.Write([]byte(`""`))
	}
	// If an error was returned by the handler, write it now.
	w.Write([]byte(`,"errors":[`))
	if he != nil {
		rm := &api.ResponseMessage{Code: he.GetRemoteCode(), Message: he.GetRemoteMsg()}
		writeJSON(rm, w)
	}
	// Write true or false for success
	w.Write([]byte(`],"messages":[],"success":`))
	if he != nil {
		w.Write([]byte(`false}`))
	} else {
		w.Write([]byte(`true}`))
	}
	w.(http.Flusher).Flush()
}

func (se *serverEndpoint) getSuccessRC() int {
	if se.successRC == 0 {
		return 200
	}
	return se.successRC
}

func writeJSON(obj interface{}, w http.ResponseWriter) {
	enc := json.NewEncoder(w)
	err := enc.Encode(obj)
	if err != nil {
		log.Errorf("Failed encoding response to JSON: %s", err)
	}
}

func newHTTPResponseWriter(r *http.Request, w http.ResponseWriter, se *serverEndpoint) *httpResponseWriter {
	return &httpResponseWriter{r: r, w: w, se: se}
}

type httpResponseWriter struct {
	r                 *http.Request
	w                 http.ResponseWriter
	se                *serverEndpoint
	writeHeaderCalled bool
	writeCalled       bool
}

func (hrw *httpResponseWriter) Header() http.Header {
	return hrw.w.Header()
}

func (hrw *httpResponseWriter) WriteHeader(scode int) {
	if !hrw.writeHeaderCalled {
		w := hrw.w
		w.Header().Set("Connection", "Keep-Alive")
		if hrw.isHead() {
			w.Header().Set("Content-Length", "0")
		} else {
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("Content-Type", "application/json")
		}
		// Write the appropriate successful status code for this endpoint
		if scode == http.StatusOK {
			scode = hrw.se.getSuccessRC()
		}
		w.WriteHeader(scode)
		hrw.writeHeaderCalled = true
	}
}

// Write writes the data to the connection as part of an HTTP reply.
func (hrw *httpResponseWriter) Write(buf []byte) (int, error) {
	if hrw.isHead() {
		return 0, nil
	}
	w := hrw.w
	hrw.WriteHeader(http.StatusOK)
	if !hrw.writeCalled {
		// Write the header of the body of the result
		b, err := w.Write([]byte(`{"result":`))
		if err != nil {
			return b, err
		}
		hrw.writeCalled = true
	}
	if buf == nil {
		return 0, nil
	}
	return w.Write(buf)
}

func (hrw *httpResponseWriter) Flush() {
	hrw.w.(http.Flusher).Flush()
}

func (hrw *httpResponseWriter) isHead() bool {
	return hrw.r.Method == "HEAD"
}

func (se *serverEndpoint) validateMethod(r *http.Request) error {
	for _, m := range se.Methods {
		if m == r.Method {
			return nil
		}
	}
	return caerrors.NewHTTPErr(405, caerrors.ErrMethodNotAllowed, "Method %s is not allowed", r.Method)
}

func getHTTPErr(err error) *caerrors.HTTPErr {
	if err == nil {
		return nil
	}
	type causer interface {
		Cause() error
	}
	curErr := err
	for curErr != nil {
		switch curErr.(type) {
		case *caerrors.HTTPErr:
			return curErr.(*caerrors.HTTPErr)
		case causer:
			curErr = curErr.(causer).Cause()
		default:
			return caerrors.CreateHTTPErr(500, caerrors.ErrUnknown, err.Error())
		}
	}
	return caerrors.CreateHTTPErr(500, caerrors.ErrUnknown, "nil error")
}
