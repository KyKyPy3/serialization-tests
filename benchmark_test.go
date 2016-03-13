package benchmark

import (
	"bytes"
	"github.com/klauspost/compress/zlib"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/ugorji/go/codec"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/vmihailenco/msgpack.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func loadJson() map[string]interface{} {
	var jsonStruct map[string]interface{}

	rawJson, err := ioutil.ReadFile("test.json")
	if err != nil {
		log.Fatalf("Failed to open file test.json: %s", err.Error())
	}

	ffjson.Unmarshal(rawJson, &jsonStruct)

	return jsonStruct
}

func validate(first interface{}, second interface{}) (result bool) {

	result = true

	switch reflect.TypeOf(first).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(first)

		if len(second.([]interface{})) == s.Len() {

			for i := 0; i < s.Len(); i++ {
				result = result && validate(s.Index(i).Interface(), second.([]interface{})[i])
			}
		} else {
			result = false
		}

	case reflect.Map:
		for _, key := range reflect.ValueOf(first).MapKeys() {
			if val := reflect.ValueOf(second).MapIndex(key); val.IsValid() {
				result = result && validate(reflect.ValueOf(first).MapIndex(key).Interface(), val.Interface())
			}
		}
	case reflect.Float64:
		result = (first.(float64) == second)
	case reflect.String:
		result = (first.(string) == second)
	}

	return
}

//
// Cbor
//

func Benchmark__Encode________Cbor(b *testing.B) {

	jsonStruct := loadJson()

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf []byte

		err := codec.NewEncoderBytes(&buf, new(codec.CborHandle)).Encode(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}
}

func Benchmark__Decode________Cbor(b *testing.B) {

	jsonStruct := loadJson()
	var buf []byte
	var data map[string]interface{}

	err := codec.NewEncoderBytes(&buf, new(codec.CborHandle)).Encode(jsonStruct)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(buf)))

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := codec.NewDecoderBytes(buf, new(codec.CborHandle)).Decode(&data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark__Roundtrip_____Cbor(b *testing.B) {

	jsonStruct := loadJson()
	var data map[string]interface{}

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf []byte

		err := codec.NewEncoderBytes(&buf, new(codec.CborHandle)).Encode(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		err = codec.NewDecoderBytes(buf, new(codec.CborHandle)).Decode(&data)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}

	if len(os.Getenv("VALIDATE")) > 0 {
		if !validate(jsonStruct, data) {
			b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
		}
	}
}

//
// MsgPack
//

func Benchmark__Encode________MsgPack(b *testing.B) {

	jsonStruct := loadJson()

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := msgpack.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}
}

func Benchmark__Decode________MsgPack(b *testing.B) {

	jsonStruct := loadJson()
	var data map[string]interface{}

	buf, err := msgpack.Marshal(jsonStruct)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(buf)))

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := msgpack.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark__Roundtrip_____MsgPack(b *testing.B) {

	jsonStruct := loadJson()
	var data map[string]interface{}

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := msgpack.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		err = msgpack.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}

	if len(os.Getenv("VALIDATE")) > 0 {
		if !validate(jsonStruct, data) {
			b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
		}
	}
}

//
// Json compressed
//

func Benchmark__Encode_______JsonCompressed(b *testing.B) {

	var jsonGZ bytes.Buffer

	zipper := zlib.NewWriter(&jsonGZ)
	jsonStruct := loadJson()

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonGZ.Reset()
		zipper.Reset(&jsonGZ)

		buf, err := ffjson.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		_, err = zipper.Write(buf)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(jsonGZ.Len()))
		}
	}

	err := zipper.Close()
	if err != nil {
		b.Fatal(err)
	}
}

func Benchmark__Decode_______JsonCompressed(b *testing.B) {

	var jsonGZ bytes.Buffer
	var out bytes.Buffer

	zipper := zlib.NewWriter(&jsonGZ)
	jsonStruct := loadJson()
	var data map[string]interface{}

	buf, err := ffjson.Marshal(jsonStruct)
	if err != nil {
		b.Fatal(err)
	}
	_, err = zipper.Write(buf)
	if err != nil {
		b.Fatal(err)
	}

	err = zipper.Close()
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(jsonGZ.Len()))

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		out.Reset()

		byteReader := bytes.NewReader(jsonGZ.Bytes())
		r, err := zlib.NewReader(byteReader)
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(&out, r)

		err = ffjson.Unmarshal(out.Bytes(), &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark__Roundtrip____JsonCompressed(b *testing.B) {

	var jsonGZ bytes.Buffer
	var out bytes.Buffer
	zipper := zlib.NewWriter(&jsonGZ)

	jsonStruct := loadJson()
	var data map[string]interface{}

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonGZ.Reset()
		zipper.Reset(&jsonGZ)

		buf, err := ffjson.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		_, err = zipper.Write(buf)
		if err != nil {
			b.Fatal(err)
		}

		out.Reset()

		byteReader := bytes.NewReader(jsonGZ.Bytes())
		r, err := zlib.NewReader(byteReader)
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(&out, r)

		err = ffjson.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(jsonGZ.Len()))
		}
	}

	err := zipper.Close()
	if err != nil {
		b.Fatal(err)
	}

	if len(os.Getenv("VALIDATE")) > 0 {
		if !reflect.DeepEqual(jsonStruct, data) {
			b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
		}
	}
}

//
// Json
//

func Benchmark__Encode_______Json(b *testing.B) {

	jsonStruct := loadJson()

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := ffjson.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}
}

func Benchmark__Decode_______Json(b *testing.B) {

	jsonStruct := loadJson()
	var data map[string]interface{}

	buf, err := ffjson.Marshal(jsonStruct)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(buf)))

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := ffjson.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark__Roundtrip____Json(b *testing.B) {

	jsonStruct := loadJson()
	var data map[string]interface{}

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := ffjson.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		err = ffjson.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}

	if len(os.Getenv("VALIDATE")) > 0 {
		if !reflect.DeepEqual(jsonStruct, data) {
			b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
		}
	}
}

//
// Bson
//

func Benchmark__Encode_______Bson(b *testing.B) {

	jsonStruct := loadJson()

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := bson.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}
}

func Benchmark__Decode_______Bson(b *testing.B) {

	jsonStruct := loadJson()
	var data map[string]interface{}

	buf, err := bson.Marshal(jsonStruct)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(buf)))

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := bson.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark__Roundtrip____Bson(b *testing.B) {

	jsonStruct := loadJson()
	var data map[string]interface{}

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := bson.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}
		err = bson.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}

	if len(os.Getenv("VALIDATE")) > 0 {
		if !reflect.DeepEqual(jsonStruct, data) {
			b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
		}
	}
}
