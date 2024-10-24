# 作用？
作用是修改你的hosts，使其自动指向一个IP地址，对github连不上有奇效
# 怎么用？
以管理员身份运行的同时带上域名参数
比如我用powershell打开
~~~bash
./hostalter.exe github.com
~~~
就会把github.com的最快ip修改进你的hosts，同时删除已有的github.com记录