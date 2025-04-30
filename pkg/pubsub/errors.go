package pubsub

import "errors"

var (
	// Client Errors.
	ErrClientNotConnected  = errors.New("nats client not connected")
	ErrSubjectNotProvided  = errors.New("subject not provided")
	ErrSubjectsNotProvided = errors.New("subjects not provided")
	// errConsumerNotProvided = errors.New("consumer name not provided")
	ErrStreamNotProvided = errors.New("stream name not provided")
	// errConsumerCreationError = errors.New("consumer creation error")
	// errFailedToDeleteStream    = errors.New("failed to delete stream")
	// errPublishError            = errors.New("publish error")
	ErrJetStreamNotConfigured  = errors.New("jStream is not configured")
	ErrJetStreamCreationFailed = errors.New("jStream creation failed")
	ErrJetStream               = errors.New("jStream error")
	// errCreateStream = errors.New("create stream error")
	// errDeleteStream = errors.New("delete stream error")
	// errGetStream    = errors.New("get stream error")
	// errCreateOrUpdateStream = errors.New("create or update stream error")
	// errHandlerError            = errors.New("handler error")
	// errConnectionError = errors.New("connection error")
	ErrSubscriptionError = errors.New("subscription error")
)
