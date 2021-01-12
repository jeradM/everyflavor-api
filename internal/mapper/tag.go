package mapper

import (
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage/model"
)

func TagFromModel(m model.Tag) view.Tag {
	return view.Tag{ID: m.ID, Tag: m.Tag}
}

func TagToModel(v view.Tag) model.Tag {
	return model.Tag{ID: v.ID, Tag: v.Tag}
}
