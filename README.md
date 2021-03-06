# Andex

一个将阿里网盘作为云盘分享的工具

## 配置文件

路径是运行目录下 `./config.json`

```json
{
  "refresh_token": "REFRESH_TOKEN",
  "root_path": "/publicDir" 
}
```
 - `refresh_token` 阿里云盘 refresh token
 
 - `root_path` 云盘内文件夹的路径，如 `/public/Andex`， 该路径将会映射到 Andex 首页路径, 默认为 `/`, 即网盘根目录  

## Feature
 - 根目录文件夹设置
 - 不直接暴露文件夹 id
 - 阿里云盘下载直链获取
 - 网盘 README.md 展示

## TODO
 - 提供后台管理/设置/初始化页面
    - README.md 编辑
    - 文件/文件夹下载加密(密码设置)
- 文件/文件夹下载次数统计
 - ip 防护/黑名单防护
 
 ## 运行效果

<img src="static/screenshot/screenshot1.png" style="width: 80%;margin-left: 10%" alt="图片1">

<img src="static/screenshot/screenshot2.png" style="width: 80%;margin-left: 10%" alt="图片2">

<img src="static/screenshot/screenshot3.png" style="width: 80%;margin-left: 10%" alt="图片3">
