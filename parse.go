package rfc9651

import (
	"fmt"

	"go.jcbhmr.com/rfc9651/internal/ascii"
	"go.jcbhmr.com/rfc9651/internal/molite"
	"go.jcbhmr.com/rfc9651/internal/orderedmap"
)

type FieldType string

const (
	FieldTypeDictionary FieldType = "dictionary"
	FieldTypeList       FieldType = "list"
	FieldTypeItem       FieldType = "item"
)

func (ft FieldType) isValid() bool {
	switch ft {
	case FieldTypeDictionary, FieldTypeList, FieldTypeItem:
		return true
	default:
		return false
	}
}

type Dictionary = *orderedmap.OrderedMap[ascii.String, molite.Tuple[IterOrInnerList, Parameters]]

type List = []molite.Tuple[molite.Either[Item, InnerList], Parameters]

type InnerList = []molite.Tuple[Item, Parameters]

func Parse(inputBytes []byte, fieldType FieldType) (output molite.Either3[List, Dictionary, Item], err error) {
	if !fieldType.isValid() {
		panic(fmt.Errorf("invalid FieldType %q", fieldType))
	}

	inputString, err := ascii.NewString(inputBytes)
	if err != nil {
		return nil, err
	}

	inputString = inputString.TrimPrefix(" ")

	if fieldType == FieldTypeList {
		result, err := ParseList(&inputString)
		if err != nil {
			return nil, err
		}
		output = molite.Either3[List, Dictionary, Item]{Which: 0, A: result}
	}

	if fieldType == FieldTypeDictionary {
		result, err := ParseDictionary(&inputString)
		if err != nil {
			return nil, err
		}
		output = molite.Either3[List, Dictionary, Item]{Which: 1, B: result}
	}

	if fieldType == FieldTypeItem {
		result, err := ParseItem(&inputString)
		if err != nil {
			return nil, err
		}
		output = molite.Either3[List, Dictionary, Item]{Which: 2, C: result}
	}

	inputString = inputString.TrimPrefix(" ")

	if inputString.Len() > 0 {
		return nil, fmt.Errorf("input string %q is not empty after parsing %v", inputString, output)
	} else {
		return output, nil
	}
}

func ParseList(inputString *ascii.String) ([]molite.Tuple[molite.Either[Item, InnerList], Parameters], error) {
	members := []molite.Tuple[molite.Either[Item, InnerList], Parameters]{}

	for inputString.Len() > 0 {
		result, err := ParseItemOrInnerList(inputString)
		if err != nil {
			return nil, err
		}
		members = append(members, molite.Tuple[molite.Either[Item, InnerList], Parameters]{A: result, B: Parameters{}})

		*inputString = inputString.TrimLeft(" \t")

		if inputString.Len() == 0 {
			return members, nil
		}

		firstChar := inputString.Get(0)
		*inputString = inputString.Slice(1, inputString.Len())
		if firstChar != ',' {
			return nil, fmt.Errorf("expected ',' after list member, got %q", firstChar)
		}

		*inputString = inputString.TrimLeft(", \t")

		if inputString.Len() == 0 {
			return nil, fmt.Errorf("expected list member after ',', got empty string")
		}
	}

	return members, nil
}

func ParseItemOrInnerList(inputString *ascii.String) (molite.Tuple[molite.Either[Item, []molite.Tuple[Item, Parameters]], Parameters], error) {
	if inputString.Get(0) == '(' {
		result, err := ParseInnerList(inputString)
		if err != nil {
			return molite.Tuple[molite.Either[Item, []molite.Tuple[Item, Parameters]], Parameters]{}, err
		}
		return result, nil
	}

	result, err := ParseItem(inputString)
	if err != nil {
		return molite.Tuple[molite.Either[Item, []molite.Tuple[Item, Parameters]], Parameters]{}, err
	}
	return result, nil
}

func ParseInnerList(inputString *ascii.String) (molite.Tuple[[]molite.Tuple[Item, Parameters], Parameters], error) {
	firstChar := inputString.Get(0)
	*inputString = inputString.Slice(1, inputString.Len())
	if firstChar != '(' {
		return molite.Tuple[[]molite.Tuple[Item, Parameters], Parameters]{}, fmt.Errorf("expected '(' at start of inner list, got %q", firstChar)
	}

	innerList := []molite.Tuple[Item, Parameters]{}

	for inputString.Len() > 0 {
		*inputString = inputString.TrimLeft(" ")

		if inputString.Get(0) == ')' {
			*inputString = inputString.Slice(1, inputString.Len())

			parameters, err := ParseParameters(inputString)
			if err != nil {
				return molite.Tuple[[]molite.Tuple[Item, Parameters], Parameters]{}, err
			}

			return molite.Tuple[[]molite.Tuple[Item, Parameters], Parameters]{A: innerList, B: parameters}, nil
		}

		result, err := ParseItem(inputString)
		if err != nil {
			return molite.Tuple[[]molite.Tuple[Item, Parameters], Parameters]{}, err
		}

		innerList = append(innerList, molite.Tuple[Item, Parameters]{A: result, B: Parameters{}})

		if firstChar := inputString.Get(0); !(firstChar == ' ' || firstChar == ')') {
			return molite.Tuple[[]molite.Tuple[Item, Parameters], Parameters]{}, fmt.Errorf("expected ' ' or ')' after inner list item, got %q", firstChar)
		}
	}

	return molite.Tuple[[]molite.Tuple[Item, Parameters], Parameters]{}, fmt.Errorf("unexpected end of input while parsing inner list")
}

func ParseDictionary(inputString *ascii.String) (*orderedmap.OrderedMap[ascii.String, molite.Tuple[]])
