package main

const htmlLogin = `
<!DOCTYPE html>
<html lang="zh-cmn-Hans">

<head>
    <meta http-equiv="content-type" content="text/html;charset=utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>VSCODE-TAODEV:LOGIN</title>
</head>
<body>
<div>
    <a href="{{.URL}}">{{.Label}}</a>
</div>
</body>
</html>
`

const htmlError = `
<!DOCTYPE html>
<html lang="zh-cmn-Hans">

<head>
    <meta http-equiv="content-type" content="text/html;charset=utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>VSCODE-TAODEV:ERROR</title>
</head>
<body>
<div>
    {{.Message}}
</div>
</body>
</html>
`
