package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var urls = make(chan string)
	var wg2 sync.WaitGroup
	var wg sync.WaitGroup

	fileName := "movies.txt"

	listURL := "https://movie.douban.com/top250?start="

	for i := 0; i < 10; i++ {
		wg2.Add(1)
		start := i * 25
		newURL := listURL + strconv.Itoa(start)

		go func() {

			getLIST(newURL, urls)
			defer func() {
				wg2.Done()
			}()
		}()
	}

	fileObject, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for i := 0; i < 250; i++ {
		go func() {
			wg.Add(1)
			url := <-urls

			info := getMOVIE(url)
			// append(infos, info)

			fileObject.Write([]byte(info))

			defer func() {
				wg.Done()
				if err := recover(); err != nil {

					fmt.Println(err)
				}

			}()
		}()
	}
	wg.Wait()
	wg2.Wait()

	fileObject.Close()

	fmt.Println("完成")
}

// getLIST 下载函数
// url 分页列表链接
//
func getLIST(url string, urls chan string) {

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	//函数结束后关闭相关链接
	defer res.Body.Close()

	// ioutil.ReadAll(res.Body)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Panicln(err.Error())
	}

	doc.Find("#content .article .grid_view .item .pic a").Each(func(i int, s *goquery.Selection) {

		movieURL, _ := s.Attr("href")

		urls <- movieURL
	})
}

// getMOVIE 下载函数
// url 下载
// string : 返回值
func getMOVIE(url string) string {

	res, err := http.Get(url)
	if err != nil {
		// return nil
		panic(err)
	}
	fmt.Println(url)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	//三种写法都可以
	// fmt.Println("名称：" + doc.Find(`#content h1`).ChildrenFiltered(`[property="v:itemreviewed"]`).Text())
	// fmt.Println("名称：" + doc.Find(`#content h1 span`).Filter(`[property="v:itemreviewed"]`).Text())

	// doc.Find("#content h1").Each(func(i int, s *goquery.Selection) {

	// 	name := s.ChildrenFiltered(`[property="v:itemreviewed"]`).Text()
	// 	fmt.Println("名称：" + name)
	// })
	str := "名称：" + doc.Find(`#content h1 span`).Filter(`[property="v:itemreviewed"]`).Text() + "\n"

	str += "年份：" + doc.Find("#content .year").Text() + "\n"

	str += "导演：" + doc.Find("#info .attrs a").Eq(0).Text() + "\n"

	pl := ""
	doc.Find("#info .attrs").Eq(1).Each(func(i int, s *goquery.Selection) {
		pl += s.Text()
	})
	str += "编剧:" + pl + "\n"

	zy := ""
	doc.Find("#info .attrs").Eq(2).Each(func(i int, s *goquery.Selection) {
		zy += s.Text()
	})
	str += "主演:" + zy + "\n"
	return str
}
