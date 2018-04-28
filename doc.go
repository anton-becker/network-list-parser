package main

const helpMessage string = `This program:
1) parses text file with set of ip addresses and/or networks,
2) normalizes them,
3) auto aggregates (optionally),
4) prints resulted set.

[Input file format]

Program scans file line by line and parse beginning of line (until first ";" or EOL) for "|"-separated list of ip addresses and/or networks.
Spaces outside ip addresses and networks are ignored.
Only IPv4 is supported.
Networks must be defined in prefix notation ("127.0.0.1/24").
If there is no prefix ("127.0.0.1") then prefix "/32" will be added.
Some parser's hacks (in order):
1) it silently ignore any line beginning with "Updated:",
2) if first char in line is ";" when this char will be removed before parsing.

[Normalization]

Normalization removes all duplications. Also if network "A" contains network "B" then network "B" will be removed.
Each network itself also will be normalized according their mask: each non significant bits will be set to 0.

[Auto aggregation]

Auto aggregation reduces number of addresses/networks in cost of some addresses which is not belong to original set will be added ("fake" IPs).
Each address which belongs to original set will belongs to resulting set. 
There are 2 type of auto aggregation: intensive and usual.

Intensive aggregation applied if all 3 conditions are satisfied:
1) result network's prefix >= "intensive-aggregation-min-prefix",
2) at least "intensive-aggregation-min-nets" nets from original set will be covered by new network,
3) no more than "intensive-aggregation-max-fake-ips" fake IPs will be covered/added by result network.

Usual aggregation computes acceptable percent ("AP") of fake IPs for resulting network.
Usual aggregation applied if all 2 conditions are satisfied:
1) number of fake IPs no more than "aggregation-max-fake-ips",
2) percent of fake IPs no more than "AP".
"AP" is non-linear function of <number of nets from original set which will be covered by new network> defined by 2 points:
1) If number of nets = 2 then "AP" = "aggregation-lo-fake-percent",
2) If number of nets = "aggregation-hi-fake-percent-nets" then "AP" = "aggregation-hi-fake-percent".
If number of nets = "aggregation-hi-fake-percent-nets" then "AP" also = "aggregation-hi-fake-percent".

[Printing result]

Result set is printed into "dst-file" in the following format:

"""
<prefix><network  1  address><delimiter><network  1  prefix or mask><suffix><network-delimiter>
<prefix><network  2  address><delimiter><network  2  prefix or mask><suffix><network-delimiter>
...
<prefix><network n-1 address><delimiter><network n-1 prefix or mask><suffix><network-delimiter>
<prefix><network  n  address><delimiter><network  n  prefix or mask><suffix>
"""

What will be printed - prefix or mask - depends on "mask-notation".
Result set is always sorted be network address.

[Usage]

`
