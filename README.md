# TFLanHttpDesktop
Transfer Files from LAN Http Desktop, 用于局域网内指定文件生成二维码或链接提供给三方设备用局域网http协议下载文件，三方设备也可以上传文件，桌面应用程序，跨平台

#### 核心功能

简单，安全，稳定，实用

- 在同内网下宿主机上设定文件下载和上传路径
- 在同wifi下移动设备扫码下载和上传
- 在同内网下共享备忘录，在不同设备上共享文本信息
- 设定密码，一次性下载等功能

#### todo
- 发布v0.1
- 优化现有功能
- 发布v0.2
- 跨平台测试 windows linux 
- 修改兼容bug
- 发布v0.3
- 性能测试，功能测试
- 修改bug,提升稳定性
- 发布v0.4

## 设计说明
1. 该应用采用http协议传输文件，内网中使用，宿主机需要开放防火墙
2. 为了安全，每次启动应用程序随机http服务端口
3. 为了安全，只能设置一个下载文件和一个上传路径
4. 设置密码没有任何限制，密码长度为0则没有密码

## 里程碑规划
v0.1 基础功能
v0.2 优化现有功能
v0.3 跨系统测试，提升兼容性
v0.4 修改bug,稳定版本,并申请应用市场
v0.5 新需求,待续...


## 体验与测试
- [优化]当前停留页面备忘录被修改需要跟新ui
- [优化]备忘录大文本会卡,方案是最多1w字,限制大文本
- [优化]备忘录正文需要有个标题
- [优化]关于能link到项目地址
- [新需求]还是需要有个分享的短链，因为复制的链接是md5编码，远程设备输入太复杂

## fyne 2.6 局限性
- 弹出层的聚焦会夺去输入框的聚焦，无法做到输入弹出联动框
- 底层私有了语言，外部无法直接操作选择语言，语言跟系统一致
- dialog createInformationDialog 私有，无法做到更灵活的定制化
- 系统托盘无法将缩小进行打开
- github.com/go-gl/gl go mod 依赖不兼容导致拉包失败
- 新老版本方法差距太大，市面上资料和文档老版本偏多
- windows上交叉编译linux,darwin环境不好搭建

## build
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
/d/go/bin/fyne.exe package -os windows -icon ./icon.png -app-id "TFLanHttpDesktop.2025.0826" -app-version 0.1.1
CertUtil -hashfile "TFLanHttpDesktop-v0.1.1-windows-arm64.zip" SHA256
TFLanHttpDesktop-v0.1.1-windows-amd64.zip.sha256


export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1
/d/go/bin/fyne.exe package -os linux -icon ./icon.png -app-id TFLanHttpDesktop.2025.0826 -app-version 0.1.1

set GOARCH=arm64
/d/go/bin/fyne.exe package -os windows -icon ./icon.png -app-id TFLanHttpDesktop.2025.0826 -app-version 0.1.1
/d/go/bin/fyne.exe package -os linux -icon ./icon.png -app-id TFLanHttpDesktop.2025.0826 -app-version 0.1.1
TFLanHttpDesktop-v0.1.1-windows-arm64.zip TFLanHttpDesktop-v0.1.1-windows-arm64.zip.sha256

需要在苹果系统上打包
fyne package -os darwin -icon ./icon.png -app-id TFLanHttpDesktop.2025.0826 -app-version 0.1.1
fyne package -os darwin -icon ./icon.png -app-id TFLanHttpDesktop.2025.0826 -app-version 0.1.1