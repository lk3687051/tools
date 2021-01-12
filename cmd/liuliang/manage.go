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
  path := client.EntitiesPath(conf.Namespace())
  opts := client.ListOptions{}
  var header http.Header
  entitys := []corev2.Entity{}

  err = c.List(path, &entitys, &opts, &header)
  if err != nil {
    fmt.Println(err)
    return
  }
  for _, entity := range entitys {
    if entity.EntityClass == corev2.EntityProxyClass {
      if len(entity.Subscriptions) == 1 {
        fmt.Println(entity.ObjectMeta.Name)
        c.DeleteEntity(entity.ObjectMeta.Namespace, entity.ObjectMeta.Name)
      }
      // fmt.Println(entity.ObjectMeta.Name)
    }
  }
}
