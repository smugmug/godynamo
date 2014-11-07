// Condition implements conditions used in Query and Scan. See:
// http://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_Condition.html
package condition

import (
	"github.com/smugmug/godynamo/types/attributevalue"
)

type Conditions map[string]*Condition

type Condition struct {
	AttributeValueList []*attributevalue.AttributeValue
	ComparisonOperator string
}

func NewConditions() Conditions {
	k := make(map[string]*Condition)
	return k
}

func NewCondition() *Condition {
	k := new(Condition)
	k.AttributeValueList = make([]*attributevalue.AttributeValue, 0)
	return k
}
