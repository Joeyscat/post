

# Post

## 发布Post
### 请求
请求方法： POST

| 参数名称 | 必选 | 类型 | 描述 |
| ---- | ---- | ---- | ---- | 
| content | 否 | String | 文本内容 |
| pics | 否 | Array of String | 图片ID数组 |
| video | 否 | String | 视频ID |

> text、pics、video 至少包含其中一个

* 请求示例

```bash
POST /post

Content-Type: application/json

{
    "content": "xxx",
    "pics": [12, 23]
}
```
* 响应示例

```json
201 OK

{
    "postId": 12
}
```

---
## 删除Post
### 请求
请求方法： DELETE

| 参数名称 | 必选 | 类型 | 描述 |
| ---- | ---- | ---- | ---- | 
| postId | 是 | Integer | Post ID |

* 请求示例

```bash
DELETE /post?postId=12

```

### 响应
* 响应示例

```bash
200 OK

{
    "postId": 12
}
```
---
## Post列表

### 请求
方法： GET

路径： /post

| 参数名称 | 必选 | 类型 | 描述 |
| ---- | ---- | ---- | ---- | 
| userId | 否 | Integer | User ID |
| page | 否 | Integer | 页码 |

* 请求示例
```bash
GET /post?userId=22&page=10
```

### 响应

| 参数名称 | 类型 | 描述 |
| ----  | ---- | ---- | 
| posts | Array of Post | 帖子数组 |
| posts.N.user | User | 用户信息 |
| posts.N.text | String | 帖子的文本内容 |
| posts.N.pics | Array of String | 帖子附加的图片URL数组 |
| posts.N.video | String | 帖子附加的视频 |
| posts.N.comments | Array of Comment | 帖子评论数组 |

* 响应示例

```json
200 OK

{
    "posts": [
        {
            "user": {},
            "text": "xxx",
            "pics": [
                "http://hi.cn/12.jpg", "http://hi.cn/23.jpg"
            ],
            "video": "http://hi.cn/v/1",
            "comments": [
                {
                    "user": {},
                    "text": "haa",
                    "time": 1672559083
                },
            ]
        }
    ]
}
```

## 评论Post
### 请求
请求方法： POST

路径： /post/{id}/comment

| 参数名称 | 必选 | 类型 | 描述 |
| ---- | ---- | ---- | ---- | 
| id | 是 | String | Post ID |
| content | 是 | String | 评论内容 |

* 请求示例

```bash
POST /post/12rmj1WtM/comment

Content-Type: application/json

{
    "content": "hello"
}
```
* 响应示例

```json
200 OK

```

# Media

## 上传Media
### 请求
请求方法： POST

路径： /media/upload

| 参数名称 | 必选 | 类型 | 描述 |
| ---- | ---- | ---- | ---- | 
| file | 是 | Binary | 媒体文件 |
| video | 否 | String | 视频ID |

* 请求示例

```bash
POST /media/upload

(file content)
```

* 响应示例

```json
201 OK

{
    "mediaId": 12
}
```

---

## Media列表

### 请求
方法： GET

路径： /media

| 参数名称 | 必选 | 类型 | 描述 |
| ---- | ---- | ---- | ---- | 
| userId | 是 | Integer | User ID |
| page | 否 | Integer | 页码 |

* 请求示例

```bash
GET /media?userId=22&page=10
```

### 响应

| 参数名称 | 类型 | 描述 |
| ----  | ---- | ---- | 
| medias | Array of Media | 媒体文件数组 |
| medias.N.type | String | Media 文件类型 |
| medias.N.url | String | 文件链接 |
| medias.N.commentNum | Integer | 评论数量 |

* 响应示例

```json
200 OK

[
  {
    "id": "bMCUw6WtM",
    "type": 1,
    "url": "http://localhost:8000/m/bMCUw6WtM.docx",
    "posted": false,
    "time": 1672658573
  },
  {
    "id": "wwLI6wWtM",
    "type": 1,
    "url": "http://localhost:8000/m/wwLI6wWtM.docx",
    "posted": false,
    "time": 1672658572
  }
]
```
---

## Media详情

### 请求
方法： GET

路径： /media/{id}

| 参数名称 | 必选 | 类型 | 描述 |
| ---- | ---- | ---- | ---- | 
| id | 是 | Integer | Media ID |

* 请求示例

```bash
GET /media/12
```

### 响应

| 参数名称 | 类型 | 描述 |
| ----  | ---- | ---- | 
| type | String | Media 文件类型 |
| url | String | 文件链接 |
| comments | Array of Comment | Media评论数组 |

* 响应示例

```json
200 OK

{
    "id": 12,
    "type": "video",
    "url": "http://hi.cn/v/1",
    "comments": [
        {
            "user": {},
            "text": "haa",
            "time": 1672559083
        },
    ],
    "time": 1672559083
}
```
---






