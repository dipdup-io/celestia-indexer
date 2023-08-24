package handler

import (
	"strings"

	"github.com/dipdup-net/indexer-sdk/pkg/storage"
)

const (
	asc  = "asc"
	desc = "desc"
)

type limitOffsetPagination struct {
	Limit  uint64 `json:"limit"  param:"limit"  query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset uint64 `json:"offset" param:"offset" query:"offset" validate:"omitempty,min=0"`
	Sort   string `json:"sort"   param:"sort"   query:"sort"   validate:"omitempty,oneof=asc desc"`
}

func (p *limitOffsetPagination) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

func pgSort(sort string) storage.SortOrder {
	switch sort {
	case asc:
		return storage.SortOrderAsc
	case desc:
		return storage.SortOrderDesc
	default:
		return storage.SortOrderAsc
	}
}

type txListRequest struct {
	Limit   uint64      `json:"limit"     param:"limit"                      query:"limit"  validate:"omitempty,min=1,max=100"`
	Offset  uint64      `json:"offset"    param:"offset"                     query:"offset" validate:"omitempty,min=0"`
	Sort    string      `json:"sort"      param:"sort"                       query:"sort"   validate:"omitempty,oneof=asc desc"`
	Status  StringArray `query:"status"   validate:"omitempty,dive,status"`
	MsgType StringArray `query:"msg_type" validate:"omitempty,dive,msg_type"`
}

func (p *txListRequest) SetDefault() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Sort == "" {
		p.Sort = asc
	}
}

type StringArray []string

func (s *StringArray) UnmarshalParam(param string) error {
	*s = StringArray(strings.Split(param, ","))
	return nil
}

type StatusArray StringArray
type MsgTypeArray StringArray
