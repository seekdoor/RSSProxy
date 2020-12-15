# RSSProxy

## 目的

实现了一个能够设置代理，中转 RSS 源的中间件，解决有一些 RSS 源需要代理才能访问的问题。

当使用 TTRSS 来订阅的时候，虽然提供了全局的代理设置，但是如果你没有搭建软路由等方法去做流量的转发，其实还是很麻烦的。同时，TTRSS 内置的插件 Options per feed ，说是可以为单独的 RSS 源走单独的代理，但是经过测试没有成功，则有了本中间件。

还有一点，RSSHUB 提供的某些网站的订阅是需要登录的，但是近期已经无法正常登录了，对应提供的方法也无效，所以还是得去 RSSHUB 上直接订阅，那么 RSSHUB 已经不能正常访问，所以你懂的。

## 如何使用

### 编译主程序

```bash
go build Main.go
```

### 修改配置文件

修改 config.yaml.sample 为 config.yaml，确保与主程序在同一级目录

```yaml
ListenPort: 1200
HttpProxy: http://127.0.0.1:10809
RSSInfos:
  pixiv_month:
    RSSUrl: https://rsshub.app/pixiv/ranking/month
  pixiv_day_male:
    RSSUrl: https://rsshub.app/pixiv/ranking/day_male
  instagram_fjamie013:
    RSSUrl: https://rsshub.app/instagram/user/fjamie013
```

#### ListenPort

本程序的监听端口

#### HttpProxy

代理服务器，需要同时设置 IP 和 PORT，必须设置，不然还用这个工具干啥···

#### RSSInfos

这里支持多个需要中转的 RSS 源，有一定 Key 命名要求，以上述配置举例。

有以下三个 Key，且会配合 RSSHUB 来使用，所以会以 "\_" 来分割，且只分割第一个 "\_"，slice[0] 是目标源网站的描述， slice[1] 则是具体的订阅路由，然后自动构建路由。

* pixiv_month
  *  slice[0] = pixiv , slice[1] = month
  * 订阅地址：http://127.0.0.1:1200/pixiv/month
* pixiv_day_male
  *  slice[0] = pixiv , slice[1] = day_male
  * 订阅地址：http://127.0.0.1:1200/pixiv/day_male
* instagram_fjamie013
  *  slice[0] = pixiv , slice[1] = day_male
  * 订阅地址：http://127.0.0.1:1200/instagram/fjamie013

### 使用其他 RSS 软件订阅中转后的 RSS 源

如上所述，直接运行主程序，那么对应的三个订阅地址则是：

* http://127.0.0.1:1200/pixiv/month
* http://127.0.0.1:1200/pixiv/day_male
* http://127.0.0.1:1200/instagram/fjamie013