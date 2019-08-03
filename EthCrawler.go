package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)



var (

	urls = make(chan string, 100)
	sourceCode = make(chan string, 100)
	killsign = make(chan string)
	count = 0
	file *os.File
	pages = 1 //爬几页
)

func main() {
	fp, err := os.Create("./demo.txt")
	if err != nil{
		fmt.Println(err)
	}
	file = fp
	defer file.Close()
	go Start(pages)
	go getSourceCode()
	go showCode()
	ks := <- killsign
	fmt.Println(ks)
	fmt.Println(count)

}

func showCode(){
	for k := range sourceCode{
		fmt.Println("ETH sourceCode:")
		fmt.Println(k)
		file.WriteString(k)
		count++
	}
	killsign <- "Crawler Stop"
}

func getSourceCode(){
	for u := range urls{
		getSource(u)
	}
	close(sourceCode)
}
func getSource(url string){
	fmt.Println("Geting ",url,"SourceCode")
	resp, err := http.Get(url)
	if err != nil{
		fmt.Println("getSource: ", err)
		close(sourceCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("getSource: ", err)
		close(sourceCode)
	}

	reg := regexp.MustCompile(`top: 5px;'>(?s:(.*?))</pre><br><script>`)
	result := reg.FindAllString(string(body),-1)
	for i:=0;i<len(result);i++{
		result[i] = strings.TrimLeft(result[i],"top: 5px;'>")
		result[i] = strings.TrimRight(result[i],"</pre><br><script>")
		sourceCode <- url+ "\n"+ result[i] + "\n"
	}
}

func Start(pages int){
	url := "http://etherscan.io/contractsVerified"
	for i := 1; i <= pages; i++{
		if i != 1{
			firstStart(url + "/" + strconv.Itoa(i))
		}else{
			firstStart(url)
		}
	}
	close(urls)
}

func firstStart(url string) {


	resp, err := http.Get(url)
	if err != nil{
		fmt.Println("firstStart: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("firstStart: ", err)
	}
	reg := regexp.MustCompile(`/address/[a-zA-z=0-9-\s]*`)

	for _, d := range reg.FindAllString(string(body), -1) {
		fmt.Println("地址收集： ", d)
		urls <- "http://etherscan.io" + d
	}
	fmt.Println("首次收集网络地址：" ,len(urls))

}
func checkRegexp(cont string, reg string, style int) (result interface{}) {
	check := regexp.MustCompile(reg)
	switch style {
	case 0:
		result = check.FindString(cont)
	case 1:
		result = check.FindAllString(cont, -1)
	default:
		result = check.FindAll([]byte(cont), -1)
	}
	return
}



func checkFile(dir string, file string) os.FileInfo {
	list, _ := ioutil.ReadDir(dir)
	for _, info := range list {
		if info.Name() == file {
			return info
		}
	}
	return list[0]
}


