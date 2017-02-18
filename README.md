# Code Poster

生成代码明信片。支持代码压缩，图片缩放，自动居中。[使用Go编写代码明信片生成器](http://cjting.me/golang/2017-02-18-%E4%BD%BF%E7%94%A8Go%E7%BC%96%E5%86%99%E4%BB%A3%E7%A0%81%E6%98%8E%E4%BF%A1%E7%89%87%E7%94%9F%E6%88%90%E5%99%A8.html)。

![](http://ww1.sinaimg.cn/large/9b85365dgy1fcujv2s1khj20m80l4n55)

## 安装

### 下载

下载相关平台的[二进制程序](http://github.com/fate-lovely/codeposter/releases)，在终端中添加执行权限就可以执行了（Windows不需要）。

```bash
$ chmod +x codeposter_darwin_amd64
$ ./codeposter_darwin_amd64 --help
```

### Go

```bash
go get -u github.com/fate-lovely/codeposter
```

## 参数

```bash
$ codeposter -h
usage: codeposter [<flags>] <source> <image>

Flags:
  -h, --help                Show context-sensitive help (also try --help-long
                            and --help-man).
      --font="Hack"         font family, please use monospace font,
      --fontsize="11.65px"  font size, valid css font size, must corresponding
                            to char width and char height
      --charwidth=7         single character width in pixels
      --charheight=14       single character height in pixels
      --width=800           output poster width in pixels
      --height=760          output poster height in pixels
      --bgcolor="#eee"      background color, valid css color
      --output=canvas       specify output format, [canvas | dom]
      --version             Show application version.

Args:
  <source>  source code path
  <image>   image path
```

- `font`：字体，默认使用`Hack`，务必选择一款等宽字体
- `fontsize`：字体大小，选个一个合适的字体大小，保证对应的字符的宽度和高度是一个整数
- `charwidth`：单个字符宽度，这个需要在浏览器中手动测量
- `charheight`：单个字符高度，这个也需要在浏览器中手动测量
- `width`：最终明信片的宽度，单位是像素，整数
- `height`：最终明星片的高度，单位是像素，整数
- `bgcolor`：背景颜色
- `output`：输出格式，目前支持`dom`和`canvas`。注意，dom格式将每个字符渲染为一个div，十分消耗性能，默认格式为canvas。

## 示例

进入`example`文件夹。

### Gopher

```bash
codeposter jquery.min.js go.png > go.html
```

![](http://ww1.sinaimg.cn/large/9b85365dgy1fcujvjjomgj20m70l37ce)

### Heart

```bash
codeposter jquery.min.js heart.png > heart.html
```

![](http://ww1.sinaimg.cn/large/9b85365dgy1fcujw0juyxj20m70l3wmc)

### Diamond

```bash
codeposter jquery.min.js diamond.png > diamond.html
```

![](http://ww1.sinaimg.cn/large/9b85365dgy1fcujw08zu4j20m70l47ck)