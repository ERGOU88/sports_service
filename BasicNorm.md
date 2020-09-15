 Golang 基本规范
---------
**规范主要是为了让代码有层次，风格统一，便于以后维护，规范仅为参考**

## gofmt
大部分的格式问题可以通过gofmt解决，gofmt自动格式化代码，保证所有的go代码一致的格式。正常情况下，采用Sublime编写go代码时，插件GoSublilme已经调用gofmt对代码实现了格式化。

## 行长
一行代码最好不超过一屏大小，80字符为好。

## 注释
在编码阶段同步写好变量、函数、包注释，注释可以通过godoc导出生成文档。
- 包注释
一个位于package子句之前的块注释或行注释。包如果有多个go文件，只需要出现在一个go文件中即可。

 ``` go
//Package regexp implements a simple library 
//for regular expressions
package regexp
 ``` 
- 可导出类型

 ``` go
// Compile parses a regular expression and returns, if successful, a Regexp
// object that can be used to match against text.
func Compile(str string) (regexp *Regexp, err error) {
 ```

## 命名
使用短命名，长名字并不会自动使得事物更易读，文档注释会比格外长的名字更有用，使用可搜索的名称：单字母名称和数字常量很难从一大堆文字中搜索出来。

- 包名
包名应该为小写单词，不要使用下划线或者混合大小写。

- 接口名
单个函数的接口名以"er"作为后缀，如Reader,Writer
接口的实现则去掉“er”
```go
	type Reader interface {
        Read(p []byte) (n int, err error)
	}
```
 两个函数的接口名综合两个函数名
``` go
	type WriteFlusher interface {
	    Write([]byte) (int, error)
	    Flush() error
	}
```
 三个以上函数的接口名，类似于结构体名
``` go
	type Car interface {
	    Start([]byte) 
	    Stop() error
	    Recover()
	}
```

- 混合大小写
采用驼峰式命名
>MixedCaps 大写开头，可导出
>mixedCaps 小写开头，不可导出

- 变量
全局变量：驼峰式，结合是否可导出确定首字母大小写
参数传递：驼峰式或全小写，小写字母开头
局部变量：可以用下划线形式或者全部小写
bool变量：使用Is-Has-Can-Allow开头

- 常量（不要出现任何magic number）
常量均需使用全部大写字母组成，并使用下划线分词
>const APP_VER = "1.0"

- 枚举类型
枚举类型的常量，需要先创建相应类型

``` go
	type Scheme string
	const (
	    HTTP  Scheme = "http"
	    HTTPS Scheme = "https"
	)
	
	type PullRequestStatus int
	const ( // 带有前缀标识
	    PULL_REQUEST_STATUS_CONFLICT PullRequestStatus = iota
	    PULL_REQUEST_STATUS_CHECKING
	    PULL_REQUEST_STATUS_MERGEABLE
	)
```
- struct 结构体
struct申明和初始化格式采用多行
定义和初始化：
``` go
	type User struct{
	    Username  string
	    Email     string
	}
	
	u := User{
	    Username: "test",
	    Email:    "test@gmail.com",
	}
```


## 控制结构
- if 
if接受初始化语句，约定如下方式建立局部变量:
```go
	if err := file.Chmod(0664); err != nil {
	    return err
	}
```
 
- for
采用短声明建立局部变量:
```go
	sum := 0
	for i := 0; i < 10; i++ {
	    sum += i
	}
```

- range
如果只需要第一项（key），就丢弃第二个:
```go
	for key := range m {
	    if key.expired() {
	        delete(m, key)
	    }
	}
```
- return
尽早return：一旦有错误发生，马上返回
```go
	f, err := os.Open(name)
	if err != nil {
	    return err
	}
	
	d, err := f.Stat()
	if err != nil {
	    f.Close()
	    return err
	}
	
	codeUsing(f, d)
```


