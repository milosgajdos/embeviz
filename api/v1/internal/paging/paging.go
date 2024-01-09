package paging

import "reflect"

// ApplyOffsetLimit applies offset and limit to items and returns the result.
func ApplyOffsetLimit(items interface{}, offset int, limit int) interface{} {
	val := reflect.ValueOf(items)
	o := ApplyOffset(val, offset)
	if reflect.ValueOf(o).Len() == 0 {
		// NOTE: we could return Zero value here,
		// reflect.Zero(reflect.TypeOf(val.Interface())).Interface()
		// but we prefer to return an empty slice.
		return reflect.MakeSlice(val.Type(), 0, 0).Interface()
	}
	return ApplyLimit(reflect.ValueOf(o), limit)
}

// ApplyOffset returns a slice of items skipped by offset.
// If offset is negative, it returns the original slice.
// If offset is bigger than the number of items it returns empty slice.
func ApplyOffset(items reflect.Value, offset int) interface{} {
	if offset > 0 {
		switch {
		case items.Len() >= offset:
			return items.Slice(offset, items.Len()).Interface()
		default:
			// NOTE: we could return Zero value of items here,
			// return reflect.Zero(reflect.TypeOf(items.Interface())).Interface()
			// but we prefer to return an empty slice.
			return reflect.MakeSlice(reflect.TypeOf(items.Interface()), 0, 0).Interface()
		}
	}
	return items.Interface()
}

// ApplyLimit returns limit number of items.
// If limit is either negative or bigger than
// the number of items it returns all items.
func ApplyLimit(items reflect.Value, limit int) interface{} {
	if limit > 0 {
		switch {
		case items.Len() >= limit:
			return items.Slice(0, limit).Interface()
		default:
			return items.Interface()
		}
	}
	return items.Interface()
}
