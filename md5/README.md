> 在过去的几年里，wasm的话题那真是从早上聊到晚上，可以说处于异常兴奋的状态，但是几年过去了，它慢慢的被大多数人们忘记，原因比较简单——落地难

> 今天就wasm能给js加多少分这个问题，做一个小型的讨论，今天的专注点是，前端js获取一个文件的md5值，也就是上传文件时所需要的秒传功能的核心

简单来说，文件上传秒传不仅仅是网盘公司的专属，平时我们上传文件给后端也是很常用的，前端通过对目标文件md5计算后与后端进行对比，如果已经上传过，则直接返回已有地址，这样，大大节省了服务器空间。基本思路如下：
* 前端input type="file"获取文件
* 通过md5工具库进行计算，得到md5值
* 请求接口，后端判断此md5是否已经在数据库里
* 如果在数据库里，则直接告诉前端，已存在（秒传）

### 今日重点
今天的重点是如何快速获取一个文件的md5值，这里就涉及到小文件，大文件的问题了。所以，我将以下面文件体积为例来测试js与wasm对文件md5计算的速度对比。wasm我使用golang进行开发，因为golang打包成wasm会把运行时也加进去，所以，打包的结果2.2M，我们暂时忽略这个体积，因为如果能落地，那么换成rust，换成c++都不是难事，如果不能落地，那么，golang不行，c++也照样不样。

### 准备工作
通过ffmeg 从一个2G+的文件上截取不同体积的文件，用于测试。
```go
ffmpeg -i /path/sourch.mp4  -fs 1M -c:v copy -c:a copy /path/1M.mp4
ffmpeg -i /path/sourch.mp4  -fs 5M -c:v copy -c:a copy /path/5M.mp4
ffmpeg -i /path/sourch.mp4  -fs 20M -c:v copy -c:a copy /path/20M.mp4
ffmpeg -i /path/sourch.mp4  -fs 50M -c:v copy -c:a copy /path/50M.mp4
ffmpeg -i /path/sourch.mp4  -fs 100M -c:v copy -c:a copy /path/100M.mp4
ffmpeg -i /path/sourch.mp4  -fs 200M -c:v copy -c:a copy /path/200M.mp4
ffmpeg -i /path/sourch.mp4  -fs 400M -c:v copy -c:a copy /path/400M.mp4
ffmpeg -i /path/sourch.mp4  -fs 600M -c:v copy -c:a copy /path/500M.mp4
ffmpeg -i /path/sourch.mp4  -fs 800M -c:v copy -c:a copy /path/800M.mp4
ffmpeg -i /path/sourch.mp4  -fs 900M -c:v copy -c:a copy /path/900M.mp4
ffmpeg -i /path/sourch.mp4  -fs 1024M -c:v copy -c:a copy /path/1024M.mp4
ffmpeg -i /path/sourch.mp4  -fs 1280M -c:v copy -c:a copy /path/1280M.mp4
ffmpeg -i /path/sourch.mp4  -fs 1536M -c:v copy -c:a copy /path/1536M.mp4
ffmpeg -i /path/sourch.mp4  -fs 1792M -c:v copy -c:a copy /path/1792M.mp4
ffmpeg -i /path/sourch.mp4  -fs 2048M -c:v copy -c:a copy /path/2048M.mp4
```
### 测试代码
#### 纯js测试代码
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>文件md5</title>
  <script src="./SparkMD5.js"></script>
</head>
<body>
  <input id="file" type="file" />
  <script>
    document.querySelector('#file').addEventListener('change', e => {
      let startTime = Date.now()
      const file = e.target.files[0];
      const fileReader = new FileReader()
      console.log('size', file.size / 1024 / 1024 / 1024, "G")
      fileReader.onprogress = e => {
        console.log(`${Math.floor((e.loaded / e.total) * 100)}%`)
      }
      let usedTime = 0
      const md5 = new SparkMD5();
      fileReader.readAsBinaryString(file);
      fileReader.onload = e => {
        md5.appendBinary(e.target.result);
        const md5Str = md5.end()
        usedTime += Date.now() - startTime
        console.log('usedTime', usedTime, 'ms')
        console.log('md5', md5Str)
      }
    });
  </script>
