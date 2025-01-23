package litereq

import "time"

// request base

const (
	minProxyExpired = 60 * time.Second
	maxProxyExpired = 5 * time.Minute
	DefaultTimeout  = 5 * time.Second
)

// do request

type doResponse int

const (
	doOK doResponse = iota
	doConnect
	doValidate
	doHandle
)

// proxy

const (
	Normal ProxyStatus = iota
	Switch ProxyStatus = 2
)

const (
	Remove ProxyChange = iota
	Create
)

const (
	Short ProxyTimelyType = iota
	Long
)
