package net

import (
	"fmt"
	"github.com/pkg/errors"
)

const (
	// Unknown error code
	ErrUnknown = 0
	// HTTP method not allowed
	ErrMethodNotAllowed = 1
	// No authorization header was found in request
	ErrNoAuthHdr = 2
	// Failed reading the HTTP request body
	ErrReadingReqBody = 3
	// HTTP request body was empty but should not have been
	ErrEmptyReqBody = 4
	// HTTP request body was of the wrong format
	ErrBadReqBody = 5
	// The token in the authorization header was invalid
	ErrBadReqToken = 6
	// The caller does not have the "hf.Revoker" attibute
	ErrNotRevoker = 7
	// Certificate to be revoked was not found
	ErrRevCertNotFound = 8
	// Certificate to be revoked is not owned by expected user
	ErrCertWrongOwner = 9
	// Identity of certificate to be revoked was not found
	ErrRevokeIDNotFound = 10
	// User info was not found for issuee of revoked certificate
	ErrRevokeUserInfoNotFound = 11
	// Certificate revocation failed for another reason
	ErrRevokeFailure = 12
	// Failed to update user info when revoking identity
	ErrRevokeUpdateUser = 13
	// Failed to revoke any certificates by identity
	ErrNoCertsRevoked = 14
	// Missing fields in the revocation request
	ErrMissingRevokeArgs = 15
	// Failed to get user's affiliation
	ErrGettingAffiliation = 16
	// Revoker's affiliation not equal to or above revokee's affiliation
	ErrRevokerNotAffiliated = 17
	// Failed to send an HTTP response
	ErrSendingResponse = 18
	// The CA (Certificate Authority) name was not found
	ErrCANotFound = 19
	// Authorization failure
	ErrAuthFailure = 20
	// No username and password were in the authorization header
	ErrNoUserPass = 21
	// Enrollment is currently disabled for the server
	ErrEnrollDisabled = 22
	// Invalid user name
	ErrInvalidUser = 23
	// Invalid password
	ErrInvalidPass = 24
	// Invalid token in authorization header
	ErrInvalidToken = 25
	// Certificate was not issued by a trusted authority
	ErrUntrustedCertificate = 26
	// Certificate has expired
	ErrCertExpired = 27
	// Certificate has been revoked
	ErrCertRevoked = 28
	// Failed trying to check if certificate is revoked
	ErrCertRevokeCheckFailure = 29
	// Certificate was not found
	ErrCertNotFound = 30
	// Bad certificate signing request
	ErrBadCSR = 31
	// Failed to get identity's prekey
	ErrNoPreKey = 32
	// The caller was not authenticated
	ErrCallerIsNotAuthenticated = 33
	// Invalid configuration setting
	ErrConfig = 34
	// The caller does not have authority to generate a CRL
	ErrNoGenCRLAuth = 35
	// Invalid RevokedAfter value in the GenCRL request
	ErrInvalidRevokedAfter = 36
	// Invalid ExpiredAfter value in the GenCRL request
	ErrInvalidExpiredAfter = 37
	// Failed to get revoked certs from the database
	ErrRevokedCertsFromDB = 38
	// Failed to get CA cert
	ErrGetCACert = 39
	// Failed to get CA signer
	ErrGetCASigner = 40
	// Failed to generate CRL
	ErrGenCRL = 41
	// Registrar does not have the authority to register an attribute
	ErrRegAttrAuth = 42
	// Registrar does not own 'hf.Registrar.Attributes'
	ErrMissingRegAttr = 43
	// Caller does not have appropriate affiliation to perform requested action
	ErrCallerNotAffiliated = 44
	// Failed to verify if caller has appropriate type
	ErrGettingType = 45
	// CA cert does not have 'crl sign' usage
	ErrNoCrlSignAuth = 46
	// Incorrect level of database
	ErrDBLevel = 47
	// Incorrect level of configuration file
	ErrConfigFileLevel = 48
	// Failed to get user from database
	ErrGettingUser = 49
	// Error processing HTTP request
	ErrHTTPRequest = 50
	// Error connecting to database
	ErrConnectingDB = 51
	// Failed to add identity
	ErrAddIdentity = 52
	// Unauthorized to perform update action
	ErrUpdateConfigAuth = 53
	// Registrar not authorized to act on type
	ErrRegistrarInvalidType = 54
	// Registrar not authorized to act on affiliation
	ErrRegistrarNotAffiliated = 55
	// Failed to remove identity
	ErrRemoveIdentity = 56
	// Failed to get boolean query parameter
	ErrGettingBoolQueryParm = 57
	// Failed to modify identity
	ErrModifyingIdentity = 58
	// Caller does not have the appropriate role
	ErrMissingRole = 59
	// Failed to add new affiliation
	ErrUpdateConfigAddAff = 60
	// Failed to remove affiliation
	ErrUpdateConfigRemoveAff = 61
	// Error occured while removing affiliation in database
	ErrRemoveAffDB = 62
	// Error occured when making a Get request to database
	ErrDBGet = 63
	// Failed to modiy affiliation
	ErrUpdateConfigModifyAff = 64
	// Error occured while deleting user
	ErrDBDeleteUser = 65
	// Certificate that is being revoked has already been revoked
	ErrCertAlreadyRevoked = 66
	// Failed to get requested certificate(s)
	ErrGettingCert = 67
	// Error occurred parsing variable as an integer
	ErrParsingIntEnvVar = 68
	// CA certificate file is not found warning message
	ErrCACertFileNotFound = 69
)

type httpErr struct {
	scode int    // HTTP status code
	lcode int    // local error code
	lmsg  string // local error message
	rcode int    // remote error code
	rmsg  string // remote error message
}

func createHTTPErr(scode, code int, format string, args ...interface{}) *httpErr {
	msg := fmt.Sprintf(format, args...)
	return &httpErr{
		scode: scode,
		lcode: code,
		lmsg:  msg,
		rcode: code,
		rmsg:  msg,
	}
}

func newHTTPErr(scode, code int, format string, args ...interface{}) error {
	return errors.Wrap(createHTTPErr(scode, code, format, args...), "")
}

func (he *httpErr) Error() string {
	return he.String()
}

func (he *httpErr) String() string {
	if he.lcode == he.rcode && he.lmsg == he.rmsg {
		return fmt.Sprintf("scode: %d, code: %d, msg: %s", he.scode, he.lcode, he.lmsg)
	}
	return fmt.Sprintf("scode: %d, local code: %d, local msg: %s, remote code: %d, remote msg: %s",
		he.scode, he.lcode, he.lmsg, he.rcode, he.rmsg)
}
