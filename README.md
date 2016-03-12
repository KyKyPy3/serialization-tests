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
Benchmark__Encode________Cbor-8         	  300000	      4101 ns/op	  27.55 MB/s	    2897 B/op	       8 allocs/op
Benchmark__Decode________Cbor-8         	  300000	      4741 ns/op	  23.83 MB/s	    3104 B/op	      11 allocs/op
Benchmark__Roundtrip_____Cbor-8         	  200000	      9750 ns/op	  11.59 MB/s	    6003 B/op	      19 allocs/op
Benchmark__Encode________MsgPack-8      	  500000	      2429 ns/op	  46.51 MB/s	     416 B/op	       5 allocs/op
Benchmark__Decode________MsgPack-8      	  500000	      3506 ns/op	  32.22 MB/s	     400 B/op	      18 allocs/op
Benchmark__Roundtrip_____MsgPack-8      	  200000	      6712 ns/op	  16.83 MB/s	     816 B/op	      23 allocs/op
Benchmark__Encode_______JsonCompressed-8	   50000	     36774 ns/op	   0.05 MB/s	     728 B/op	       5 allocs/op
Benchmark__Decode_______JsonCompressed-8	   50000	     24181 ns/op	   5.50 MB/s	   40656 B/op	      19 allocs/op
Benchmark__Roundtrip____JsonCompressed-8	   30000	     54881 ns/op	   0.04 MB/s	   41395 B/op	      24 allocs/op
Benchmark__Encode_______Json-8          	  500000	      2878 ns/op	  51.07 MB/s	     712 B/op	       5 allocs/op
Benchmark__Decode_______Json-8          	  200000	      6296 ns/op	  23.34 MB/s	     416 B/op	      12 allocs/op
Benchmark__Roundtrip____Json-8          	  200000	     10343 ns/op	  14.21 MB/s	    1128 B/op	      17 allocs/op
Benchmark__Encode_______Bson-8          	  500000	      3307 ns/op	  52.91 MB/s	     608 B/op	       8 allocs/op
Benchmark__Decode_______Bson-8          	  300000	      5464 ns/op	  32.02 MB/s	     752 B/op	      43 allocs/op
Benchmark__Roundtrip____Bson-8          	  200000	      9503 ns/op	  18.41 MB/s	    1360 B/op	      51 allocs/op
```
