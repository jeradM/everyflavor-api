package model

import "time"

type Tag struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Tag       string    `db:"tag"`
}

func (t Tag) GetID() uint64 {
	return t.ID
}

func (t Tag) TableName() string {
	return "tags"
}

func (t Tag) SelectFields() []string {
	return []string{
		"tags.id",
		"tags.created_at",
		"tags.updated_at",
		"tags.tag",
	}
}

func (t Tag) InsertMap() map[string]interface{} {
	return map[string]interface{}{"tag": t.Tag}
}

func (t Tag) UpdateMap() map[string]interface{} {
	return t.InsertMap()
}

func (t Tag) SelectJoins() []string {
	return []string{}
}
