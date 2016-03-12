package benchmark

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/klauspost/compress/zlib"
	"github.com/ugorji/go/codec"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Nested struct {
	Child1 string
	Child2 string
}

type BenchStruct struct {
	Name      string
	Age       int
	Hobby     []string
	Married   bool
	Weight    float64
	Childrens Nested
}

func getStructRepresentation() BenchStruct {
	return BenchStruct{
		Name:    "Jack",
		Age:     15,
		Hobby:   []string{"books", "films", "skiing", "swimming"},
		Married: true,
		Weight:  1.65,
		Childrens: Nested{
			Child1: "Robert",
			Child2: "Maks",
		},
	}
}

//
// Cbor
//

func Benchmark__Encode________Cbor(b *testing.B) {

	jsonStruct := getStructRepresentation()

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

	jsonStruct := getStructRepresentation()
	var buf []byte
	var data BenchStruct

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

	jsonStruct := getStructRepresentation()
	var data BenchStruct

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

	if !reflect.DeepEqual(jsonStruct, data) {
		b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
	}
}

//
// MsgPack
//

func Benchmark__Encode________MsgPack(b *testing.B) {

	jsonStruct := getStructRepresentation()

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

	jsonStruct := getStructRepresentation()
	var data BenchStruct

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

	jsonStruct := getStructRepresentation()
	var data BenchStruct

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

	if !reflect.DeepEqual(jsonStruct, data) {
		b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
	}
}

//
// Json compressed
//

func Benchmark__Encode_______JsonCompressed(b *testing.B) {

	var jsonGZ bytes.Buffer

	zipper := zlib.NewWriter(&jsonGZ)
	jsonStruct := getStructRepresentation()

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonGZ.Reset()
		zipper.Reset(&jsonGZ)

		buf, err := json.Marshal(jsonStruct)
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
	jsonStruct := getStructRepresentation()
	var data BenchStruct

	buf, err := json.Marshal(jsonStruct)
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

		err = json.Unmarshal(out.Bytes(), &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark__Roundtrip____JsonCompressed(b *testing.B) {

	var jsonGZ bytes.Buffer
	var out bytes.Buffer
	zipper := zlib.NewWriter(&jsonGZ)

	jsonStruct := getStructRepresentation()
	var data BenchStruct

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonGZ.Reset()
		zipper.Reset(&jsonGZ)

		buf, err := json.Marshal(jsonStruct)
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

		err = json.Unmarshal(buf, &data)
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

	if !reflect.DeepEqual(jsonStruct, data) {
		b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
	}
}

//
// Json
//

func Benchmark__Encode_______Json(b *testing.B) {

	jsonStruct := getStructRepresentation()

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := json.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}
}

func Benchmark__Decode_______Json(b *testing.B) {

	jsonStruct := getStructRepresentation()
	var data BenchStruct

	buf, err := json.Marshal(jsonStruct)
	if err != nil {
		b.Fatal(err)
	}

	b.SetBytes(int64(len(buf)))

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark__Roundtrip____Json(b *testing.B) {

	jsonStruct := getStructRepresentation()
	var data BenchStruct

	b.ReportAllocs()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := json.Marshal(jsonStruct)
		if err != nil {
			b.Fatal(err)
		}

		err = json.Unmarshal(buf, &data)
		if err != nil {
			b.Fatal(err)
		}

		if i == 0 {
			b.SetBytes(int64(len(buf)))
		}
	}

	if !reflect.DeepEqual(jsonStruct, data) {
		b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
	}
}

//
// Bson
//

func Benchmark__Encode_______Bson(b *testing.B) {

	jsonStruct := getStructRepresentation()

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

	jsonStruct := getStructRepresentation()
	var data BenchStruct

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

	jsonStruct := getStructRepresentation()
	var data BenchStruct

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

	if !reflect.DeepEqual(jsonStruct, data) {
		b.Fatalf("Unmarshaled object differed:\n%v\n%v", jsonStruct, data)
	}
}
