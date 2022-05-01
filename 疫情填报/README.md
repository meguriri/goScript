## 智慧沈航-疫情信息填报脚本

>一些设计借鉴了[Github链接](https://github.com/DOEMsy/MyScripts/tree/master/SAU%E8%87%AA%E5%8A%A8%E7%AD%BE%E5%88%B0/normal), 在此特别感谢!!!
>
>并加入了一些新的内容.如: 解析响应json,添加了配置文件方便使用者填入自己参数等.



### 使用方法

#### 1. 在 info.yaml 中填入自己的信息

```yaml
#user
user:
  username: 学号
  password: 智慧沈航密码

#form
form:
  name: 姓名
  stunum: 学号
  tel: 电话
  college: 学院
  province: 当前省份
  city: 当前城市
  id: 见下文
```

#### 2. form中id的获取

这个id需要自行获取:

URL:[链接](https://app.sau.edu.cn/form/wap/default/index?formid=10)

用浏览器自带的抓包抓取POST(save?formid=10)中的表单参数,即可得到id值

#### 3. 结果

提交正确得结果为:

```
获取Cookies:  操作成功
上传填报信息:  操作成功 
```



### 注:

仅供go net/http学习
