# aliyun-ddns
通过阿里云的API接口，为动态公网IP实现域名绑定（DDNS）

## 安装

```bash
go get github.com/lonord/aliyun-ddns
```

如果需要编译其他平台的二进制文件（比如为树莓派编译），进入目录，使用`build.sh`来编译

```bash
./build.sh -l=arm    # 编译arm版linux的二进制文件
```

## 使用

### 参数

```bash
aliyun-ddns --help
Usage of ./dist/darwin/amd64/aliyun-ddns:
  -domain string
        Domain name (like google.com)
  -key string
        Access Key ID
  -region string
        Region ID (default "cn-hangzhou")
  -rr string
        Resource record (RR)
  -secret string
        Access Key Secret
  -type string
        Domain type (A,CNAME,MX,etc...) (default "A")
  -v    Show version
```

#### -domain
指定DNS解析的根域名，也可以通过环境变量`ALIDNS_DOMAIN`来设置

#### -key
阿里云API服务访问key，也可以通过环境变量`ALIDNS_ACCESS_KEY`来设置

#### -region
*可选参数*

阿里云API服务地区，一般选择最近的区域，默认为`cn-hangzhou`，（服务区域ID可以在[阿里云官网](https://help.aliyun.com/document_detail/40654.html?spm=5176.10695662.1996646101.1.2b4a33dcFrxth0)找到），也可以通过环境变量`ALIDNS_REGION`来设置

#### -rr
指定需要动态更新的主机记录，也可以通过环境变量`ALIDNS_RR`来设置

#### -secret
阿里云API服务访问密钥，也可以通过环境变量`ALIDNS_ACCESS_SECRET`来设置

#### -type
*可选参数*

指定需要动态更新的主机记录的类型（A，CNAME，MX等），默认为`A`，也可以通过环境变量`ALIDNS_DDMAIN_TYPE`来设置

### 放在定时调度（cron）中使用

例如每10分支更新一次`subdomain.abc.com`这个域名的DDNS记录，在`/etc/cron.d/`下创建一个文件`updatedns`

```bash
# The first element of the path is a directory where the debian-sa1
# script is located
PATH=/root/bin:/usr/sbin:/usr/sbin:/usr/bin:/sbin:/bin

ALIDNS_REGION=cn-shanghai
ALIDNS_ACCESS_KEY=xxxxxxxxxxx
ALIDNS_ACCESS_SECRET=xxxxxxxxxxxxx
ALIDNS_DOMAIN=abc.com
ALIDNS_DDMAIN_TYPE=A
ALIDNS_RR=subdomain

5-55/10 * * * * root aliyun-ddns | logger -i -t update-dns
```

## License

MIT
