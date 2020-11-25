# jsonpathgo

Package jsonpathgo provides simple path based getters for accessing JSON primitives

## Functions

### func [GetJsonPathBool](/json.go#L188)

`func GetJsonPathBool(path string, jsonInput []byte) (*bool, error)`

GetJsonPathString returns a bool at a path from JSON.
Given JSON

```go
{"foo":{"123":{"bar":[true,false]}}}
```

path

```go
foo.*.bar[1]
```

would return false

### func [GetJsonPathNumber](/json.go#L209)

`func GetJsonPathNumber(path string, jsonInput []byte) (*float64, error)`

GetJsonPathString returns a number at a path from JSON.
Given JSON

```go
{"foo":{"123":{"bar":[41,42]}}}
```

path

```go
foo.*.bar[1]
```

would return 42

### func [GetJsonPathString](/json.go#L167)

`func GetJsonPathString(path string, jsonInput []byte) (*string, error)`

GetJsonPathString returns a string at a path from JSON.
Given JSON

```go
{"foo":{"123":{"bar":["41","42"]}}}
```

path

```go
foo.*.bar[1]
```

would return "42"

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
