package model

import "time"

type OutboundOrder struct {
	ID       uint32    `json:"id"`
	Creator  uint32    `json:"creator"`
	Status   uint8     `json:"status"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}
