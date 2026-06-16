# 喵玉汉化馆帖文列表获取 API

## 请求

|方法|URL|
|-|-|
|GET|/plate/localization|

### 响应

#### Status: 200 OK

|物理名|逻辑名|类型|备注|
|-|-|-|-|
|id|编号|number||
|title|标题|string||
|thumbnail|封面|string|图床中的地址|

```json
[{
  "id": 0,
  "title": "sample",
  "thumbnail": "http://example.com/example.jpg"
}]
```
