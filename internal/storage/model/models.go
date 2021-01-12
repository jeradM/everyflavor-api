package model

type Table interface {
	TableName() string
}

type Entity interface {
	Table
	GetID() uint64
	SelectFields() []string
	InsertMap() map[string]interface{}
	UpdateMap() map[string]interface{}
	SelectJoins() []string
}

type OwnedEntity interface {
	Entity
	OwnerField() string
}

type PublishableEntity interface {
	Entity
	PublicField() string
}

type SharedEntity interface {
	OwnedEntity
	CollaboratorIDsQuery(uint64) (string, []interface{})
}

type ListResult struct {
	Result interface{}
	Count  uint64
}
