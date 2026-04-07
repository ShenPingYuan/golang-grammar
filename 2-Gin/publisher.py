import xmlrpc.client
import os
import glob
import re

# ============ 配置 ============
BLOG_URL = "https://rpc.cnblogs.com/metaweblog/michaelshen"
USERNAME = "沈平元"
TOKEN = "4BDED7F725F681A86DCD332AEE430330702380D7893D801C885D7BC012E2884C"
MD_DIR = "./posts"  # markdown 文件所在目录
# ==============================

server = xmlrpc.client.ServerProxy(BLOG_URL)

def clean_title(filename):
    """去掉序号前缀和.md后缀
    '1-Gin 框架进阶系列（一）：第一个路由.md' -> 'Gin 框架进阶系列（一）：第一个路由'
    """
    name = filename.replace(".md", "")
    name = re.sub(r"^\d+-", "", name)
    return name

def publish(title, content, categories=["Gin","Go","[Gin 框架进阶系列]"], publish=True):
    post = {
        "title": title,
        "description": content,
        "categories": categories,
        "mt_keywords": "Go语言,Gin",  # 标签
    }
    post_id = server.metaWeblog.newPost(
        "dummy",   # blogid（cnblogs 忽略此参数）
        USERNAME,
        TOKEN,
        post,
        publish    # True=公开, False=草稿
    )
    return post_id

# 按文件名排序，依次发布
files = sorted(glob.glob(os.path.join(MD_DIR, "*.md")))

for f in files:
    filename = os.path.basename(f)
    # 用文件名（去掉.md）作为标题，也可以自定义
    title = clean_title(filename)

    with open(f, "r", encoding="utf-8") as fp:
        content = fp.read()

    try:
        post_id = publish(title, content, publish=False)  # 先发草稿，确认没问题再改True
        print(f"✅ {title} -> ID: {post_id}")
    except Exception as e:
        print(f"❌ {title} -> {e}")