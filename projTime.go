package main

import (
	"encoding/xml"
	"time"
)

type projTime struct {
	time.Time
}

func (pt *projTime) UnmarshalXMLAttr(attr xml.Attr) error {
	if parse, e := time.Parse("2006-01-02T15:04:05", attr.Value); e == nil {
		*pt = projTime{parse}
	} else {
		return e
	}
	return nil
}
