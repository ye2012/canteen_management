package dto

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	a := OrderNode{ID: "123", Name: "343", Price: 111}
	str, err := json.Marshal(a)
	if err != nil {
		fmt.Println("Marshal Failed|Err:", err)
		return
	}
	fmt.Println(string(str))
}
