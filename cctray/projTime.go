package cctray

import (
	"encoding/xml"
	"time"
)

type ProjTime struct {
	time.Time
}

func (pt *ProjTime) UnmarshalXMLAttr(attr xml.Attr) error {
	if parse, e := time.Parse("2006-01-02T15:04:05", attr.Value); e == nil {
		*pt = ProjTime{parse}
	} else {
		return e
	}
	return nil
}
