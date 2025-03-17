package controllers

import (
	"fmt"
	"github.com/asdine/storm/q"
)

type Condition struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func explainQueryCondition(conditions []*Condition) []q.Matcher {
	var matcher []q.Matcher
	for _, v := range conditions {
		var match q.Matcher
		switch v.Type {
		case "eq":
			match = q.Eq(v.Name, v.Value)
		case "re":
			match = q.Re(v.Name, fmt.Sprintf("^%s", v.Value))
		case "lt":
			match = q.Lt(v.Name, v.Value)
		case "lte":
			match = q.Lte(v.Name, v.Value)
		case "gt":
			match = q.Gt(v.Name, v.Value)
		case "gte":
			match = q.Gte(v.Name, v.Value)
		case "in":
			match = q.In(v.Name, v.Value)
		}
		matcher = append(matcher, match)
	}

	return matcher
}