##函数
- 函数采用命名的多值返回
- 传入变量和返回变量以小写字母开头  
 >func nextInt(b []byte, pos int) (value, nextPos int) {
  
  在godoc生成的文档中，带有返回值的函数返回多个参数的时声明更利于理解
  多个函数定义之间留一个空行
```go
  func Add(a, b int) (result int ) {
	  return a+b
  }
  
  // 留一个空行
  func Sub(a, b int) (result int) {
	  return a-b
  }
```
  函数中代码，最好使用空行分块
```go
  func (dm *RegModule) ItgCLogin() (e error){
	  var login, password string
	  if e := req.Parse("name", &login, "pwd", &password); e != nil {
		  return 
	  }

	  r, e := re.Login(login, password, req.IP())
	  if e != nil {
	      return 
	  }

	  result["res"] = r
	  return
  }
```
函数中定义多个变量可以用var ()定义一组

```go
var (
   num1 int
   name string
   x    float
)
```

## 错误处理
- error作为函数的值返回,必须对error进行处理
- 错误描述如果是英文必须为小写，不需要标点结尾
- 采用独立的错误流进行处理

 不要采用这种方式
``` go
    if err != nil {
        // error handling
    } else {
        // normal code
    }
```

 而要采用下面的方式
``` go
	if err != nil {
        // error handling
        return // or continue, etc.
    }
    // normal code
```
如果返回值需要初始化，则采用下面的方式
``` go
	x, err := f()
	if err != nil {
	    // error handling
	    return
	}
	// use x
```

##panic

- 尽量不要使用panic，除非你知道你在做什么

## import
- 对import的包进行分组管理，而且标准库作为第一组
``` go
	package main
	
	import (
	    "fmt"
	    "hash/adler32"
	    "os"
	
	    "appengine/user"
	    "appengine/foo"
	
	    "code.google.com/p/x/y"
	    "github.com/foo/bar"
	)
```

## 缩写
- 采用全部大写或者全部小写来表示缩写单词
比如对于url这个单词，不要使用
``` go
	UrlPony  
``` 
 
 而要使用
``` go
	urlPony 或者 URLPony  
```

## 参数传递

- 对于少量数据，不要传递指针
- 对于大量数据的struct可以考虑使用指针
- 传入参数是map，slice，chan不要传递指针
因为map，slice，chan是引用类型，不需要传递指针的指针

## 接受者

- 名称
统一采用单字母'p'而不是this，me或者self
``` go
	type T struct{} 

	func (p *T)Get(){}
```

- 类型
对于go初学者，接受者的类型如果不清楚，统一采用指针型
``` go
	func (p *T)Get(){}
```
而不是
``` go
	func (p T)Get(){}   
```
在某些情况下，出于性能的考虑，或者类型本来就是引用类型，有一些特例

- 如果接收者是map,slice或者chan，不要用指针传递
``` go
	//Map
	package main
	
	import (
	    "fmt"
	)
	
	type mp map[string]string
	
	func (m mp) Set(k, v string) {
	    m[k] = v
	}
	
	func main() {
	    m := make(mp)
	    m.Set("k", "v")
	    fmt.Println(m)
	}
```
 
``` go
	//Channel
	package main
	
	import (
	    "fmt"
	)
	
	type ch chan interface{}
	
	func (c ch) Push(i interface{}) {
	    c <- i
	}
	
	func (c ch) Pop() interface{} {
	    return <-c
	}
	
	func main() {
	    c := make(ch, 1)
	    c.Push("i")
	    fmt.Println(c.Pop())
	}
```

```go
// 如果需要对slice进行修改，通过返回值的方式重新赋值
// Slice
package main

import (
    "fmt"
)

type slice []byte

func main() {
    s := make(slice, 0)
    s = s.addOne(42)
    fmt.Println(s)
}

func (s slice) addOne(b byte) []byte {
    return append(s, b)
}
```
 
```go
//如果接收者是含有sync.Mutex或者类似同步字段的结构体，必须使用指针传递避免复制
package main

import (
    "sync"
)

type T struct {
    m sync.Mutex
}

func (t *T) lock() {
    t.m.Lock()
}

/*
Wrong !!!
func (t T) lock() {
    t.m.Lock()
}
*/

func main() {
    t := new(T)
    t.lock()
}
```

```go
// 如果接收者是大的结构体或者数组，使用指针传递会更有效率。

package main

import (
    "fmt"
)

type T struct {
    data [1024]byte
}

func (t *T) Get() byte {
    return t.data[0]
}

func main() {
    t := new(T)
    fmt.Println(t.Get())
}
```

