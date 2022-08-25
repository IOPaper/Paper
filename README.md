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
    <a href="README.zh.md">简体中文</a>
</p>

## RESTFul API Doc

#### Generic Response Structure
##### response example
```json5
{
  "status": true,
  "msg": "msg is nullable field",
  "data": {
    "msg": "data is nullable field"
  }
}
```
##### response field type
 - `status<bool>`
 - `msg<string, nullable>`
 - `data<any, nullable>`

_as a `data` field with `any` type, its role is to carry unknown data_

-------

### GetPaperList

**method:** `GET`

**path:** `/paper/list`

**query options:**
 - `before<int, nullable>`
 - `limit<int, max:10, nullable>`

**success response**
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

**error response**
 - `status<bool>`
 - `msg<string>`

-------

### GetPaper

**method:** `GET`

**path:** `/paper/:paper_id`

**url param**
 - `paper_id<string>`

**success response**
 - `status<bool>`
 - `data<PaperExport>`

**error response**
 - `status<bool>`
 - `msg<string>`
 
 ## Roadmap
  - **Support for paper encryption**
  - **Support for blockchain storage**
