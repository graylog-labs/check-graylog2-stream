package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "os"
  "io/ioutil"
  "net/http"
  "net/url"
  "crypto/tls"
  "strings"
  "github.com/fractalcat/nagiosplugin"
)

var condition *string
var stream    *string
var api_url   *string
var user      *string
var pass      *string

func init() {
  condition = flag.String("condition", "<ID>", "Condition ID, set only to check a single alert (optional)")
  stream    = flag.String("stream",    "<ID>", "Stream ID (mandatory)")
  api_url   = flag.String("url",       "http://localhost:12900", "URL to Graylog2 api (optional)")
  user      = flag.String("user",      "<username>", "API username (mandatory)")
  pass      = flag.String("password",  "<password>", "API password (mandatory)")
}

func main() {
  flag.Parse()
  checkArguments()

  check := nagiosplugin.NewCheck()
  defer check.Finish()

  var data map[string]interface{}
  queryApi(parseUrl(*api_url) + "/streams/" + url.QueryEscape(*stream) + "/alerts/check", url.QueryEscape(*user), url.QueryEscape(*pass), &data)

  total   := data["total_triggered"].(float64)
  results := data["results"].([]interface{})

  var stream_title string

  queryApi(parseUrl(*api_url) + "/streams/" + url.QueryEscape(*stream), url.QueryEscape(*user), url.QueryEscape(*pass), &data)
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

func parseUrl(unparsed_url string) string {
  parsed_url, err := url.Parse(unparsed_url)
  if err != nil {
    nagiosplugin.Exit(nagiosplugin.UNKNOWN, "Can not parse given URL")
  }

  if !strings.Contains(parsed_url.Host, ":") {
    nagiosplugin.Exit(nagiosplugin.UNKNOWN, "Please give the API port number in the form http://hostname:port")
  }

  connection_string := parsed_url.Scheme + "://" + parsed_url.Host + parsed_url.Path
  return connection_string
}

func queryApi(api_url string, user string, pass string, data *map[string]interface{}) {
  tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }
  client := &http.Client{Transport: tr}
  req, err := http.NewRequest("GET", api_url, nil)
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
