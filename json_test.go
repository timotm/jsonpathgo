package jsonpathgo

import (
	"fmt"
	"testing"
)

func TestGetIndexedKey(t *testing.T) {
	key, index, err := getIndexedKey("foo[3]")
	if err != nil {
		t.Errorf("Expected success, got %+v", err)
	} else if key == nil || index == nil {
		t.Errorf("Expected non-nil return values, got %+v/%+v", key, index)
	} else if *key != "foo" || *index != 3 {
		t.Errorf("Expected 'foo' and 3, got '%s' and %d", *key, *index)
	}
}

func TestGetJsonString(t *testing.T) {
	input := []byte(`{"foo":"bar"}`)

	v, err := getJsonPath("foo", input)

	if err != nil {
		t.Errorf("Expected success, got error '%s'", err)
	} else {
		switch vv := v.(type) {
		case jsonString:
			if vv.Value != "bar" {
				t.Errorf("Expected value of 'bar', got '%+v'", vv)
			}
		default:
			t.Errorf("Expected type JsonString, got %T", v)
		}
	}
}

func TestGetJsonStringFromArray(t *testing.T) {
	input := []byte(`{"foo":["0", "1", "bar"]}`)

	v, err := getJsonPath("foo[2]", input)

	if err != nil {
		t.Errorf("Expected success, got error '%s'", err)
	} else {
		switch vv := v.(type) {
		case jsonString:
			if vv.Value != "bar" {
				t.Errorf("Expected value of 'bar', got '%+v'", vv)
			}
		default:
			t.Errorf("Expected type JsonString, got %T", v)
		}
	}
}

func TestGetString(t *testing.T) {
	input := []byte(`{"foo":{"123":{"bar":["41","42"]}}}`)

	v, err := GetJsonPathString("foo.*.bar[1]", input)

	if err != nil {
		t.Errorf("Expected success, got error '%s'", err)
	} else {
		if v == nil {
			t.Errorf("Expected value, got nil")
		} else {
			if *v != "42" {
				t.Errorf("Expected value of '42', got '%+v'", *v)
			}
		}
	}
}
func TestGetNumber(t *testing.T) {
	input := []byte(`{"foo":42}`)

	v, err := GetJsonPathNumber("foo", input)

	if err != nil {
		t.Errorf("Expected success, got error '%s'", err)
	} else {
		if v == nil {
			t.Errorf("Expected value, got nil")
		} else {
			if *v != 42.0 {
				t.Errorf("Expected value of 42, got '%+v'", v)
			}
		}
	}
}

func TestGetStringNull(t *testing.T) {
	input := []byte(`{"foo":null}`)

	v, err := GetJsonPathString("foo", input)

	if err != nil {
		t.Errorf("Expected success, got error '%s'", err)
	} else {
		if v != nil {
			t.Errorf("Expected nil, got %+v", v)
		}
	}
}
func TestGetJsonNull(t *testing.T) {
	input := []byte(`{"foo":null}`)

	v, err := getJsonPath("foo", input)

	if err != nil {
		t.Errorf("Expected success, got error '%s'", err)
	} else {
		if v != nil {
			t.Errorf("Expected nil, got %+v", v)
		}
	}
}

func TestGetInvalidJson(t *testing.T) {
	input := []byte(`i don't even`)

	_, err := getJsonPath("foo", input)

	if err == nil {
		t.Errorf("Expected failure, got success")
	}
}

func TestGetWildcardPathWithNumber(t *testing.T) {
	input := []byte(`{"foo":{"bar":42}}`)

	v, err := getJsonPath("*.bar", input)

	if err != nil {
		t.Errorf("Expected success, got %s", err)
	} else {
		switch vv := v.(type) {
		case jsonNumber:
			if vv.Value != 42 {
				t.Errorf("Expected value of 42, got %+v", vv)
			}
		default:
			t.Errorf("Expected JsonNumber, got %T", v)
		}
	}
}

