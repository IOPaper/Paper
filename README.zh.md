<img align="center" src=".doc/io.paper.png" />

<p align="center">
    <a href="https://github.com/IOPaper/Paper/releases">
        <img src="https://img.shields.io/github/v/release/IOPaper/Paper?include_prereleases&style=flat-square"/>
    </a>
    <a href="https://goreportcard.com/report/github.com/IOPaper/Paper">
        <img src="https://goreportcard.com/badge/github.com/IOPaper/Paper?style=flat-square">
    </a>
</p>

<p align="left">
    <a href="README.md">English</a>
</p>

## RESTFul 接口文档

#### 通用响应结构
##### 响应示例
```json5
{
  "status": true,
  "msg": "msg is nullable field",
  "data": {
    "msg": "data is nullable field"
  }
}
```
##### 响应字段类型
- `status<bool>`
- `msg<string, nullable>`
- `data<any, nullable>`

_作为`any`类型的`data`字段，其作用是携带未知数据_

-------

### GetPaperList

**method:** `GET`

**path:** `/paper/list`

**query options:**
- `before<int, nullable>`
- `limit<int, max:10, nullable>`

**成功响应**
- `status<bool>`
- `data<Array<PaperExport>>`
- PaperExport
    - `paper_id<string>`
    - `title<string>`
    - `content<string>`
    - `tags<Array<string>, nullable>`
    - `attachment<Array<string>, nullable>`
    - `author<string>`
    - `sign<string[base64], nullable>`
    - `date_create<date>`
    - `date_modified<date, nullable>`

**错误响应**
- `status<bool>`
- `msg<string>`

-------

### GetPaper

**method:** `GET`

**path:** `/paper/:paper_id`

**url param**
- `paper_id<string>`

**成功响应**
- `status<bool>`
- `data<PaperExport>`

**错误响应**
- `status<bool>`
- `msg<string>`