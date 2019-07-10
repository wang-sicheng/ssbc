package net

import (
	"crypto/x509"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type serverRequestContextImpl struct {
	req            *http.Request
	resp           http.ResponseWriter
	endpoint       *serverEndpoint

	enrollmentID   string
	enrollmentCert *x509.Certificate

	body           struct {
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


//func (ctx *serverRequestContextImpl) BasicAuthentication() (string, error) {
//	r := ctx.req
//	// Get the authorization header
//	authHdr := r.Header.Get("authorization")
//	if authHdr == "" {
//		return "", caerrors.NewHTTPErr(401, caerrors.ErrNoAuthHdr, "No authorization header")
//	}
//	// Extract the username and password from the header
//	username, password, ok := r.BasicAuth()
//	if !ok {
//		return "", caerrors.NewAuthenticationErr(caerrors.ErrNoUserPass, "No user/pass in authorization header")
//	}
//	// Get the CA that is targeted by this request
//	ca, err := ctx.GetCA()
//	if err != nil {
//		return "", err
//	}
//	// Error if max enrollments is disabled for this CA
//	log.Debugf("ca.Config: %+v", ca.Config)
//	caMaxEnrollments := ca.Config.Registry.MaxEnrollments
//	if caMaxEnrollments == 0 {
//		return "", caerrors.NewAuthenticationErr(caerrors.ErrEnrollDisabled, "Enroll is disabled")
//	}
//	// Get the user info object for this user
//	ctx.ui, err = ca.registry.GetUser(username, nil)
//	if err != nil {
//		return "", caerrors.NewAuthenticationErr(caerrors.ErrInvalidUser, "Failed to get user: %s", err)
//	}
//	// Check the user's password and max enrollments if supported by registry
//	err = ctx.ui.Login(password, caMaxEnrollments)
//	if err != nil {
//		return "", caerrors.NewAuthenticationErr(caerrors.ErrInvalidPass, "Login failure: %s", err)
//	}
//	// Store the enrollment ID associated with this server request context
//	ctx.enrollmentID = username
//	ctx.caller, err = ctx.GetCaller()
//	if err != nil {
//		return "", err
//	}
//	// Return the username
//	return username, nil
//}