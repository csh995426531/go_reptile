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
	var wg1 sync.WaitGroup
	var wg sync.WaitGroup
	file := "movies.txt"
	// var infos = [250]string{}

	fl, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	fl.Write([]byte("start\n"))
	fl.Close()

	listURL := "https://movie.douban.com/top250?start="

	for i := 0; i < 10; i++ {

		start := i * 25
		newURL := listURL + strconv.Itoa(start)
		wg1.Add(1)
		go func() {

			getLIST(newURL, urls)

			defer func() {
				wg1.Done()
			}()
		}()

	}
	wg1.Wait()

	for i := 0; i < 250; i++ {
		go func() {
			wg.Add(1)
			url := <-urls

			info := getMOVIE(url)
			// append(infos, info)
			// fmt.Println(info)
			fl, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return
			}

			fl.Write([]byte(info))
			defer func() {
				wg.Done()
				if err := recover(); err != nil {
					// fl.Write([]byte(p))
					fmt.Println(err)
				}
				fl.Close()

			}()
		}()
	}
	wg.Wait()
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
		fmt.Println(urls)
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
