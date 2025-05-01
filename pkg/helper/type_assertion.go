package helper

import (
	"fmt"
	"strconv"
)

func ChangeID(idValue interface{}) (int, error) {
	switch v := idValue.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		parsedID, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return parsedID, nil

	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}

}