</body>
</html>
```
#### wasm（go）源码
请参考：

[https://github.com/butoften/wasm](https://github.com/butoften/wasm)

#### js+wasm测试代码
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>文件md5</title>
  <script src="./wasm_exec.js"></script>
</head>
<body>
  <script>
    function handleSayHello(message) {
      console.lof('str from go', message)
    }
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch('wasm/md5.wasm'), go.importObject)
      .then(res => {
        go.run(res.instance);
      });
  </script>
  <input id="file" type="file" />
  <script>
    document.querySelector('#file').addEventListener('change', e => {
      let startTime = Date.now()
      const file = e.target.files[0];
      const fileReader = new FileReader()
      console.log('size', file.size / 1024 / 1024 / 1024, "G")
      fileReader.onprogress = e => {
        console.log(`${Math.floor((e.loaded / e.total) * 100)}%`)
      }
      let usedTime = 0
      fileReader.readAsArrayBuffer(file);
      fileReader.onload = e => {
        const bytes = new Uint8Array(e.target.result)
        wasmMd5Add(bytes)
        const md5Hash = wasmMd5End()
        usedTime += Date.now() - startTime
        console.log('usedTime', usedTime, 'ms')
        console.log('md5', md5Hash)
      }
    });
  </script>
</body>
</html>
```
### 测试条件
* 从FileReader开始读取算起到md5计算结束，因为现实中，我们需要做loading条动画比例
* mac 2.7 GHz 双核Intel Core i5
* mac 8 GB 1867 MHz DDR3
### 测试目标
#### chrome （版本：103.0.5060.114）
 * 2048M 测试5次分别用时：
 * 如何分段计算，每段使用512M
     | 序号 | 纯js |纯js分段 | js+wasm |js+wasm分段 | 
     | :-: | :-:| :-: |:-:| :-:|
     | 1 | 37477 ms | 25638 ms | 31680 ms | 22898 ms |
     | 2 | 32926 ms | 28088 ms | 32516 ms | 25168 ms |
     | 3 | 33413 ms | 31412 ms | 33424 ms | 20547 ms |
     | 4 | 35054 ms | 35821 ms | 33906 ms | 23130 ms |
     | 5 | 35986 ms | 36895 ms | 29014 ms | 22011 ms |
 * 1792M 测试5次分别用时：
     | 序号 | 纯js |纯js分段 | js+wasm |js+wasm分段 | 
     | :-: | :-:| :-: |:-:| :-:|
     | 1 | 16298 ms | 19441 ms | 27322 ms | 19233 ms |
     | 2 | 11593 ms | 29424 ms | 28955 ms | 18602 ms |
     | 3 | 24589 ms | 28685 ms | 28192 ms | 18472 ms |
     | 4 | 24725 ms | 29892 ms | 28931 ms | 18260 ms |
     | 5 | 24695 ms | 31453 ms | 36166 ms | 19474 ms |
 * 1536M 测试5次分别用时：
     | 序号 | 纯js |纯js分段 | js+wasm |js+wasm分段 | 
     | :-: | :-:| :-: |:-:| :-:|
     | 1 | 19856 ms | 19591 ms | 21259 ms | 15920 ms |
     | 2 | 15119 ms | 26283 ms | 20821 ms | 15634 ms |
     | 3 | 21387 ms | 25861 ms | 22473 ms | 16893 ms |
     | 4 | 19550 ms | 25797 ms | 21793 ms | 17239 ms |
     | 5 | 20363 ms | 26402 ms | 20782 ms | 15786 ms |
 * 1280M 测试5次分别用时：
     | 序号 | 纯js |纯js分段 | js+wasm |js+wasm分段 | 
     | :-: | :-:| :-: |:-:| :-:|
     | 1 | 6449 ms | 12169 ms | 22856 ms | 16621 ms |
     | 2 | 14695 ms | 17558 ms | 19147 ms | 18014 ms |
     | 3 | 17792 ms | 20326 ms | 17203 ms | 14683 ms |
     | 4 | 18094 ms | 16452 ms | 18396 ms | 14399 ms |
     | 5 | 15830 ms | 19006 ms | 19241 ms | 14119 ms |
 * 1024M 测试5次分别用时：
     | 序号 | 纯js |纯js分段 | js+wasm |js+wasm分段 | 
     | :-: | :-:| :-: |:-:| :-:|
     | 1 | 5003 ms | 9441 ms | 16233 ms | 9252 ms |
     | 2 | 6240 ms | 14917 ms | 11145 ms | 9316 ms |
     | 3 | 8563 ms | 10849 ms | 12653 ms | 10963 ms |
     | 4 | 10261 ms | 12155 ms | 11607 ms | 9108 ms |
     | 5 | 8775 ms | 11138 ms | 9869 ms | 10451 ms |
 * 900M  测试5次分别用时：
     | 序号 | 纯js |纯js分段 | js+wasm |js+wasm分段 | 
     | :-: | :-:| :-: |:-:| :-:|
     | 1 | 4632 ms | 7721 ms | 9590 ms | 7887 ms |
     | 2 | 5858 ms | 3312 ms | 7161 ms | 7963 ms |
     | 3 | 2859 ms | 10808 ms | 7646 ms | 7973 ms |
     | 4 | 3531 ms | 8614 ms | 7904 ms | 8197 ms |
     | 5 | 5744 ms | 7612 ms | 7131 ms | 10714 ms |
 * 800M  测试5次分别用时：
     | 序号 | 纯js |纯js分段 | js+wasm |js+wasm分段 | 
     | :-: | :-:| :-: |:-:| :-:|
     | 1 | 3329 ms | 5884 ms | 9318 ms | 7270 ms |
     | 2 | 7222 ms | 9917 ms | 6897 ms | 7096 ms |
     | 3 | 2602 ms | 6066 ms | 6295 ms | 6908 ms |
     | 4 | 2757 ms | 6662 ms | 6551 ms | 8164 ms |
     | 5 | 2509 ms | 8730 ms | 7126 ms | 7039 ms |
 * 600M  测试5次分别用时：
     | 序号 | 纯js |纯js分段 | js+wasm |js+wasm分段 | 
     | :-: | :-:| :-: |:-:| :-:|
     | 1 | 2721 ms | 2824 ms | 6557 ms | 5019 ms |
     | 2 | 3241 ms | 6867 ms | 4943 ms | 5026 ms |
     | 3 | 1803 ms | 3012 ms | 4902 ms | 5052 ms |
     | 4 | 1930 ms | 3010 ms | 5007 ms | 5022 ms |
     | 5 | 1807 ms | 2885 ms | 4881 ms | 5238 ms |
 * 400M  测试5次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 6406 ms | 3358 ms |
     | 2 | 6435 ms | 3599 ms |
     | 3 | 6450 ms | 3283 ms |
     | 4 | 6286 ms | 3952 ms |
     | 5 | 6408 ms | 3207 ms |
 * 200M  测试5次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 3497 ms | 1705 ms |
     | 2 | 3412 ms | 1643 ms |
     | 3 | 3263 ms | 1825 ms |
     | 4 | 3284 ms | 1710 ms |
     | 5 | 3376 ms | 1768 ms |
 * 100M  测试5次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 1873 ms |  923 ms |
     | 2 | 1776 ms | 928 ms |
     | 3 | 1772 ms | 913 ms |
     | 4 | 1682 ms | 923 ms |
     | 5 | 1742 ms | 898 ms |
 * 50M  测试5次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 1043 ms | 516 ms |
     | 2 | 877 ms | 479 ms |
     | 3 | 907 ms | 504 ms |
     | 4 | 872 ms | 459 ms |
     | 5 | 865 ms | 495 ms |
 * 20M  测试5次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 487 ms | 209 ms |
     | 2 | 387 ms | 209 ms |
     | 3 | 410 ms | 225 ms |
     | 4 | 512 ms | 268 ms |
     | 5 | 399 ms | 225 ms |
 * 5M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 147 ms | 92 ms |
     | 2 | 133 ms | 90 ms |
     | 3 | 177 ms | 94 ms |
     | 4 | 157 ms | 42 ms |
     | 5 | 175 ms | 84 ms |
 * 1M  测试5次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 71 ms | 20 ms |
     | 2 | 66 ms | 24 ms |
     | 3 | 45 ms | 33 ms |
     | 4 | 80 ms | 30 ms |
     | 5 | 97 ms | 29 ms |

