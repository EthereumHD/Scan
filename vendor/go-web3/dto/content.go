package dto

import (
	"encoding/json"
)

type Content struct {
	Pending map[string]map[string]TransactionResponse
	Queued  map[string]map[string]TransactionResponse
}

func (c *Content) UnmarshalJSON(data []byte) error {
	type Alias Content
	temp := &struct {
		Pending map[string]map[string]TransactionResponse `json:"pending"`
		Queued  map[string]map[string]TransactionResponse `json:"queued"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	//fmt.Println("data:", string(data), "\ntemp", temp)
	c.Pending = temp.Pending
	c.Queued = temp.Queued

	return nil
}
