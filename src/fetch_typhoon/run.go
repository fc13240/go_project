/*抓取中央台网网的数据*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"time"
)

var dirCurrent = filepath.Dir(os.Args[0])

var dirSave = filepath.Join(dirCurrent, "data")

var urlTyhoon = "http://typhoon.nmc.cn/weatherservice/typhoon/jsons/"

type typhoon struct {
}
type typhoonList struct {
	TyphoonList [][]interface{} `json:"typhoonList"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func setDefer() {
	defer func() {
		if e, ok := recover().(error); ok {
			log.Printf("WARN: panic in %v", e)
			log.Println(string(debug.Stack()))
		}
	}()
}

/*保存文件内容*/
func saveFile(path, content string) {
	path = filepath.Join(dirSave, path)

	dir, _ := filepath.Split(path)
	if _, err := os.Stat(dir); err != nil {
		os.MkdirAll(dir, 0777)
	}
	ioutil.WriteFile(path, []byte(content), 0644)

	fmt.Printf("\n%s saved!\n", path)
}

/*抓取一个URL的内容*/
func fetch(url string, fn func(string)) {
	fmt.Printf("\nfetching: %s\n", url)
	resp, err := http.Get(url)

	setDefer()
	defer resp.Body.Close()

	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	result := string(body)
	// fmt.Println(result)

	// 处理jsonp格式的响应，并防止内容中出现括号
	reg := regexp.MustCompile(`^[^(]+\(+([\s\S]+)\)+$`)

	arr := reg.FindStringSubmatch(result)

	if len(arr) == 2 {
		jsonStr := arr[1]

		fn(jsonStr)
	}
}

/*得到台风路径列表*/
func getTyphoonList() {
	for i := time.Now().Year(); i >= 1949; i-- {
		year := strconv.Itoa(i)
		fileName := "list_" + year
		fetch(urlTyhoon+fileName+"?"+strconv.Itoa(int(time.Now().Unix())), func(content string) {
			saveFile(fileName, content)
			var result typhoonList
			err := json.Unmarshal([]byte(content), &result)
			if err == nil {
				for _, typhoon := range result.TyphoonList {
					code := strconv.Itoa(int(typhoon[0].(float64)))

					getTyphoonDetail(year, code)
				}
			}
		})
	}
}

/*得到单个台风的详细数据*/
func getTyphoonDetail(year, code string) {
	fileName := "view_" + code
	fetch(urlTyhoon+fileName+"?"+strconv.Itoa(int(time.Now().Unix())), func(content string) {
		// fileName := filepath.Join(year, fileName)
		saveFile(fileName, content)
	})
}
func main() {
	if len(os.Args) > 1 {
		dirSave = os.Args[1]
	}
	fmt.Printf("dirCurrent = %s, dirSave = %s\n", dirCurrent, dirSave)
	getTyphoonList()
}
