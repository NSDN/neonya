# 创建漫画主题 API（草案）

## 说明

提交主题执行此 API。

## 请求

|方法|URL|
|-|-|
|PSOT|/topic/commic|

### 请求参数

|物理名|逻辑名|类型|必须|
|-|-|-|-|
|plateRouteName|版块路由名|string|true|
|topicTitle|主题标题|string|true|
|pages|图片列表|Commic.Page[]|true|

#### Commic.Page

|物理名|逻辑名|类型|必须|说明|
|-|-|-|-|-|
|id|主键|number|true|数组中的元素下标|
|path|图标地址|string|true|初始为空或为用户图床中的地址|
|name|图片名称|string|true||
|content|图片文件本体|Blob|false|用户上传文件时使用此属性|

## Response

- 正常： 状态码
- 异常： 状态码 + 错误信息
