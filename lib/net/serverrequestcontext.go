package net

import (
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type serverRequestContextImpl struct {
	req      *http.Request
	resp     http.ResponseWriter
	endpoint *serverEndpoint

	enrollmentID   string
	enrollmentCert *x509.Certificate

	body struct {
		read bool   // true after body is read
		buf  []byte // the body itself
		err  error  // any error from reading the body
	}
	callerRoles map[string]bool
}

func newServerRequestContext(r *http.Request, w http.ResponseWriter, se *serverEndpoint) *serverRequestContextImpl {
	return &serverRequestContextImpl{
		req:      r,
		resp:     w,
		endpoint: se,
	}
}

func (ctx *serverRequestContextImpl) ReadBody(body interface{}) error {
	empty, err := ctx.TryReadBody(body)
	if err != nil {
		return err
	}
	if empty {
		return newHTTPErr(400, ErrEmptyReqBody, "Empty request body")
	}
	return nil
}

func (ctx *serverRequestContextImpl) TryReadBody(body interface{}) (bool, error) {
	buf, err := ctx.ReadBodyBytes()
	if err != nil {
		return false, err
	}
	empty := len(buf) == 0
	if !empty {
		err = json.Unmarshal(buf, body)
		if err != nil {
			return true, newHTTPErr(400, ErrBadReqBody, "Invalid request body: %s; body=%s",
				err, string(buf))
		}
	}
	return empty, nil
}

func (ctx *serverRequestContextImpl) ReadBodyBytes() ([]byte, error) {
	if !ctx.body.read {
		r := ctx.req
		buf, err := ioutil.ReadAll(r.Body)
		ctx.body.buf = buf
		ctx.body.err = err
		ctx.body.read = true
	}
	err := ctx.body.err
	if err != nil {
		return nil, newHTTPErr(400, ErrReadingReqBody, "Failed reading request body: %s", err)
	}
	return ctx.body.buf, nil
}
