package main
import (
  "fmt"
  // "net/http"
  "github.com/sensu/sensu-go/cli/client"
  // "github.com/sensu/sensu-go/cli/client/config"
  // corev2 "github.com/sensu/sensu-go/api/core/v2"
  "github.com/sensu/sensu-go/cli/client/config/inmemory"
)

func main()  {
  conf := inmemory.New("http://127.0.0.1:32613")
  c := client.New(conf)
  fmt.Println(conf.Namespace())
  path := client.EntitiesPath(conf.Namespace())
  // opts := client.ListOptions{}
  // var header http.Header
  // entitys := []corev2.Entity{}

  // err := c.List(path, &entitys, &opts, &header)
  // if err != nil {
  //   fmt.Println(err)
  //   return
  // }
  // for _, entity := range entitys {
  //   fmt.Println(entity.ObjectMeta.Name)
  // }
}
