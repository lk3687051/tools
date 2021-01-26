package main

import (
	"encoding/json"
	"os/exec"
	"fmt"
	"io/ioutil"
  "os"
	"gopkg.in/yaml.v2"
)

type ModuleArgs struct {
	Path   string
}
type SensuCongig struct {
	Annotations       map[string]string `yaml:"annotations"`
	BackendURL        []string          `yaml:"backend-url"`
	CacheDir          string            `yaml:"cache-dir"`
	DisableAPI        bool              `yaml:"disable-api"`
	DisableSockets    bool              `yaml:"disable-sockets"`
	KeepaliveInterval int               `yaml:"keepalive-interval"`
	Labels            map[string]string `yaml:"labels"`
	LogLevel          string            `yaml:"log-level"`
	Name              string            `yaml:"name"`
	StatsdDisable     bool              `yaml:"statsd-disable"`
	Subscriptions     []string          `yaml:"subscriptions"`
}

type Response struct {
	Msg     string `json:"msg"`
	Changed bool   `json:"changed"`
	Failed  bool   `json:"failed"`
}

func ExitJson(responseBody Response) {
	returnResponse(responseBody)
}

func FailJson(responseBody Response) {
	responseBody.Failed = true
	returnResponse(responseBody)
}

func returnResponse(responseBody Response) {
	var response []byte
	var err error
	response, err = json.Marshal(responseBody)
	if err != nil {
		response, _ = json.Marshal(Response{Msg: "Invalid response object"})
	}
	fmt.Println(string(response))
	if responseBody.Failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func RemoveIndex(s []string, index int) []string {
    ret := make([]string, 0)
    ret = append(ret, s[:index]...)
    return append(ret, s[index+1:]...)
}

func removeRep(slc []string) []string {
    result := []string{}
    tempMap := map[string]int{}
    for _, e := range slc {
			_, ok := tempMap[e]
			if ok {
				continue
			}
      result = append(result, e)
			tempMap[e] = 1
    }
    return result
}

func restart() {
	var response Response
	cmd := exec.Command("sudo", "systemctl", "restart", "sensu-agent")
	err := cmd.Start()
	if err != nil {
		response.Msg = "Can not retart"
		FailJson(response)
	}
	err = cmd.Wait()
}
func config() {
	var response Response
	file, _ := ioutil.ReadFile("/apps/svr/sensu-agent-ops/conf/agent.yml")
	if len(file) == 0 {
		response.Msg = "Config file can not be empty"
		FailJson(response)
	}
	c := SensuCongig{}
	err := yaml.Unmarshal([]byte(file), &c)
	if err != nil {
		response.Msg = "Can not parse config"
		FailJson(response)
	}

	if c.DisableAPI && c.DisableSockets && c.StatsdDisable {
		response.Msg = "Success"
		response.Changed  = false
		ExitJson(response)
	}
	// 初始化默认值
	c.DisableAPI = true
	c.DisableSockets = true
	c.StatsdDisable = true
	c.LogLevel = "fatal"
	restart()
	// 删除重复的Subscriptions
	c.Subscriptions = removeRep(c.Subscriptions)

	d, err := yaml.Marshal(&c)
	if err != nil {
		response.Msg = "Can not create config"
		response.Changed  = false
		FailJson(response)
	}
	ioutil.WriteFile("/apps/svr/sensu-agent-ops/conf/agent.yml", d, 0644)
}

func main() {
	var response Response

	if len(os.Args) != 2 {
		response.Msg = "No argument file provided"
		FailJson(response)
	}

	argsFile := os.Args[1]

	text, err := ioutil.ReadFile(argsFile)
	if err != nil {
		response.Msg = "Could not read configuration file: " + argsFile
		FailJson(response)
	}
	var moduleArgs ModuleArgs
	err = json.Unmarshal(text, &moduleArgs)
	if err != nil {
		response.Msg = "Configuration file not valid JSON: " + argsFile
		FailJson(response)
	}

	config()
	response.Msg = "Success"
	response.Changed  = true
	ExitJson(response)
}
