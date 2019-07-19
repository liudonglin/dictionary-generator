package tpl

import (
	"dg-server/core"
	"strings"
)

// Enum ...
type Enum struct {
	Name  string
	Title string
	Kvds  []*KVD
}

// KVD ...
type KVD struct {
	Key   string
	Value string
	Des   string
}

func getEnums(cols []*core.Column) []*Enum {
	enums := make([]*Enum, 0)
	for _, col := range cols {
		enum := analysisEnumString(col.Enum)
		if enum != nil {
			enum.Name = fn.ToCamelString(col.Name)
			enum.Title = col.Title
			enums = append(enums, enum)
			col.IsEnum = true
		}
	}
	return enums
}

func analysisEnumString(s string) *Enum {
	if s == "" {
		return nil
	}
	items := strings.Split(s, ";")
	if len(items) == 0 {
		return nil
	}

	kvds := make([]*KVD, 0)
	for _, item := range items {
		kvd := strings.Split(item, ":")
		if len(kvd) == 3 {
			kvds = append(kvds, &KVD{Key: kvd[0], Value: kvd[1], Des: kvd[2]})
		}
	}

	if len(kvds) == 0 {
		return nil
	}
	return &Enum{Kvds: kvds}
}
