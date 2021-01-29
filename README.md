# RSSProxy

## 目的

实现了一个能够设置代理，中转 RSS 源的中间件，解决有一些 RSS 源需要代理才能访问的问题。

当使用 TTRSS 来订阅的时候，虽然提供了全局的代理设置，但是如果你没有搭建软路由等方法去做流量的转发，其实还是很麻烦的。同时，TTRSS 内置的插件 Options per feed ，说是可以为单独的 RSS 源走单独的代理，但是经过测试没有成功，则有了本中间件。

还有一点，RSSHUB 提供的某些网站的订阅是需要登录的，但是近期已经无法正常登录了，对应提供的方法也无效，所以还是得去 RSSHUB 上直接订阅，那么 RSSHUB 已经不能正常访问，所以你懂的。

## 特性

* 支持设置代理（强制）
* 支持 instagram 订阅用户的 RSS
* 支持 twitter 订阅用户的 RSS
* 支持中转 rsshub （完全转发）

## 如何使用

### 编译主程序

```bash
go build .
```

### 修改配置文件

注意，启动后再修改配置文件，会自动触发重载逻辑，便于更新关注的列表，或者是用户名密码什么的。

修改 config.yaml.sample 为 config.yaml，确保与主程序在同一级目录

```yaml
ListenPort: 1200
HttpProxy: http://127.0.0.1:10809
EveryTime: 30m

RSSInfos:
  pixiv_week_r18: https://rsshub.app/pixiv/ranking/week_r18
  pixiv_day_male: https://rsshub.app/pixiv/ranking/day_male
  
 IGInfo:
  UserName: username
  PassWord: password
  FeedMaxItems: 30
  InstagramUsers:
    - fjamie013
    
 TwiInfo:
  FeedMaxItems: 200
  ExcludeReplies: true
  PhotoOnly: true
  TwitterUsers:
    - baby_eiss
```

看上面示例哈。做了一个缓存 EveryTime 是 cron 的定时刷新时间，默认 30min。

#### ListenPort

本程序的监听端口

#### HttpProxy

代理服务器，需要同时设置 IP 和 PORT，必须设置，不然还用这个工具干啥···

#### RSSInfos

看上面的示例哈

> http://127.0.0.1:1200/rss?key=pixiv_day_male		-- https://rsshub.app/pixiv/ranking/day_male
>
> http://127.0.0.1:1200/rss?key=pixiv_month			  -- https://rsshub.app/pixiv/ranking/month

#### IGInfo

看上面的示例哈

> http://127.0.0.1:1200/rss/instagram?key=fjamie013			  -- https://www.instagram.com/fjamie013

#### TwiInfo

ExcludeReplies: 跳过回复

看上面的示例哈

> http://127.0.0.1:1200/rss/twitter?key=baby_eiss			  -- https://twitter.com/baby_eiss

### Docker 部署

参考项目内的 docker-compose 即可。需要注意的是，如果使用 docker 来部署，本程序的 config 文件中的端口号记得不要改，默认就是 1200，然后容器外面是啥你再定。

```yaml
version: '3.0'

services:
  rssproxy:
    image: allanpk716/RSSProxy:latest
    volumes:
      - /mnt/user/appdata/rssproxy/config.yaml:/app/config.yaml
    ports:
      - 1201:1200
    restart: always
```

## TTRSS 可能遇到的问题

如果你用的 TTRSS 是2020年12月左右的（个人最近才更新过一次），你可能会发现 TTRSS 无法订阅 80、433 以外的 RSS 源。启动 TTRSS docker 的时候添加以下环境即可：

```yaml
    environment:
      - ALLOW_PORTS=1200
```

## 致谢

感谢下列项目的 Instagram 的相关代码，抄的很开心

* [instafeed](https://github.com/falzm/instafeed)
* [goinsta](https://github.com/ahmdrz/goinsta)
* [twitter2rss](https://github.com/n0madic/twitter2rss)
