# 论坛版块列表获取 API

## 请求

|方法|URL|
|-|-|
|GET|/plates|

### 响应

#### Status: 200 OK

|物理名|逻辑名|类型|备注|
|-|-|-|-|
|id|编号|int||
|name|显示名|string||
|routeName|路由名|string||
|nameColor|显示名颜色|string|同 CSS|
|background|背景图片|string|图床地址|
|pageType|画面类型|PageType|PageType: 枚举 ["commic", "article"]|

```json
[{
  id:         0,
  name:       "喵玉汉化馆",
  routeName:  "PlateLocalization",
  nameColor:  "#fff",
  background: "https://static-event.benghuai.com/new_mihoyo_homepage/images/download/cg/origin/2020-10-22.jpg",
  pageType: "commic"
}]
```
