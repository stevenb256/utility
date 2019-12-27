package utility

import (
	"bytes"
	"encoding/gob"
)

// GobEncode used to encode an object in golang gob format
func GobEncode(object interface{}) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(object)
	if nil != err {
		return nil, err
	}
	return w.Bytes(), nil
}

// GobDecode used to deserialize an object encoded with GobEncode
func GobDecode(buf []byte, object interface{}) error {
	return gob.NewDecoder(bytes.NewReader(buf)).Decode(object)
}