#### firefox （版本号：103.0.1 (64 位)）
* 2048M 加载到52%时页面崩溃
    * 采用Blob.slice方式分段计算
    * 每512M为一段，测试5次
         | 序号 | 纯js分段 | js+wasm分段 |
         | :-: | :-:| :-: |
         | 1 | 51398 ms | 17338 ms |
         | 2 | 41282 ms | 16385 ms |
         | 3 | 42358 ms | 16966 ms |
         | 4 | 43363 ms | 15843 ms |
         | 5 | 40802 ms | 16551 ms |
* 1792M 加载到59%时页面崩溃
    * 采用Blob.slice方式分段计算
    * 每512M为一段，测试5次
         | 序号 | 纯js分段 | js+wasm分段 |
         | :-: | :-:| :-: |
         | 1 | 33690 ms | 13251 ms |
         | 2 | 37423 ms | 13636 ms |
         | 3 | 42903 ms | 13487 ms |
         | 4 | 32684 ms | 13662 ms |
         | 5 | 36691 ms | 14984 ms |
* 1536M 加载到69%时页面崩溃
    * 采用Blob.slice方式分段计算
    * 每512M为一段，测试5次
         | 序号 | 纯js分段| js+wasm分段 |
         | :-: | :-:| :-: |
         | 1 | 28051 ms | 11425 ms |
         | 2 | 27822 ms | 11337 ms |
         | 3 | 28331 ms | 12508 ms |
         | 4 | 30089 ms | 11520 ms |
         | 5 | 32890 ms | 11507 ms |
