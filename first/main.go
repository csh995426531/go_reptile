package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func hello(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "hello ,this is from HelloServer func")
}

func main() {
	var urls []string

	listURL := "https://movie.douban.com/top250?start="

	for i := 0; i < 10; i++ {

		start := i * 25
		newURL := listURL + strconv.Itoa(start)

		urls = getLIST(newURL)

		for _, url := range urls {
			getMOVIE(url)
		}
	}
}

// getLIST 下载函数
// url 下载
// byte : 返回值
func getLIST(url string) []string {
	var urls []string

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	defer res.Body.Close()

	// ioutil.ReadAll(res.Body)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Panicln(err.Error())
	}

	doc.Find("#content .article .grid_view .item .pic a").Each(func(i int, s *goquery.Selection) {
		// fmt.Printf("%v", s)
		movieURL, _ := s.Attr("href")
		urls = append(urls, movieURL)
	})
	return urls
}

// getMOVIE 下载函数
// url 下载
// byte : 返回值
func getMOVIE(url string) {
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		// return nil
		panic(err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	//三种写法都可以
	fmt.Println("名称：" + doc.Find(`#content h1`).ChildrenFiltered(`[property="v:itemreviewed"]`).Text())
	fmt.Println("名称：" + doc.Find(`#content h1 span`).Filter(`[property="v:itemreviewed"]`).Text())

	doc.Find("#content h1").Each(func(i int, s *goquery.Selection) {

		name := s.ChildrenFiltered(`[property="v:itemreviewed"]`).Text()
		fmt.Println("名称：" + name)
	})

	fmt.Println("年份：" + doc.Find("#content h1 .year").Eq(1).Text())

	fmt.Println("导演：" + doc.Find("#info .attrs a").Eq(1).Text())

	//这里不太懂
	pl := ""
	doc.Find("#info span:nth-child(3) span.attrs").Each(func(i int, s *goquery.Selection) {
		pl += s.Text()
	})
	fmt.Println("编剧:" + pl)
}
