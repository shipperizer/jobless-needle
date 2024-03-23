## Single endpoint invocation

```bash
~/GolandProjects/jobless-needle$ curl http://127.0.0.1:8000/api/v0/tp/50 | jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  1181  100  1181    0     0   890k      0 --:--:-- --:--:-- --:--:-- 1153k

{
  "count": "50",
  "data": "[{key:Task-1 value:10} {key:Task-0 value:19} {key:Task-9 value:19} {key:Task-2 value:19} {key:Task-3 value:19} {key:Task-4 value:19} {key:Task-5 value:19} {key:Task-6 value:19} {key:Task-7 value:19} {key:Task-8 value:19} {key:Task-13 value:19} {key:Task-10 value:19} {key:Task-11 value:19} {key:Task-12 value:19} {key:Task-15 value:19} {key:Task-16 value:19} {key:Task-17 value:19} {key:Task-21 value:26} {key:Task-18 value:27} {key:Task-19 value:28} {key:Task-20 value:29} {key:Task-25 value:31} {key:Task-22 value:32} {key:Task-23 value:33} {key:Task-24 value:33} {key:Task-29 value:35} {key:Task-26 value:35} {key:Task-27 value:35} {key:Task-28 value:35} {key:Task-31 value:35} {key:Task-30 value:35} {key:Task-32 value:35} {key:Task-33 value:35} {key:Task-49 value:50} {key:Task-34 value:50} {key:Task-35 value:50} {key:Task-36 value:50} {key:Task-37 value:50} {key:Task-38 value:50} {key:Task-39 value:50} {key:Task-40 value:50} {key:Task-41 value:50} {key:Task-42 value:50} {key:Task-43 value:50} {key:Task-44 value:50} {key:Task-45 value:50} {key:Task-46 value:50} {key:Task-47 value:50} {key:Task-48 value:50} {key:Task-14 value:50}]",
  "status": "200"
}


```

## Test con Hey
```bash
~/GolandProjects/jobless-needle$ hey -c 100 -n 1000 http://127.0.0.1:8000/api/v0/tp/50

Summary:
  Total:        0.0911 secs
  Slowest:      0.0502 secs
  Fastest:      0.0002 secs
  Average:      0.0079 secs
  Requests/sec: 10979.9114
  
  Total data:   680821 bytes
  Size/request: 680 bytes

Response time histogram:
  0.000 [1]     |
  0.005 [559]   |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.010 [170]   |■■■■■■■■■■■■
  0.015 [84]    |■■■■■■
  0.020 [44]    |■■■
  0.025 [50]    |■■■■
  0.030 [57]    |■■■■
  0.035 [27]    |■■
  0.040 [4]     |
  0.045 [3]     |
  0.050 [1]     |


Latency distribution:
  10% in 0.0004 secs
  25% in 0.0009 secs
  50% in 0.0038 secs
  75% in 0.0111 secs
  90% in 0.0243 secs
  95% in 0.0291 secs
  99% in 0.0339 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0002 secs, 0.0502 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0000 secs
  req write:    0.0001 secs, 0.0000 secs, 0.0054 secs
  resp wait:    0.0075 secs, 0.0002 secs, 0.0501 secs
  resp read:    0.0002 secs, 0.0000 secs, 0.0071 secs

Status code distribution:
  [200] 1000 responses

```
