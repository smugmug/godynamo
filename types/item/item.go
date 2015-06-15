// Item as described in docs for various endpoints.
package item

import (
	"errors"
	"fmt"
	"github.com/smugmug/godynamo/types/attributevalue"
)

type Item attributevalue.AttributeValueMap

type item Item

// Item is already a reference type
func NewItem() Item {
	a := attributevalue.NewAttributeValueMap()
	return Item(a)
}

// Copy an Item
func (i Item) Copy(ic Item) error {
	if i == nil {
		return errors.New("Item.Copy: pointer receiver is nil")
	}
	if ic == nil {
		return errors.New("Item.Copy: copy target Item instance is nil")
	}

	for k, av := range i {
		ac := attributevalue.NewAttributeValue()
		if ac == nil {
			return errors.New("Item.Copy: copy target attributeValue is nil")
		}
		cp_err := av.Copy(ac)
		if cp_err != nil {
			e := fmt.Sprintf("Item.Copy:%s", cp_err.Error())
			return errors.New(e)
		}
		ic[k] = ac
	}
	return nil
}

// ItemLike is an interface for those structs you wish to map back and forth to Items.
// This is currently provided instead of the lossy translation advocated by the
// JSON document mapping as described by AWS.
type ItemLike interface {
	ToItem(interface{}) (Item, error)
	FromItem(Item) (interface{}, error)
}

// GetItem and UpdateItem share a Key type which is another alias to AttributeValueMap
type Key attributevalue.AttributeValueMap

func NewKey() Key {
	a := attributevalue.NewAttributeValueMap()
	return Key(a)
}
