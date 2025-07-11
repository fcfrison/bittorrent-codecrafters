package main

import "strconv"

type KeyValue struct {
	key   []byte
	value []byte
}

type KeyValueSlice []KeyValue

func (kvp *KeyValueSlice) AddSorted(kv KeyValue) {
	i := 0
	kvsl := *kvp
	for i < len(kvsl) {
		if string(kv.key) < string(kvsl[i].key) {
			break
		}
		i++
	}
	var newSlice KeyValueSlice = make(KeyValueSlice, 0)
	if i > 0 {
		newSlice = append(newSlice, kvsl[:i]...)
	}
	newSlice = append(newSlice, kv)
	if i < len(kvsl) {
		newSlice = append(newSlice, kvsl[i:]...)
	}
	*kvp = newSlice

}
func (kvp *KeyValueSlice) Pop() (KeyValue, error) {
	var kv KeyValue
	var kvsl KeyValueSlice = *kvp
	if len(*kvp) <= 0 {
		return kv, nil
	}
	kv = kvsl[0]
	if len(*kvp) == 1 {
		kvp = new(KeyValueSlice)
		return kv, nil
	}
	*kvp = kvsl[1:]
	return kv, nil

}
func EncodeString(str []byte) ([]byte, error) {
	strLength := strconv.Itoa(len(str)) + ":"
	return append([]byte(strLength), str...), nil

}
func EncodeInteger(value int) ([]byte, error) {
	valueStr := "i" + strconv.Itoa(value) + "e"
	return []byte(valueStr), nil
}
func EncodeDictionary(dict map[string]any) ([]byte, error) {
	kvsl := new(KeyValueSlice)
	var kv KeyValue
	var str []byte
	for key, value := range dict {
		keyEncoded, err := EncodeString([]byte(key))
		switch i := value.(type) {
		case []byte:
			str, err = EncodeString(i)
			if err != nil {
				return nil, err
			}

		case int:
			str, err = EncodeInteger(i)
		case map[string]any:
			str, err = EncodeDictionary(i)
		case []any:
			str, err = EncodeList(i)
		}
		if err != nil {
			return nil, err
		}
		finalValue := append(keyEncoded, str...)
		kv = KeyValue{
			key:   []byte(key),
			value: finalValue,
		}
		if kv.key != nil {
			kvsl.AddSorted(kv)
		}
	}
	encodedDict := []byte{'d'}
	for _, val := range *kvsl {
		value := val.value
		encodedDict = append(encodedDict, value...)
	}
	encodedDict = append(encodedDict, 'e')
	return encodedDict, nil
}
func EncodeList(list []any) ([]byte, error) {
	encodedList := make([]byte, 0)
	encodedList = append(encodedList, 'l')
	var err error
	var encoded []byte
	for _, value := range list {
		switch i := value.(type) {
		case int:
			encoded, err = EncodeInteger(i)
		case []any:
			encoded, err = EncodeList(i)
		case map[string]any:
			encoded, err = EncodeDictionary(i)
		case []byte:
			encoded, err = EncodeString(i)
		}
		if err != nil {
			return encodedList, err
		}
		encodedList = append(encodedList, encoded...)
	}
	encodedList = append(encodedList, 'e')
	return encodedList, nil
}
