package doc

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"sports_service/global/app/errdef"
	"sports_service/global/consts"
)

func ApiCode(c *gin.Context) {
	const tpl = `
<!DOCTYPE html>
<html lang="en">
 <head>
  <meta charset="UTF-8" />
  <title>Open API Code Doc</title>
  <style type="text/css">

table.dataintable {
	margin-top:15px;
	border-collapse:collapse;
	border:1px solid #aaa;
	width:100%;
	}

table.dataintable th {
	vertical-align:baseline;
	padding:5px 15px 5px 6px;
	background-color:#3F3F3F;
	border:1px solid #3F3F3F;
	text-align:left;
	color:#fff;
	}

table.dataintable td {
	vertical-align:text-top;
	padding:6px 15px 6px 6px;
	border:1px solid #aaa;
	}

table.dataintable tr:nth-child(odd) {
	background-color:#F5F5F5;
}

table.dataintable tr:nth-child(even) {
	background-color:#fff;
}
	</style>
 </head>
 <body>
  <table class="dataintable">
   <tbody>
    <tr>
     <th>Code</th>
     <th>描述</th>
    </tr>
	{{range $index, $elem := .Codes}}
		<tr>
			 <td>{{$index}}</td>
			 <td>{{$elem}}</td>
		</tr>
	{{end}}
   </tbody>
  </table>
 </body>
</html>
`
	type lists struct {
		Codes map[int]string
	}

	t, _ := template.New("webpage").Parse(tpl)
	data := lists{
		Codes: errdef.MsgFlags,
	}
	c.Header("Content-Type", "text/html; charset=utf-8")

	_ = t.Execute(c.Writer, data)
}

func NotifyDoc(c *gin.Context) {
	const tpl = `
<!DOCTYPE html>
<html lang="en">
 <head>
  <meta charset="UTF-8" />
  <title>Open API Code Doc</title>
  <style type="text/css">

table.dataintable {
	margin-top:15px;
	border-collapse:collapse;
	border:1px solid #aaa;
	width:100%;
	}

table.dataintable th {
	vertical-align:baseline;
	padding:5px 15px 5px 6px;
	background-color:#3F3F3F;
	border:1px solid #3F3F3F;
	text-align:left;
	color:#fff;
	}

table.dataintable td {
	vertical-align:text-top;
	padding:6px 15px 6px 6px;
	border:1px solid #aaa;
	}

table.dataintable tr:nth-child(odd) {
	background-color:#F5F5F5;
}

table.dataintable tr:nth-child(even) {
	background-color:#fff;
}
	</style>
 </head>
 <body>
  <table class="dataintable">
   <tbody>
    <tr>
     <th>MsgType</th>
     <th>业务场景</th>
    </tr>
	{{range $index, $elem := .Codes}}
		<tr>
			 <td>{{$index}}</td>
			 <td>{{$elem}}</td>
		</tr>
	{{end}}
   </tbody>
  </table>
 </body>
</html>
`
	type lists struct {
		Codes map[consts.MessageType]string
	}

	t, _ := template.New("webpage").Parse(tpl)
	data := lists{
		Codes: consts.NotifyDoc,
	}
	c.Header("Content-Type", "text/html; charset=utf-8")

	_ = t.Execute(c.Writer, data)
}
