package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	// "time"
	"encoding/json"
	"github.com/tealeg/xlsx/v3"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Dish struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type DailyMenu struct {
	DailyName string
	Dishes    []Dish
}

type Menu struct {
	Dailys []DailyMenu
}

var rules []Dish

func LoadRule() {
	file, _ := ioutil.ReadFile("./rules.json")
	_ = json.Unmarshal([]byte(file), &rules)
}

func SaveRule() {
	file, _ := json.MarshalIndent(rules, "", " ")
	_ = ioutil.WriteFile("./rules.json", file, 0644)
}

func GetRuleIndex(name string) int {
	for i, d := range rules {
		if d.Name == name {
			return i
		}
	}
	d := Dish{
		Name: name,
	}
	rules = append(rules, d)
	return len(rules) - 1
}
func getImg(url string, name string) {
	imgPath := "./images/"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 1024*1024)
	file, err := os.Create(imgPath + name + ".jpg")
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)
}
func getImageurls(content string) []string {
	urls := []string{}
	reg := regexp.MustCompile(`middleURL":".*?"`)
	reg2 := regexp.MustCompile(`https.*?.jpg`)
	if reg != nil {
		s := reg.FindAllStringSubmatch(content, -1) //-1表示全部匹配
		for _, _s := range s {
			s2 := reg2.FindAllStringSubmatch(_s[0], -1)
			urls = append(urls, s2[0][0])
		}
		// fmt.Println(s)						//[[abc] [aac] [a.c] [a7c] [a c]]
	}
	return urls
}

// func download(name string)  string {
// 	client := &http.Client{}
// 	url := "https://image.baidu.com/search/index?ct=201326592&z=&tn=baiduimage&word="+name+"&pn=0&ie=utf-8&oe=utf-8&cl=2&lm=-1&fr=&se=&sme=&width=640&height=480"
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; BIDUBrowser 2.6)")
// 	req.Header.Set("Referer","https://www.baidu.com")
// 	req.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9")
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Http get err:", err)
//         return ""
// 	}
// 	if resp.StatusCode != 200 {
// 		fmt.Println("Http status code:", resp.StatusCode)
// 		return ""
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Read error", err)
// 		return ""
// 	}
// 	urls := getImageurls(string(body))
// 	return string(body)
// }

func GetDishUrl(name string) string {
	client := &http.Client{}
	url := "https://image.baidu.com/search/index?ct=201326592&z=&tn=baiduimage&word=" + name + "&pn=0&ie=utf-8&oe=utf-8&cl=2&lm=-1&fr=&se=&sme=&width=640&height=480"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; BIDUBrowser 2.6)")
	req.Header.Set("Referer", "https://www.baidu.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}
	urls := getImageurls(string(body))
	return urls[0]
}
func parseexcel() *Menu {
	menu := &Menu{
		Dailys: []DailyMenu{},
	}
	wb, err := xlsx.OpenFile("./菜单.xlsx")
	if err != nil {
		panic(err)
	}
	// wb now contains a reference to the workbook
	// show all the sheets in the workbook
	sh := wb.Sheets[0]
	for i := 0; i <= 10; i++ {
		theCell, _ := sh.Cell(1, i)
		if theCell.String() == "" {
			continue
		}
		daily := DailyMenu{
			DailyName: theCell.String(),
			Dishes:    []Dish{},
		}
		for j := 2; j <= 10; j++ {
			theCell, _ := sh.Cell(j, i)
			if theCell.String() != "" {
				dish := Dish{
					Name: theCell.String(),
				}
				daily.Dishes = append(daily.Dishes, dish)
			}
		}
		menu.Dailys = append(menu.Dailys, daily)
	}
	return menu
}

func GetDishType(name string) string {
	var t string
	fmt.Printf("1: 大荤  2: 小荤  3: 水果  4: 主食 5: 蔬菜\n请输入 [%s] 类型: ", name)
	fmt.Scanln(&t)
	switch t {
	case "1":
		return "大荤"
	case "2":
		return "小荤"
	case "3":
		return "水果"
	case "4":
		return "主食"
	case "5":
		return "蔬菜"
	default:
		return t
	}
}

func main() {
	LoadRule()
	defer SaveRule()
	// defer time.Sleep(60 * time.Second)

	fmt.Println("亲爱的老婆程序启动了，请稍等片刻")
	menu := parseexcel()
	for index, _ := range menu.Dailys {
		for index2, _ := range menu.Dailys[index].Dishes {
			name := menu.Dailys[index].Dishes[index2].Name
			rindex := GetRuleIndex(name)
			if rules[rindex].Type == "" {
				t := GetDishType(name)
				menu.Dailys[index].Dishes[index2].Type = t
				rules[rindex].Type = t
			} else {
				menu.Dailys[index].Dishes[index2].Type = rules[rindex].Type
			}

			if rules[rindex].Url == "" {
				url := GetDishUrl(name)
				menu.Dailys[index].Dishes[index2].Url = url
				rules[rindex].Url = url
			} else {
				menu.Dailys[index].Dishes[index2].Url = rules[rindex].Url
			}
		}
	}

	var filename = "./output.html"
	f, err1 := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	defer f.Close()
	// 读取templlate文件
	t, err := template.ParseFiles("./layouts/main.tmpl", "./layouts/daily.tmpl")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.ExecuteTemplate(f, "daily.tmpl", menu)
	fmt.Println("亲爱的老婆转换完成了，复制ouput.html中的内容到 https://bj.96weixin.com/")
	fmt.Println("可以关闭这个界面了")
}
