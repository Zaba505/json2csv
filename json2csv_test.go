package json2csv

import (
  "bytes"
  "strings"
  "testing"
)

func TestConvert(t *testing.T) {
  testCases := []struct{
    Name string
    RawJson string
    ExpectedCsv string
    Options []Option
  }{
    {
      Name: "ArrayOfObjects",
      RawJson: `[
  {
    "hello": "world",
    "goodbye": "world"
  },
  {
    "hello": "bob",
    "goodbye": "bob"
  }
]`,
      ExpectedCsv: `hello,goodbye
world,world
bob,bob
`,
      Options: []Option{
        MapFieldToColumn("hello", 1),
        MapFieldToColumn("goodbye", 2),
      },
    },
  }

  for _, testCase := range testCases {
    t.Run(testCase.Name, func(subT *testing.T) {
      var out bytes.Buffer
      src := strings.NewReader(testCase.RawJson)

      err := Convert(&out, src, testCase.Options...)
      if err != nil {
        subT.Error(err)
        return
      }

      if out.String() != testCase.ExpectedCsv {
        subT.Log(out.String())
        subT.Fail()
        return
      }
    })
  }
}
