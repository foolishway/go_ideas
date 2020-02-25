package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	site, blogSelector, pagSelector *string
	rootPath                        string
)

func main() {
	site = flag.String("s", "", "the site to fetch")
	blogSelector = flag.String("b", "", "the blog link css selector")
	pagSelector = flag.String("p", "", "the pagination selector")

	flag.Parse()

	if *site == "" || *blogSelector == "" || *pagSelector == "" {
		flag.Usage()
		return
	}
	if (*site)[len(*site)-1:] == "/" {
		*site = (*site)[:len(*site)-1]
	}
	rootPathUrl, _ := url.Parse(*site)
	rootPath = rootPathUrl.Host

	//判断根目录是否存在，存在则删除
	if _, err := os.Stat("./" + rootPath); !os.IsNotExist(err) {
		os.RemoveAll(rootPath)
	}

	//create root dic
	err := Mkdir(rootPath)
	if err != nil {
		panic("Create root dictory error.")
	}

	var pages []string

	htmlBytes, err := getHtmlBytes(*site)
	if err != nil {
		panic(fmt.Sprintf("Read %v error %v.", *site, err))
	}
	reader := bytes.NewReader(htmlBytes)
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		panic(err)
	}

	doc.Find(*pagSelector).Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if _, isExist := Find(pages, href); !isExist {
			pages = append(pages, href)
		}
	})
	fmt.Printf("%#v\n", pages)
	var wg sync.WaitGroup
	for i := 0; i < len(pages); i++ {
		//todo 加载每页数据，找到blogs href，获取html
		wg.Add(1)
		go fetchPage(pages[i], i, &wg)
	}
	wg.Wait()
}
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
func fetchPage(pageUrl string, pageNum int, wg *sync.WaitGroup) {
	log.Printf("开始下载第%d页数据~", pageNum)
	err := Mkdir(rootPath + "/" + strconv.Itoa(pageNum))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Fetch page %s\n", *site+pageUrl)
	htmlBytes, err := getHtmlBytes(*site + pageUrl)
	if err != nil {
		log.Fatalf("Fetch %s error:%v", pageUrl, err)
	}
	reader := bytes.NewReader(htmlBytes)
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		log.Fatal(err)
	}
	var childWG sync.WaitGroup

	blogs := doc.Find(*blogSelector)
	childWG.Add(blogs.Length())
	blogs.Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		var path string = href
		//if filepath.IsAbs(href) {
		//	path = href
		//} else {
		//	path = *site + href
		//}

		go func() {
			defer childWG.Done()
			html, err := getHtmlBytes(path)
			if err != nil {
				panic(err)
			}
			var fileName string

			sp := strings.Split(path, "/")
			//_, fileName := filepath.Split(path)
			if path[len(path)-1:] == "/" {
				fileName = sp[len(sp)-2]
			} else {
				fileName = sp[len(sp)-1]
			}
			log.Printf("开始下载%s~\n", fileName)
			fullName := fmt.Sprintf("%s/%s/%s.html", rootPath, strconv.Itoa(pageNum), fileName)
			log.Printf("%s is downloading.", fullName)
			err = SaveLocalFile(html, fullName)

			if err != nil {
				panic(err)
			}
		}()
	})
	childWG.Wait()
	wg.Done()
}

func getHtmlBytes(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	return bytes, err
}

func SaveLocalFile(body []byte, fileName string) error {
	err := ioutil.WriteFile(fileName, body, 0777)
	return err
}

func Mkdir(dirName string) error {
	if dirName == "" {
		return errors.New("Directory name is required.")
	}
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		log.Printf("Create dir error: %s\n", err)
		return err
	}
	err = os.Chmod(dirName, 0777)

	if err != nil {
		log.Printf("Change mod error: %s\n", err)
		return err
	}
	return nil
}
