package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "os"
  "io/ioutil"
  "net/http"
	"github.com/fractalcat/nagiosplugin"
)

var condition *string
var stream    *string
var url       *string
var user      *string
var pass      *string

func init() {
  condition = flag.String("condition", "<ID>", "Condition ID, set only to check a single alert (optional)")
  stream    = flag.String("stream",    "<ID>", "Stream ID (mandatory)")
  url       = flag.String("url",       "http://localhost:12900", "URL to Graylog2 api (optional)")
  user      = flag.String("user",      "<username>", "API username (mandatory)")
  pass      = flag.String("password",  "<password>", "API password (mandatory)")
}

func main() {
  flag.Parse()
  checkArguments()

  check := nagiosplugin.NewCheck()
  defer check.Finish()

  var data map[string]interface{}
  queryApi(*url + "/streams/" + *stream + "/alerts/check", *user, *pass, &data)

  total   := data["total_triggered"].(float64)
  results := data["results"].([]interface{})

  var stream_title string
  queryApi(*url + "/streams/" + *stream, *user, *pass, &data)
  stream_title = data["title"].(string)

  for i, result := range results {
    i = i
    mappedResult    := result.(map[string]interface{})
    mappedCondition := mappedResult["condition"].(map[string]interface{})
    if *condition != "<ID>" && *condition == mappedCondition["id"] && mappedResult["triggered"] == true {
      nagiosplugin.Exit(nagiosplugin.CRITICAL, "Alert triggered for stream '" + stream_title + "' condition: " + *condition)
    }
  }

  if total > 0  && *condition == "<ID>"{
    nagiosplugin.Exit(nagiosplugin.CRITICAL, fmt.Sprintf("%g", total) + " alert/s triggered for stream " + stream_title)
  }

  check.AddResult(nagiosplugin.OK, "No stream alerts triggered for stream: " + stream_title)
}

func checkArguments() {
  if *stream == "<ID>" {
    fmt.Println("usage:")
    flag.PrintDefaults()
    os.Exit(1)
  }
}

func queryApi(url string, user string, pass string, data *map[string]interface{}) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", url, nil)
  req.SetBasicAuth(user, pass)
  res, err := client.Do(req)
  if err != nil {
    nagiosplugin.Exit(nagiosplugin.UNKNOWN, "Can not connect to Graylog2 API")
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    nagiosplugin.Exit(nagiosplugin.UNKNOWN, "Got no response from Graylog2 API")
  }

  err = json.Unmarshal(body, data)
  if err != nil {
    nagiosplugin.Exit(nagiosplugin.UNKNOWN, "Can not parse JSON from Graylog2 API")
  }

  if res.StatusCode != 200 {
    nagiosplugin.Exit(nagiosplugin.UNKNOWN, "Got wrong return code from Graylog2 API, please check all command line parameters")
  }
}
