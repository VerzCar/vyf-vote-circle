package model

type EventOperation string

const (
	EventOperationCreated EventOperation = "CREATED"
	EventOperationUpdated EventOperation = "UPDATED"
	EventOperationDeleted EventOperation = "DELETED"
)
