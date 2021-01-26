package main
import (
  "fmt"
  "strings"
  "io/ioutil"
  "net/http"
  "github.com/sensu/sensu-go/cli/client"
  "tools/helper/config"
  // "github.com/sensu/sensu-go/cli/client/config"
  corev2 "github.com/sensu/sensu-go/api/core/v2"
  // "github.com/sensu/sensu-go/cli/client/config/inmemory"
)

func WriteHosts(hosts []string) {
  var s string = "[agent]\n"
  for _, host := range hosts {
    s = fmt.Sprintf("%s%s\n",s,host)
  }
  s = fmt.Sprintf("%s%s\n",s,"[agent:vars]\nansible_user=aiops\nansible_password=aiops密码\n")
  ioutil.WriteFile("hosts", []byte(s), 0644)
}
func main()  {
  hosts := []string{}
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
    if entity.EntityClass == "agent" && !strings.Contains(entity.ObjectMeta.Name, "sensu-proxy") {
      hosts = append(hosts, entity.ObjectMeta.Name)
    }
  }
  WriteHosts(hosts)
}
