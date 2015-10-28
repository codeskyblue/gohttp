## Self Signed key and Nginx configuration

文章参考：<http://wangye.org/blog/archives/732/>

## 网站证书
1. 首先需要一个根证书

     准备公钥私钥
        
        openssl genrsa -out ca.key 2048

     利用ca.key创建根证书，其中/C是国家,ST：State L: Locality Name也就是City，O公司名，OU这里的内容会显示在颁发者这一栏，填网址比较好

        openssl req -new -x509 -days 36500 -key ca.key -out ca.crt -subj \
        "/C=CN/ST=Jiangsu/L=Yangzhou/O=Your Company Name/OU=Your Root CA"

2. 准备网站的请求证书

     创建ssl私钥

        openssl genrsa -out server.key 2048

     创建csr证书签名请求文件,其中CN必须为网站域名

        openssl req -new -key server.key -out server.csr -subj \
        "/C=CN/ST=Zhejiang/L=Hangzhou/O=Netease/OU=netease.com/CN=teamtoy.mt.nie.netease.com"

3. 用根证书签名

     做些准备工作：

        mkdir demoCA
        cd demoCA
        mkdir newcerts
        touch index.txt
        echo '01' > serial
        cd ..

     利用刚才的CA证书签下名， 有提示的话，直接yes就可以了
     
        openssl ca -in server.csr -out server.crt -cert ca.crt -keyfile ca.key

## nginx上的配置

一个典型的配置是

```
server {
     listen 443 ssl;
     server_name localhost;
     ssl on;
     ssl_certificate /etc/server.crt;
     ssl_certificate_key /etc/server.key;

     location / {
        proxy_pass http://localhost:4000;
        proxy_redirect     off;
        proxy_set_header   Host             $host;
        proxy_set_header   X-Real-IP        $remote_addr;
        proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
    }
}
```

或者增强一点，开一个http的然后所有的请求都跳转到https

```
server {
     listen 80;
     server_name localhost;
     return 301 https://localhost$request_uri;
}
```