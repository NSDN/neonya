# 概述

## 技术栈

|技术栈|功能|
|-|-|
|vite|项目工具|
|vue3|JS 框架|
|TypeScript|编程语言|
|vue-router|路由|
|pinia|store|
|axios|HTML 请求|
|Naive UI（暂未使用）|UI 框架|
|tiptap|富文本编缉器|
|marked|Markdown 编缉器|
|crypto-js|hash 预处理工具|

## 书写规范

- 使用 [JavaScript Standard Style](https://standardjs.com/readme-zhcn.html)
- 文件名由小写字母及下划线、短横线所组成。短横线表示词组的连接；下划线表示所属。
- 变量名尽量将单词拼写完整，如果一个词组有常用缩写可以使用（如 HyperText Markup Language > HTML）。

## 目录结构

```shell
src
├─assets
│  ├─css
│  └─img
│      └─icon
├─components
│  └─button
├─config
├─constant
│  └─enums
├─context
├─data
├─features
│  ├─authorization
│  │  ├─components
│  │  ├─hooks
│  │  ├─services
│  │  └─types
│  └─plate
│      ├─components
│      ├─hooks
│      ├─services
│      └─types
├─hooks
├─layouts
├─lib
│  ├─http
│  ├─router
│  └─store
├─utils
└─views
    ├─computer
    └─phone
```

|目录|描述|备注|
|-|-|-|
|assets|放置样式、图片等静态文件||
|components|放置共通组件|可以按组件的种类创建子目录；添加新内容时需要更新开发文档|
|config|放置设置文件||
|constant|放置常量文件，包括枚举||
|context|放置上下文文件||
|data|放置静态数据||
|**features**|放置功能模块|核心目录；如果有不明白的地方务必询问冰龍|
|hooks|放置全局勾子||
|layouts|放置布局组件||
|lib|放置第三方库||
|utils|放置工具函数|添加新内容时需要更新开发文档|
|views|放置画面本体组件|目前的想法是桌面端移动端分别做一套|