* 1280M 加载到83%时页面崩溃
    * 采用Blob.slice方式分段
    * 计算512M为一段
         | 序号 | 纯js分段 | js+wasm分段 |
         | :-: | :-:| :-: |
         | 1 | 25680 ms | 9571 ms |
         | 2 | 23956 ms | 9549 ms |
         | 3 | 28829 ms | 10070 ms |
         | 4 | 23518 ms | 9449 ms |
         | 5 | 23200 ms | 9540 ms |
* 1024M 测试10次分别用时：
    | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 38277 ms | 7776 ms |
     | 2 | 40936 ms | 11254 ms |
     | 3 | 29861 ms | 7653 ms |
     | 4 | 25630 ms | 7517 ms |
     | 5 | 18934 ms | 11443 ms |
     | 6 | 24849 ms | 8039 ms |
     | 7 | 18214 ms | 7727 ms |
     | 8 | 18617 ms | 12987 ms |
     | 9 | 33281 ms | 7523 ms |
     | 10 | 40757 ms | 8895 ms |
 * 900M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 22752 ms | 8605 ms |
     | 2 | 16669 ms | 9313 ms |
     | 3 | 15716 ms | 6678 ms |
     | 4 | 16940 ms | 6521 ms |
     | 5 | 16732 ms | 9269 ms |
     | 6 | 15805 ms | 6582 ms |
     | 7 | 15718 ms | 6519 ms |
     | 8 | 15795 ms| 9377 ms |
     | 9 | 15641 ms | 6773 ms |
     | 10 | 15622 ms | 7489 ms |
 * 800M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 15181 ms | 8333 ms |
     | 2 | 14031 ms | 5880 ms |
     | 3 | 14214 ms | 5987 ms |
     | 4 | 33812 ms | 5935 ms |
     | 5 | 14167 ms | 8666 ms |
     | 6 | 14666 ms | 8031 ms |
     | 7 | 28640 ms | 5991 ms |
     | 8 | 13992 ms| 5840 ms |
     | 9 | 13926 ms | 6032 ms |
     | 10 | 14216 ms | 6637 ms |
 * 600M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 11418 ms | 4457 ms |
     | 2 | 11199 ms | 5370 ms |
     | 3 | 10717 ms | 4654 ms |
     | 4 | 10607 ms | 4436 ms |
     | 5 | 10611 ms | 4479 ms |
     | 6 | 10718 ms | 4368 ms |
     | 7 | 10560 ms | 5494 ms |
     | 8 | 11519 ms| 5044 ms |
     | 9 | 10802 ms | 4426 ms |
     | 10 | 11779 ms | 4971 ms |
 * 400M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 8362 ms | 2981 ms |
     | 2 | 7516 ms | 2999 ms |
     | 3 | 7335 ms | 3030 ms |
     | 4 | 7357 ms | 3150 ms |
     | 5 | 7444 ms | 3001 ms |
     | 6 | 8456 ms | 3223 ms |
     | 7 | 7376 ms | 3120 ms |
     | 8 | 7313 ms | 3072 ms |
     | 9 | 7349 ms | 3240 ms |
     | 10 | 7447 ms | 3352 ms |

 * 200M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 4066 ms | 1525 ms |
     | 2 | 4440 ms | 1516 ms |
     | 3 | 4223 ms | 1510 ms |
     | 4 | 3916 ms | 1610 ms |
     | 5 | 3917 ms | 1509 ms |
     | 6 | 4028 ms | 1588 ms |
     | 7 | 3964 ms | 1514 ms |
     | 8 | 4037 ms| 1507 ms |
     | 9 | 3957 ms | 1506 ms |
     | 10 | 3987 ms | 1642 ms |
 * 100M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 2280 ms | 761 ms |
     | 2 | 2331 ms | 820 ms |
     | 3 | 2193 ms | 798 ms |
     | 4 | 2242 ms | 777 ms |
     | 5 | 2197 ms | 752 ms |
     | 6 | 2330 ms | 769 ms |
     | 7 | 2236 ms | 758 ms |
     | 8 | 2364 ms| 798 ms |
     | 9 | 2278 ms | 783 ms |
     | 10 | 2384 ms | 785 ms |
 * 50M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 1366 ms | 397 ms |
     | 2 | 1355 ms | 378 ms |
     | 3 | 1445 ms | 460 ms |
     | 4 | 1468 ms | 437 ms |
     | 5 | 1417 ms | 406 ms |
     | 6 | 1525 ms | 478 ms |
     | 7 | 1381 ms | 393 ms |
     | 8 | 1450 ms| 430 ms |
     | 9 | 1417 ms | 428 ms |
     | 10 | 1378 ms | 431 ms |
 * 20M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 921 ms | 168 ms |
     | 2 | 871 ms | 162 ms |
     | 3 | 859 ms | 163 ms |
     | 4 | 864 ms | 162 ms |
     | 5 | 1025 ms | 177 ms |
     | 6 | 910 ms | 158 ms |
     | 7 | 904 ms | 150 ms |
     | 8 | 931 ms| 187 ms |
     | 9 | 1014 ms | 182 ms |
     | 10 | 871 ms | 159 ms |
 * 5M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 127 ms | 48 ms |
     | 2 | 124 ms | 50 ms |
     | 3 | 140 ms | 44 ms |
     | 4 | 129 ms | 47 ms |
     | 5 | 127 ms | 51 ms |
     | 6 | 129 ms | 50 ms |
     | 7 | 126 ms | 46 ms |
     | 8 | 119 ms| 54 ms |
     | 9 | 121 ms | 46 ms |
     | 10 | 118 ms | 50 ms |
 * 1M  测试10次分别用时：
     | 序号 | 纯js | js+wasm |
     | :-: | :-:| :-: |
     | 1 | 46 ms | 18 ms |
     | 2 | 41 ms | 22 ms |
     | 3 | 43 ms | 13 ms |
     | 4 | 40 ms | 15 ms |
     | 5 | 44 ms | 11 ms |
     | 6 | 47 ms | 15 ms |
     | 7 | 42 ms | 11 ms |
     | 8 | 42 ms| 20 ms |
     | 9 | 45 ms | 13 ms |
     | 10 | 44 ms | 16 ms |


