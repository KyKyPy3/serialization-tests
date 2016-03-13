# serialization-tests
Benchmark tests popular json serializers

#### Tested serializers:
* Native JSON (like etalon)
* JSON + zlib compression
* MsgPack
* Cbor
* Bson

#### Results
```
Benchmark__Encode________Cbor-8         	   50000	     27732 ns/op	 113.66 MB/s	    8368 B/op	      10 allocs/op
Benchmark__Decode________Cbor-8         	   10000	    152000 ns/op	  20.74 MB/s	   16624 B/op	     773 allocs/op
Benchmark__Roundtrip_____Cbor-8         	   10000	    180315 ns/op	  17.48 MB/s	   24992 B/op	     783 allocs/op
Benchmark__Encode________MsgPack-8      	   10000	    116670 ns/op	  27.02 MB/s	   20941 B/op	     351 allocs/op
Benchmark__Decode________MsgPack-8      	   20000	     79670 ns/op	  39.58 MB/s	   19922 B/op	     606 allocs/op
Benchmark__Roundtrip_____MsgPack-8      	   10000	    194603 ns/op	  16.20 MB/s	   41005 B/op	     957 allocs/op
Benchmark__Encode_______JsonCompressed-8	   10000	    168569 ns/op	   0.01 MB/s	   22745 B/op	     533 allocs/op
Benchmark__Decode_______JsonCompressed-8	   10000	    190395 ns/op	   7.20 MB/s	   64403 B/op	     533 allocs/op
Benchmark__Roundtrip____JsonCompressed-8	    5000	    321516 ns/op	   0.01 MB/s	   86867 B/op	    1056 allocs/op
Benchmark__Encode_______Json-8          	   10000	    134020 ns/op	  26.03 MB/s	   22664 B/op	     533 allocs/op
Benchmark__Decode_______Json-8          	   10000	    125861 ns/op	  27.71 MB/s	   23801 B/op	     516 allocs/op
Benchmark__Roundtrip____Json-8          	    5000	    264933 ns/op	  13.17 MB/s	   46455 B/op	    1049 allocs/op
Benchmark__Encode_______Bson-8          	   20000	     88840 ns/op	  42.76 MB/s	   25225 B/op	     359 allocs/op
Benchmark__Decode_______Bson-8          	   10000	    157969 ns/op	  24.05 MB/s	   40362 B/op	    1264 allocs/op
Benchmark__Roundtrip____Bson-8          	   10000	    256740 ns/op	  14.80 MB/s	   65595 B/op	    1623 allocs/op
```
