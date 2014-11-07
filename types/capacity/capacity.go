// Package Capacity implements the Capacity type. See:
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Capacity.html
package capacity

type ConsumedCapacityUnit float32

type ConsumedCapacityUnit_struct struct {
	CapacityUnits ConsumedCapacityUnit
}

type ConsumedCapacity struct {
	CapacityUnits          ConsumedCapacityUnit                   `json:",omitempty"`
	GlobalSecondaryIndexes map[string]ConsumedCapacityUnit_struct `json:",omitempty"`
	LocalSecondaryIndexes  map[string]ConsumedCapacityUnit_struct `json:",omitempty"`
	Table                  *ConsumedCapacityUnit_struct           `json:",omitempty"`
	TableName              string                                 `json:",omitempty"`
}

func NewConsumedCapacity() *ConsumedCapacity {
	c := new(ConsumedCapacity)
	c.GlobalSecondaryIndexes = make(map[string]ConsumedCapacityUnit_struct)
	c.LocalSecondaryIndexes = make(map[string]ConsumedCapacityUnit_struct)
	c.Table = new(ConsumedCapacityUnit_struct)
	return c
}

type ReturnConsumedCapacity string
