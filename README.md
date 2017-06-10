# json2appcache

Tool for generating cache manifest for the released Egret app via its
`manifest.json`, basically a monkey patch for Egret's
[resourcemanager](https://github.com/egret-labs/resourcemanager).

## CLI usage

```bash
$ json2appcache manifest.json > example.appcache
```

## 关于 Egret 4.1.0 资源管理 (`resouremanager`) 与项目发布

目前个人使用的 NPM 模块 ID 为 `egret-resource-manager@4.0.4-25`,
资源发布流程如官方描述一致

```bash
$ egret build   # 构建项目
$ res build     # 构建资源
$ egret publish # 发布项目
$ res publish . bin-release/web/<版本号>    # 发布资源
```

### 官方的 `manifest.json` 怎么使用?

- 按照官方进度, 目前未知 :(
- 目前的解决方案为使用 application cache 进行缓存, 在根目录放置 `*.appcache`
文件配置需要缓存的资源, 而需要缓存的资源都在官方的 `manifest.json` 中列出
- 使用此工具生成一个 `example.appcache` 文件, 放置在根目录下
(工具默认将文件打印到标准输出中, 需要重定向到文件)
- 在代码中要注意同步客户端的 *第一次缓存*, 因为客户端脚本的 `src` 及相关样式表的
`href` 请求的是没有带散列值后缀的文件, 所以缓存没有下载完成时会 404.

```js
function runEgret (e) {
  if (e.type) {
    switch (e.type) {
      case 'cached':        // 第一次下载完毕
      case 'noupdate':      // 不需要更新缓存
      case 'updateready':   // 更新缓存完毕
        egret.runEgret({renderMode: "webgl", audioType: 0, retina: true})
        break
    }
  }
}

window.applicationCache.addEventListener('cached', runEgret, false)
window.applicationCache.addEventListener('noupdate', runEgret, false)
window.applicationCache.addEventListener('updateready', runEgret, false)
```

- 配置服务器的时候, 虽然在我的 Chrome 59 上默认配置是可行的,
但出于客户端兼容性考虑, 应该给服务器配置相应的 MIME types, 例如在 Nginx 的
`mime.types` 中加入字段

```
text/cache-manifest                   appcache;
```

- 对于一并兼容其他格式的 Manifest 文件, 可参考
[server-configs-nginx](https://github.com/h5bp/server-configs-nginx/blob/master/mime.types), 相应字段如下

```
  # Manifest files

    application/manifest+json             webmanifest;
    application/x-web-app-manifest+json   webapp;
    text/cache-manifest                   appcache;
```

- 需要注意的是, 在 Nginx 中服务 appcache 这个静态文件是默认没有 `Cache-Control`
头的, 需要服务器缓存控制的话, 只需要更改静态文件的超时规则即可, 可以参考
[Nginx Caching](https://serversforhackers.com/nginx-caching)

```
location ~* \.(?:manifest|appcache)$ {
  expires -1;   # 服务器不缓存 (推荐)
  # expires 1h; # 缓存, 1 小时后超时
}
```

- 注: 根据
[Using the application cache](https://developer.mozilla.org/en-US/docs/Web/HTML/Using_the_application_cache),
该标准已经是 *deprecated* 状态, 但是
[兼容性](http://caniuse.com/#feat=offline-apps) 算是相当的好, 而 Service Worker
兼容性仍不稳定, 即便它有更好的版本管理和缓存操作接口

### 缓存怎么做到的?

查看 `example.appcache` 文件可知, 文件散列值已经基本可以作为该 appcache
文件的唯一标识, 并且在 `CACHE:` 表项中直接表明了带有散列值的缓存文件.
客户端在请求没有散列值后缀的文件时, 使用 `FALLBACK:`
表项可以让客户端重定向到其相应的文件, 实现缓存读取.

- appcache 文件例子:

```
CACHE MANIFEST
# v42 - 版本号已经不必要

CACHE:
/resource/config_d34db33f.json

NETWORK:
*

FALLBACK:
/resource/config.json /resource/config_d34db33f.json
```

## License

MIT
