package models

import (
	"encoding/json"
	"fmt"
	"pubsub-ckg-tb/internal/config"
)

const (
	PUBSUB_CONSUME = iota
	PUBSUB_PRODUCE
)

// PubSubObject is the base class for CKG data objects
type PubSubObjectWrapper[T PubSubObject] struct {
	CKGObject bool
	Type      int
	Data      []T `json:"data"`
}

type PubSubObject interface {
	FromMap(data map[string]interface{})
	ToMap() map[string]interface{}
}

func NewPubSubConsumerWrapper[T PubSubObject]() PubSubObjectWrapper[T] {
	return PubSubObjectWrapper[T]{
		Type: PUBSUB_CONSUME,
	}
}

func NewPubSubProducerWrapper[T PubSubObject](data []T) PubSubObjectWrapper[T] {
	return PubSubObjectWrapper[T]{
		Type: PUBSUB_PRODUCE,
		Data: data,
	}
}

// FromMap creates a PubSubObject from a map
func (t *PubSubObjectWrapper[T]) FromMap(obj map[string]interface{}) *PubSubObjectWrapper[T] {
	cfg := config.GetConfig()

	markerValueObj, ok := obj[cfg.CKG.MarkerField].(string)
	if !ok || markerValueObj == "" {
		return t
	}

	markerValueStruct := ""
	switch t.Type {
	case PUBSUB_PRODUCE:
		markerValueStruct = cfg.CKG.MarkerProduce
	case PUBSUB_CONSUME:
		markerValueStruct = cfg.CKG.MarkerConsume
	}

	// Check if this is a CKG object
	if markerValueObj == markerValueStruct {
		t.CKGObject = true
		t.Data = make([]T, 0)

		if data, ok := obj["data"].([]interface{}); ok {
			for _, item := range data {
				var x T
				if dataMap, ok := item.(map[string]interface{}); ok {
					x.FromMap(dataMap)
					t.Data = append(t.Data, x)
				}
			}
		}
	}

	return t
}

func (t *PubSubObjectWrapper[T]) FromJSON(jsonStr string) error {
	data := make(map[string]interface{})

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	t.FromMap(data)
	return nil
}

func (t *PubSubObjectWrapper[T]) ToMap() map[string]interface{} {
	data := make(map[string]interface{})

	items := make([]interface{}, 0)
	for _, item := range t.Data {
		items = append(items, item.ToMap())
	}

	data["data"] = items
	return data
}

// ToJSON converts PubSubObject to JSON string
func (t *PubSubObjectWrapper[T]) ToJSON() (string, error) {
	data := t.ToMap()
	cfg := config.GetConfig()

	// Add marker field
	if t.Type == PUBSUB_PRODUCE {
		data[cfg.CKG.MarkerField] = cfg.CKG.MarkerProduce
	} else {
		data[cfg.CKG.MarkerField] = cfg.CKG.MarkerConsume
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %v", err)
	}

	return string(jsonBytes), nil
}

func (t *PubSubObjectWrapper[T]) IsCKGObject() bool {
	return t.CKGObject
}
