package utility

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/numbatx/gn-core/core"
	"github.com/elodina/go-avro"
)

// HexSliceToByteSlice outputs a decoded byte slice representation of a hex string encoded slice input
func HexSliceToByteSlice(in []string) ([][]byte, error) {
	if in == nil {
		return nil, nil
	}
	out := make([][]byte, len(in))

	for i := range in {
		tmp, err := hex.DecodeString(in[i])
		if err != nil {
			return nil, err
		}
		out[i] = tmp
	}

	return out, nil
}

// UIntSliceToIntSlice outputs the int64 slice representation of a uint64 slice input
func UIntSliceToIntSlice(in []uint64) []int64 {
	if in == nil {
		return nil
	}

	out := make([]int64, len(in))

	for i := range in {
		out[i] = int64(in[i])
	}

	return out
}

// GetBytes returns the bytes representation of a big int input if not nil, otherwise returns []byte{}
func GetBytes(val *big.Int) []byte {
	if val != nil {
		return val.Bytes()
	}

	return big.NewInt(0).Bytes()
}

// EncodePubKey returns a byte slice of the encoded pubKey input, using a pub key converter
func EncodePubKey(pubKeyConverter core.PubkeyConverter, pubKey []byte) []byte {
	return []byte(pubKeyConverter.Encode(pubKey))
}

// Encode returns a byte slice representing the binary encoding of the input avro record
func Encode(record avro.AvroRecord) ([]byte, error) {
	writer := avro.NewSpecificDatumWriter()
	writer.SetSchema(record.Schema())

	buffer := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(buffer)

	err := writer.Write(record, encoder)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode tries to decode a data buffer, read it and store it on the input record.
// If successfully, the record is filled with data from the buffer, otherwise an error might be returned
func Decode(record avro.AvroRecord, buffer []byte) error {
	reader := avro.NewSpecificDatumReader()
	reader.SetSchema(record.Schema())

	decoder := avro.NewBinaryDecoder(buffer)
	return reader.Read(record, decoder)
}

// MetaChainShardAddress returns core.MetachainShardId as a 62 byte array address(by padding with zeros).
// This is needed, since all addresses from avro schema are required to be 62 fixed byte array
func MetaChainShardAddress() []byte {
	ret := make([]byte, 62)
	copy(ret, fmt.Sprintf("%d", core.MetachainShardId))
	return ret
}