### 分段计算测试代码
#### 纯js
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>文件md5</title>
  <script src="./SparkMD5.js"></script>
</head>
<body>
  <input id="file" type="file" />
  <script>
  document.querySelector('#file').addEventListener('change', e => {
      let startTime = Date.now()
      const file = e.target.files[0];
      const fileReader = new FileReader()
      console.log('size', file.size / 1024 / 1024 / 1024, "G")
      fileReader.onprogress = e => {
        console.log(`${Math.floor((e.loaded / e.total) * 100)}%`)
      }
      let usedTime = 0
      const md5 = new SparkMD5();
      let index = 0
      const chunkSize = 512 * 1024 * 1024;//file.size / count
      let count = Math.ceil(file.size / chunkSize)
      console.log('分几份', count)
      loadSliceFile();
      function loadSliceFile() {
        const sliceFile = file.slice(index * chunkSize, index * chunkSize + chunkSize)
        fileReader.readAsBinaryString(sliceFile);
      }
      fileReader.onload = e => {
        index += 1;
        md5.appendBinary(e.target.result);
        if (index < count) {
          loadSliceFile()
        }
        else {
          const md5Str = md5.end()
          usedTime += Date.now() - startTime
          console.log('usedTime', usedTime, 'ms')
          console.log('md5', md5Str)
        }
      }
    });
  </script>
