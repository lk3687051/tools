package main
import (
	"fmt"
	"bufio"
	"os"
	"io"
	"time"
	"regexp"
	"net/http"
	"io/ioutil"
	"html/template"
	"github.com/tealeg/xlsx/v3"
)

type Dish struct {
	Url string
	Type string
	Name string
}

type DailyMenu struct {
	DailyName string
	Dishes  []Dish
}

type Menu struct {
	Dailys []DailyMenu
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
	reader := bufio.NewReaderSize(res.Body, 1024 * 1024)
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
	if(reg != nil){
		s := reg.FindAllStringSubmatch(content,-1)   //-1表示全部匹配
		for _, _s := range s {
			s2 := reg2.FindAllStringSubmatch(_s[0], -1)
			urls = append(urls,s2[0][0])
		}
		// fmt.Println(s)						//[[abc] [aac] [a.c] [a7c] [a c]]
	}
	return urls
}
func download(name string)  string {
	client := &http.Client{}
	url := "https://image.baidu.com/search/index?ct=201326592&z=&tn=baiduimage&word="+name+"&pn=0&ie=utf-8&oe=utf-8&cl=2&lm=-1&fr=&se=&sme=&width=640&height=480"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; BIDUBrowser 2.6)")
	req.Header.Set("Referer","https://www.baidu.com")
	req.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9")
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
	getImg(urls[0], name)
	return string(body)
}

func GetMenuUrl(name string)  string {
	client := &http.Client{}
	url := "https://image.baidu.com/search/index?ct=201326592&z=&tn=baiduimage&word="+name+"&pn=0&ie=utf-8&oe=utf-8&cl=2&lm=-1&fr=&se=&sme=&width=640&height=480"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; BIDUBrowser 2.6)")
	req.Header.Set("Referer","https://www.baidu.com")
	req.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.9")
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
		theCell,_ := sh.Cell(1, i)
		if theCell.String() == "" {
			continue
		}
		daily := DailyMenu{
			DailyName: theCell.String(),
			Dishes: []Dish{},
		}
		for j := 2; j <= 10; j++ {
			theCell,_ := sh.Cell(j, i)
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

func ReadRule() map[string]Dish {
	rules := map[string]Dish{}
	wb, _ := xlsx.OpenFile("./rule.xlsx")
	sh := wb.Sheets[0]

	for i := 0; i < sh.MaxRow; i++ {
		dish := Dish{}
		NameCell,_ := sh.Cell(i, 0)
		dish.Name = NameCell.String()
		TypeCell,_ := sh.Cell(i, 1)
		dish.Type = TypeCell.String()
		UrlCell,_ := sh.Cell(i, 2)
		dish.Url = UrlCell.String()
		rules[dish.Name] = dish
	}

	return rules
}
func main()  {
	fmt.Println("亲爱的老婆程序启动了，请稍等片刻")
	defer time.Sleep(60 * time.Second) 
	rules := ReadRule()
	var filename = "./output.html"
	f, err1 := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	defer f.Close()

	menu := parseexcel()
	failed := false
	for index, _ := range menu.Dailys {
		for index2, _ := range menu.Dailys[index].Dishes {
			dish, ok := rules[menu.Dailys[index].Dishes[index2].Name]
			if !ok || dish.Type == "" {
				fmt.Println(menu.Dailys[index].Dishes[index2].Name)
				failed = true
				continue
			}

			menu.Dailys[index].Dishes[index2].Type = dish.Type
			if dish.Url != "" {
				menu.Dailys[index].Dishes[index2].Url = dish.Url
			} else {
				url := GetMenuUrl(dish.Name)
				menu.Dailys[index].Dishes[index2].Url = url
			}
		}
	}
	if failed {
		fmt.Println("将上述菜名写入规则库，谢谢，亲爱的宝贝")
		return
	}
	// 读取templlate文件
	t, err := template.ParseFiles("./layouts/main.tmpl","./layouts/daily.tmpl")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.ExecuteTemplate(f, "daily.tmpl", menu)
	fmt.Println("亲爱的老婆转换完成了，复制ouput.html中的内容到 https://bj.96weixin.com/")
	fmt.Println("可以关闭这个界面了")
}
