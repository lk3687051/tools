package main
import (
  "fmt"
  // "strings"
  "net/http"
  "github.com/sensu/sensu-go/cli/client"
  "tools/helper/config"
  // "github.com/sensu/sensu-go/cli/client/config"
  corev2 "github.com/sensu/sensu-go/api/core/v2"
  // "github.com/sensu/sensu-go/cli/client/config/inmemory"
)


func main()  {
  conf := config.New("http://127.0.0.1:32613")
  c := client.New(conf)
  tokens, err := c.CreateAccessToken(
    conf.APIUrl(), "admin", "U4tmN2R*gf3k",
  )
  conf.SaveTokens(tokens)
  path := client.ChecksPath(conf.Namespace())
  opts := client.ListOptions{}
  var header http.Header
  checks := []corev2.CheckConfig{}

  err = c.List(path, &checks, &opts, &header)
  if err != nil {
    fmt.Println(err)
    return
  }
  for _, check := range checks {
    if check.Ttl != 0 || check.Timeout == 0 {
      fmt.Printf("The check is %s, ttl is %d\n", check.ObjectMeta.Name, check.Ttl)
      check.Ttl = 0
      check.Timeout = check.Interval
      c.UpdateCheck(&check)
    }
  }
  
  for _, check := range checks {
    if check.ObjectMeta.Name == "check_network" {
      check.Ttl = 0
      check.Interval = 300
      c.UpdateCheck(&check)
    }
  }
}
