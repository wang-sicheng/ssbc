sflow [![Circle CI](https://circleci.com/gh/Cistern/sflow.svg?style=svg&circle-token=e6b83ff5665619ad7615dd1e472c958f542cca3c)](https://circleci.com/gh/Cistern/sflow) [![GoDoc](https://godoc.org/github.com/PreetamJinka/sflow?status.svg)](https://godoc.org/github.com/PreetamJinka/sflow) [![BSD License](https://img.shields.io/pypi/l/Django.svg)](https://github.com/PreetamJinka/sflow/blob/master/LICENSE)
====

An [sFlow](http://sflow.org/) v5 encoding and decoding package for Go.

Usage
---

```go
// Create a new decoder that reads from an io.Reader.
d := sflow.NewDecoder(r)

// Attempt to decode an sFlow datagram.
dgram, err := d.Decode()
if err != nil {
	log.Println(err)
	return
}

for _, sample := range dgram.Samples {
	// Sample is an interface type
	if sample.SampleType() == sflow.TypeCounterSample {
		counterSample := sample.(sflow.CounterSample)

		for _, record := range counterSample.Records {
			// While there is a record.RecordType() method,
			// you can always check types directly.
			switch record.(type) {
			case sflow.HostDiskCounters:
				fmt.Printf("Max used percent of disk space is %d.\n",
					record.MaxUsedPercent)
			}
		}
	}
}
```

API guarantees
---
API stability is *not guaranteed*. Vendoring or using a dependency manager is suggested.

Reporting issues
---
Bug reports are greatly appreciated. Please provide raw datagram dumps when possible.

License
---
BSD (see LICENSE)
