package client

// Interface for odata Entity type
type IEntity interface {
	Key__() interface{}
	Data__() interface{}
	SetClient__(*Client)
}
