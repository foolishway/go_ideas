## crawler.go

按页抓取博客数据并保存到本地

usage

./crawler -s https://zhangwei412827.github.io -b .posts-list-item-title -p .page-link

本地目录生成结构如下：

![](https://zhangwei412827.github.io/images/crawler.png)

## errer_handler.go

该例子来自于Rob Pike的一篇博客“Errors are values”，地址：
<https://blog.golang.org/errors-are-values/>

## file_server.go

FTP demo

通过设置-r指定目录，通过-p指定监听端口，即可对该目录设置FTP服务，绑定在-p设置的端口上

## goroutine_calc.go

goroutine demo
