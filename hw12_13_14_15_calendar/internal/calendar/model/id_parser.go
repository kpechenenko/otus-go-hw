package model

import (
	"fmt"

	"github.com/google/uuid"
)

func ParseEventID(s string) (EventID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return EventID{}, fmt.Errorf("fail to parse eventID: %v", err)
	}
	return EventID(id), nil
}