func TestGetWildcardEmptyStruct(t *testing.T) {
	input := []byte(`{"foo":{}}`)

	v, err := getJsonPath("*", input)

	if err != nil {
		t.Errorf("Expected success, got %+v", err)
	}

	if v != nil {
		t.Errorf("Expected nil, got %+v", v)
	}
}
func TestGetBoolean(t *testing.T) {
	input := []byte(`{"foo":true, "bar":false}`)
	cases := []struct {
		key      string
		expected bool
	}{
		{key: "foo", expected: true},
		{key: "bar", expected: false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%+v", c.expected), func(t *testing.T) {
			v, err := getJsonPath(c.key, input)

			if err != nil {
				t.Errorf("Expected success, got error '%s'", err)
			} else {
				switch vv := v.(type) {
				case jsonBool:
					if vv.Value != c.expected {
						t.Errorf("Expected value of '%+v', got '%+v'", c.expected, vv)
					}
				default:
					t.Error("Expected type JsonBool, got something else")
				}
			}
		})
	}
}

func TestOverun(t *testing.T) {
	input := []byte(`{"foo":"bar"}`)

	v, err := getJsonPath("foo.humppa", input)

	if err == nil {
		t.Errorf("Expected failure, got %+v", err)
	}

	if v != nil {
		t.Errorf("Expected nil, got %+v", v)
	}

}
func TestComplex(t *testing.T) {
	input := []byte(`{
		"012345678901": {
		  "ident": {
			"type": {
			  "key_localized": "Devicetype",
			  "value_raw": 12,
			  "value_localized": "Oven"
			},
			"deviceName": "",
			"deviceIdentLabel": {
			  "fabNumber": "012345678",
			  "fabIndex": "00",
			  "techType": "H7464BP",
			  "matNumber": "123456",
			  "swids": [
				"4953",
				"20553",
				"25229",
				"4857",
				"25300",
				"25307",
				"25247",
				"20436",
				"25223",
				"4875",
				"20366",
				"20462"
			  ]
			},
			"xkmIdentLabel": {
			  "techType": "EK037",
			  "releaseVersion": "03.85"
			}
		  },
		  "state": {
			"ProgramID": {
			  "value_raw": 24,
			  "value_localized": "",
			  "key_localized": "Program Id"
			},
			"status": {
			  "value_raw": 5,
			  "value_localized": "In use",
			  "key_localized": "State"
			},
			"programType": {
			  "value_raw": 1,
			  "value_localized": "Own program",
			  "key_localized": "Program type"
			},
			"programPhase": {
			  "value_raw": 3073,
			  "value_localized": "PreHeat",
			  "key_localized": "Phase"
			},
			"remainingTime": [
			  0,
			  0
			],
			"startTime": [
			  0,
			  0
			],
			"targetTemperature": [
			  {
				"value_raw": 18000,
				"value_localized": 180,
				"unit": "Celsius"
			  },
			  {
				"value_raw": -32768,
				"value_localized": null,
				"unit": "Celsius"
			  },
			  {
				"value_raw": -32768,
				"value_localized": null,
				"unit": "Celsius"
			  }
			],
			"temperature": [
			  {
				"value_raw": 6967,
				"value_localized": 69.67,
				"unit": "Celsius"
			  },
			  {
				"value_raw": -32768,
				"value_localized": null,
				"unit": "Celsius"
			  },
			  {
				"value_raw": -32768,
				"value_localized": null,
				"unit": "Celsius"
			  }
			],
			"signalInfo": false,
			"signalFailure": false,
			"signalDoor": false,
			"remoteEnable": {
			  "fullRemoteControl": true,
			  "smartGrid": false,
			  "mobileStart": true
			},
			"light": 1,
			"elapsedTime": [
			  0,
			  1
			],
			"spinningSpeed": {
			  "unit": "rpm",
			  "value_raw": null,
			  "value_localized": null,
			  "key_localized": "Spinning Speed"
			},
			"dryingStep": {
			  "value_raw": null,
			  "value_localized": "",
			  "key_localized": "Drying level"
			},
			"ventilationStep": {
			  "value_raw": null,
			  "value_localized": "",
			  "key_localized": "Power Level"
			},
			"plateStep": [],
			"ecoFeedback": null,
			"batteryLevel": null
		  }
		}
	  }`)
	v, err := GetJsonPathNumber("*.state.temperature[0].value_localized", input)
	if err != nil {
		t.Errorf("Expected success, got %+v", err)
	} else {
		if *v != 69.67 {
			t.Errorf("Expected 69.67, got %+v", *v)
		}
	}
}
