package main
import (
  "fmt"
  "encoding/json"
)
type Endpoint struct {
  Host `json:"host"`
  Name `json:"host"`
}

var Endpoints = []Endpoint{}
func init()  {
  jsonFile, err := os.Open("./endpoints.json")
  if err != nil {
      log.Error(err)
      return
  }
  defer jsonFile.Close()
  byteValue, _ := ioutil.ReadAll(jsonFile)
  err = json.Unmarshal([]byte(byteValue), &Endpoints)
  if err != nil {
    log.Error(err)
  }
  log.Infof("The prerule is  %+v\n", Endpoints)
  }
}

func main()  {
  for _, endpoint := range Endpoints {
    fmt.Printf("%v", endpoint)
  }
}
