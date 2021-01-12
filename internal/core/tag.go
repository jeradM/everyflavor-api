package core

import (
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/mapper"
	"github.com/pkg/errors"
)

func (a *App) ListTags() ([]view.Tag, error) {
	t, err := a.Store.Tag().List()
	if err != nil {
		return nil, errors.Wrap(err, "App.ListTags() failed")
	}
	tags := make([]view.Tag, len(t))
	for idx, tag := range t {
		tags[idx] = mapper.TagFromModel(tag)
	}
	return tags, nil
}
