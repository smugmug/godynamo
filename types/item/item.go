// Item as described in docs for various endpoints.
package item

// Note:
// In the DynamoDB feature announcement here:
// http://aws.amazon.com/blogs/aws/dynamodb-update-json-and-more/
// They describe an example of creating an Item directly from a JSON document.
// As should be obvious from reviewing the documentation for the AttributeValue type,
// DynamoDB does not actually have a JSON document core type. The SDK methods for
// creating Items from JSON appear to just be a convenience that makes assumptions about
// how data should be mapped to existing AttributeValue fields (the addition
// of "L" and "M" fields merely make it *possible* to map arbitrary JSON).
// This SDK feature seems to only support trivial conversions that do not encompass
// the full breadth of the AttributeValue type. For example, it isn't clear how
// the SDK differentiates between B/BS and S/SS types (both are stringy).
//
// I recommend just writing ToItem() and FromItem() methods where you need them
// if you wish to have this functionality in your own programs to map from
// your own types to and from Items. The interface ItemLike below can be used
// to implement this.

import (
	"github.com/smugmug/godynamo/types/attributevalue"
)

type Item attributevalue.AttributeValueMap

type item Item

// Item is already a reference type
func NewItem() Item {
	a := attributevalue.NewAttributeValueMap()
	return Item(a)
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
