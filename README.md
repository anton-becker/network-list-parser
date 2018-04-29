# network-list-parser
Parse, normalize &amp; aggregate list of IPv4 networks/addresses

This program:
1) parses text file with set of ip addresses and/or networks,
2) normalizes them,
3) auto aggregates (optionally),
4) saves result set.

Example of normalization & auto aggregation with default settings:
    
    File contains: 267593 networks/IPs with 15974399 IPs coverage.
    After normalization: 63642 networks/IPs with 14740205 IPs coverage.
    After auto aggregate: 15919 networks/IPs with 15501918 IPs coverage. Removed 75% records (47723). Add 5.2% IPs coverage (761713).

See [doc.go](../master/doc.go) for details.
