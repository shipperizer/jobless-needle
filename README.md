# jobless-needle


```
shipperizer in ~/shipperizer/identity-platform-admin-ui on IAM-767 ● λ http :8000/api/v0/counter/5
HTTP/1.1 200 OK
Content-Length: 141
Content-Type: text/plain; charset=utf-8
Date: Fri, 22 Mar 2024 12:42:19 GMT

{
    "count": "50",
    "data": "[0 0 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1 1 1 2 2 2 2 2 2 2 2 2 2 3 3 3 3 3 3 3 3 3 3 4 4 4 4 4 4 4 4 4 4]",
    "status": "200"
}


```

```
shipperizer in ~/shipperizer/identity-platform-admin-ui on IAM-767 ● λ hey -c 100 -n 1000 http://127.0.0.1:8000/api/v0/counter/5

Summary:
  Total:	15.7140 secs
  Slowest:	4.7075 secs
  Fastest:	0.0158 secs
  Average:	1.5222 secs
  Requests/sec:	63.6375
  
  Total data:	141000 bytes
  Size/request:	141 bytes

Response time histogram:
  0.016 [1]	|
  0.485 [36]	|■■■■
  0.954 [70]	|■■■■■■■
  1.423 [404]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.892 [314]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  2.362 [84]	|■■■■■■■■
  2.831 [46]	|■■■■■
  3.300 [18]	|■■
  3.769 [13]	|■
  4.238 [7]	|■
  4.707 [7]	|■


Latency distribution:
  10% in 0.9324 secs
  25% in 1.1607 secs
  50% in 1.4146 secs
  75% in 1.7399 secs
  90% in 2.2801 secs
  95% in 2.7629 secs
  99% in 3.8808 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0027 secs, 0.0158 secs, 4.7075 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:	0.0005 secs, 0.0000 secs, 0.0566 secs
  resp wait:	1.5189 secs, 0.0157 secs, 4.6484 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0004 secs

Status code distribution:
  [200]	1000 responses



```
