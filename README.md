# CSVReader

Requires Go 1.20

## Running

```
go build
./csvreader -e test.csv
```

Table format specification: `go doc teval`

Evaluating tables:

### Input

```
,A,B,Cell
1,1,0,1
2,2,=A1+Cell30,0
30,0,=B1+A1,5
```

### Output

```
,A,B,Cell
1,1,0,1
2,2,6,0
30,0,1,5
```
