package model

import "fmt"

type ListParams interface {
	GetLimit() uint64
	GetOffset() uint64
	GetSort() string
	GetGroup() string
}

type BaseListParams struct {
	Limit  uint64 `json:"limit" form:"limit"`
	Offset uint64 `json:"offset" form:"offset"`
	Sort   string `json:"sort" form:"sort"`
	Order  string `json:"order" form:"order"`
	Group  string
}

func (b BaseListParams) GetLimit() uint64 {
	return b.Limit
}

func (b BaseListParams) GetOffset() uint64 {
	return b.Offset
}

func (b BaseListParams) GetSort() string {
	if b.Sort == "" {
		return ""
	}
	order := b.Order
	if order != "desc" {
		order = "asc"
	}
	return fmt.Sprintf("%s %s", b.Sort, order)
}

func (b BaseListParams) GetGroup() string {
	return b.Group
}
