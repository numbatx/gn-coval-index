package covalent

import "errors"

// ErrNilPubKeyConverter signals that a pub key converter input parameter is nil
var ErrNilPubKeyConverter = errors.New("received nil input value: pub key converter")

// ErrNilAccountsAdapter signals that an accounts adapter input parameter is nil
var ErrNilAccountsAdapter = errors.New("received nil input value: accounts adapter")

// ErrBlockBodyAssertion signals that an error occurred when trying to assert BodyHandler interface of type block body
var ErrBlockBodyAssertion = errors.New("error asserting BodyHandler interface of type block body")

// ErrNilHasher signals that a nil hasher has been provided
var ErrNilHasher = errors.New("received nil input value: hasher")

// ErrNilMarshaller signals that a nil marshaller has been provided
var ErrNilMarshaller = errors.New("received nil input value: marshaller")

// ErrNilMiniBlockHandler signals that a nil mini block handler has been provided
var ErrNilMiniBlockHandler = errors.New("received nil input value: mini block handler")

// ErrNilShardCoordinator signals that a shard coordinator input parameter is nil
var ErrNilShardCoordinator = errors.New("received nil input value: shard coordinator")

// ErrCannotCastAccountHandlerToUserAccount signals an error when trying to cast from AccountHandler to UserAccountHandler
var ErrCannotCastAccountHandlerToUserAccount = errors.New("cannot cast AccountHandler to UserAccountHandler")

// ErrNilDataHandler signals that a nil data handler handler has been provided
var ErrNilDataHandler = errors.New("received nil input value: data handler")

// ErrNilHTTPServer signals that a nil http server has been provided
var ErrNilHTTPServer = errors.New("received nil input value: http server")
