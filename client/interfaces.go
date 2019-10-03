package client

// Interface for odata Entity type
type IEntity interface {
	Key__() interface{}
	Data__() interface{}
}

// Interface for odata PrimaryKey type
type IPrimaryKey interface {
	APIEntityType() string
	Serialize() string
}

// Interface for odata Function type
type IFunction interface {
	Name() string
	Parameters() string
}
