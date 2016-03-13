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
Benchmark__Encode________Cbor-8         	   20000	     71203 ns/op	 117.48 MB/s	   32944 B/op	      12 allocs/op
Benchmark__Decode________Cbor-8         	    3000	    432005 ns/op	  19.36 MB/s	   40216 B/op	    2159 allocs/op
Benchmark__Roundtrip_____Cbor-8         	    3000	    531434 ns/op	  15.74 MB/s	   73160 B/op	    2171 allocs/op
Benchmark__Encode________MsgPack-8      	    5000	    306260 ns/op	  27.36 MB/s	   49907 B/op	     968 allocs/op
Benchmark__Decode________MsgPack-8      	   10000	    214967 ns/op	  38.98 MB/s	   59649 B/op	    1690 allocs/op
Benchmark__Roundtrip_____MsgPack-8      	    3000	    547434 ns/op	  15.31 MB/s	  109754 B/op	    2658 allocs/op
Benchmark__Encode_______JsonCompressed-8	    3000	    422254 ns/op	   0.00 MB/s	   54790 B/op	    1478 allocs/op
Benchmark__Decode_______JsonCompressed-8	    3000	    464569 ns/op	   4.50 MB/s	  117206 B/op	    1446 allocs/op
Benchmark__Roundtrip____JsonCompressed-8	    2000	    881305 ns/op	   0.00 MB/s	  171795 B/op	    2915 allocs/op
Benchmark__Encode_______Json-8          	    5000	    391100 ns/op	  24.37 MB/s	   54520 B/op	    1478 allocs/op
Benchmark__Decode_______Json-8          	    3000	    392785 ns/op	  24.27 MB/s	   76634 B/op	    1430 allocs/op
Benchmark__Roundtrip____Json-8          	    2000	    796115 ns/op	  11.97 MB/s	  131150 B/op	    2908 allocs/op
Benchmark__Encode_______Bson-8          	   10000	    249024 ns/op	  40.42 MB/s	   70085 B/op	     982 allocs/op
Benchmark__Decode_______Bson-8          	    3000	    524408 ns/op	  19.19 MB/s	  124777 B/op	    3580 allocs/op
Benchmark__Roundtrip____Bson-8          	    2000	    712524 ns/op	  14.13 MB/s	  195334 B/op	    4562 allocs/op
```