</body>
</html>

```
#### js+wasm
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>文件md5</title>
  <script src="./wasm_exec.js"></script>
  <!-- <script src="./wasm_exec_tiny.js"></script> -->
</head>
<body>
  <script>
    function handleSayHello(message) {
      console.lof('str from go', message)
    }
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch('wasm/md5.wasm'), go.importObject)
      .then(res => {
        go.run(res.instance); // 执行 golang里 main 方法
      });
  </script>
  <input id="file" type="file" />
  <script>
    document.querySelector('#file').addEventListener('change', e => {
      let startTime = Date.now()
      const file = e.target.files[0];
      const fileReader = new FileReader()
      console.log('size', file.size / 1024 / 1024 / 1024, "G")
      fileReader.onprogress = e => {
        console.log(`${Math.floor((e.loaded / e.total) * 100)}%`)
      }
      let usedTime = 0
      let index = 0
      const sliceSize = 512
      const chunkSize = sliceSize * 1024 * 1024;//file.size / count
      let count = Math.ceil(file.size / chunkSize)
      console.log('分几份', count)

      loadSliceFile();
      function loadSliceFile() {
        const sliceFile = file.slice(index * chunkSize, index * chunkSize + chunkSize)
        fileReader.readAsArrayBuffer(sliceFile);
      }
      fileReader.onload = e => {
        index += 1;
        const bytes = new Uint8Array(e.target.result)
        wasmMd5Add(bytes)
        if (index < count) {
          loadSliceFile()
        }
        else {
          const md5Hash = wasmMd5End()
          usedTime += Date.now() - startTime
          console.log('usedTime', usedTime, 'ms')
          console.log('md5', md5Hash)
        }
      }
    });
  </script>
</body>
</html>

```
### 测试结论
#### firefox
* 超过1G的文件，直接崩溃，只能通过分段计算最终合并计算
* 从1M到2G，wasm的速度是纯js计算的2-3倍
* 20M，wasm是纯js的 6倍
#### chrome
* 0-400M时，wasm是纯js的2-3倍
* 600M-1024M时，纯js不分段比wasm要快
    * 分段js比不分段wasm快一点点
    * 分段js比分段wasm慢一点点
* 1280M，差不太多
* 大于1280M,js比wasm分段慢
* 对于js，分段要慢一些
* 对于wasm,分段要快一些


### 最终结论
* chrome对js的优化，使得在600M-1024M期间的大文件纯js计算md5速度要快于wasm，其他范围还是wasm性能好一些
* 由于firefox超过1G就崩溃了，所以我们平时写代码时，还是要做分段加载的。
* 业务中，还是可以使用wasm来提升性能的
* 可以针对 chrome与其他浏览器来制作不同的方案
* 其实golang 计算md5基本上是js的7-9倍，但js给wasm复制数据的时间占用了太多，导致wasm被降低了速度，文件越大，复制时间越长，越慢

### 结束语
wasm 还是可以使用的，众观全局，速度提升2-3倍。chrome可以针对性处理
