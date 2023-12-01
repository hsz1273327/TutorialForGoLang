
# 使用golang连接常用的数据库和消息中间件

写服务几乎必定会用到数据库和消息中间件,常用到的有:

+ 数据库,用于保存和查询数据,最常见的是关系数据库
+ 缓存,用于保存临时数据或多个服务共享数据,最常见的是redis
+ 消息中间件,用于构造数据流,生产消费模式等,通常不对数据做处理,只是用于构造更加高性能高可维护性的结构件

## 关系数据库

我们最常用的数据库还是关系数据库,他们都使用sql语言作为操作语言,应用最多的数据库产品就是`mysql`,`postgresql`和`sqlite`.一般写㐏也不直接连接使用SQL控制这些数据库,而是使用orm,这样虽然有一定的性能损失,但可以防止sql注入,同时代码也更加容易维护.我更推荐使用[xorm](https://xorm.io/zh/),这个库简单够用,同时提供一个周边工具可以把已有的数据库结构生成为对应的结构体.它还有个威力加强版[github.com/xormplus/xorm](https://www.kancloud.cn/xormplus/xorm/167077)定义了更多的操作,当然也更重些.

使用orm更加适合多语言联合开发项目,向go基本是做服务端开发,python一般做原型开发或者mvp,dba则是直接使用客户端在命令行中操作数据库,数据分析数据挖掘的通常也是使用sql语句或者python连接数据库做处理取数据.可以看出go只是整个数据流中的一环中使用到,而且往往是做的写入操作,更加偏向于业务.相对于其他orm框架xorm最大的好处是

1. 可以直接执行sql语句应付复杂请求
2. 有工具直接导出已存在的数据库表到结构体,不需要围绕这个orm重新设计数据库.

本部分例子在[这里](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/xorm_test)
xorm需要安装如下3个部分

+ orm本体

```
go get github.com/go-xorm/xorm
```

+ 扩展

```
go get github.com/xormplus/xorm
```

+ orm对应工具

```
go get github.com/go-xorm/cmd/xorm
```

+ 数据库驱动(选需要的装)

    + Mysql: `github.com/go-sql-driver/mysql`
    + Postgres: `github.com/lib/pq`
    + SQLite: `github.com/mattn/go-sqlite3`
    
    
使用的时候这么引入(本文以使用postgres做例子)

```golang
import (
     _ "github.com/lib/pq"
    "github.com/go-xorm/xorm"
)
```

xorm使用`NewEngine`函数初始化一个数据库客户端对象
```golang
db, err := xorm.NewEngine("postgres","postgres://postgres:postgres@localhost:5432/test?sslmode=disable")
```


如果需要设置连接池的空闲数大小,可以使用`db.SetMaxIdleConns()`来实现


如果需要设置最大打开连接数则可以使用`db.SetMaxOpenConns()`来实现


我们也可以手动关闭这个数据库客户端对象
```golang
defer db.Close()
```

如果我们希望验证数据库操作,要打印出操作对应的sql语句,可以使用接口`ShowSQL(true)`
```golang
db.ShowSQL(true)
```

### 定义表结构

xorm使用结构体定义表,我们可以在结构体中对应的字段上使用如下例的注解来细化定义表结构

```
`xorm:"varchar(25) notnull unique 'usr_name'"`
```

```golang
type Goods struct {
    Id    int `xorm:"not null pk autoincr INTEGER"`
    Price uint
}

```

### 将数据库中已有的表结构导出

xorm相比起其他orm框架最大的优势就是他有工具支持将数据库中已有的表结构导出.这就让他有了支持先定义数据库再写业务逻辑这样工作流的能力.

```bash
xorm reverse [-s] driverName datasourceName tmplPath [generatedPath] [tableFilterReg]
```

+ `driverName`是数据库的名字,比如`postgres`
+ `datasourceName`是数据库的连接配置字符串,比如`"postgres://postgres:postgres@localhost:5432/test?sslmode=disable"`
+ `tmplPath`是模板所在的位置,通常在`GOPATH/src/github.com/go-xorm/cmd/xorm/templates/goxorm`
+ `generatedPath`是将模板放在什么位置,比如在当前目录就设置`.`,不填就会生成在当前目录下的`model`文件夹下
+ `tableFilterReg`则是用于筛选要导出的表名的re字符串.如果不填就会导出数据库中所有的表

### 表操作

常用的表操作有:

+ 同步表`db.Sync2(...string|*struct)err`参数为一个或多个空的对应Struct的指针或表名字符串.Sync2函数将进行如下的同步操作:

    + 自动检测和创建表,这个检测是根据表的名字
    + 自动检测和新增表中的字段,这个检测是根据字段名,同时对表中多余的字段给出警告信息
    + 自动检测,创建和删除索引和唯一索引,这个检测是根据索引的一个或多个字段名,而不根据索引名称.因此这里需要注意,如果在一个有大量数据的表中引入新的索引,数据库可能需要一定的时间来建立索引.
    + 自动转换varchar字段类型到text字段类型,自动警告其它字段类型在模型和数据库之间不一致的情况.
    + 自动警告字段的默认值,是否为空信息在模型和数据库之间不匹配的情况
    
    这是同步数据库的最简单方法,但由于其隐藏了大量细节,并不建议使用.通常这个用来代替`CreateTables`操作.

+ 判断表是否存在`db.IsTableExist(...string|*struct) (bool,error)`参数为一个或多个空的对应Struct的指针或表名字符串.

+ 判断表是否是空的`db.IsTableEmpty(...string|*struct) (bool,error)`参数为一个或多个空的对应Struct的指针或表名字符串.

+ 创建表`db.CreateTables(...string|*struct) error`参数为一个或多个空的对应Struct的指针或表名字符串

+ 删除表`db.DropTables(...string|*struct) error`参数为一个或多个空的对应Struct的指针或表名字符串

```golang
func sync_table() {
    db, err := xorm.NewEngine("postgres","postgres://postgres:postgres@localhost:5432/test?sslmode=disable")
    if err != nil {
        fmt.Printf("%v", err)
    } else {
        defer db.Close()
        db.ShowSQL(true)
        var ok bool
        ok, err = db.IsTableExist("goods")
        if err != nil{
            fmt.Println("table goods IsTableExist error", err)
        }else{
            fmt.Println("table goods is exist :%v", ok)
            err = db.Sync2(&Goods{})
            if err != nil{
                fmt.Println("table goods Sync2 error", err)
            }else{
                ok, err = db.IsTableExist("goods")
                if err != nil{
                    fmt.Println("table goods IsTableExist error", err)
                }else{
                    fmt.Println("table goods is exist :%v", ok)
                }
            }
                
        }  
    }
}
sync_table()
```

### 写操作

> 插入数据

插入数据可以一次插入一条也可以一次插入多条

+ 插入一条

```golang
goods := &Goods{
		Id:    0,
		Price: 4,
	}

affected, err := db.Insert(goods)

```

+ 插入多条

```golang
goods := []*Goods{{
		Id:    1,
		Price: 3,
	}, {
		Id:    2,
		Price: 2,
	}, {
		Id:    3,
		Price: 1,
	},
	}

affected, err := db.Insert(goods)

```

> 更新数据

更新数据有两种方式:

+ 指定列

```golang
good := new(Goods)
good.Price = 15
affected, err := db.Id(1).Cols("price").Update(good)
```

+ 指定表后传入map,在map中指定要更新的列

```golang
affected, err := db.Table(new(Goods)).Id(2).Update(map[string]interface{}{"price":20})
```

> 删除数据

直接使用`Delete`接口就好

```golang
affected, err := db.Where("id >?", 2).Delete(new(Goods))
```

> 直接执行sql语句

```golang
sql = "update goods set price=? where id=?"
res, err := engine.Exec(sql, 40, 1) 
```

### 读操作

> 查询数据

xorm有几个方法用于查询数据,这几个方法需要在语句最后调用作为一个链式查询的终结.他们是:

+ `Get` 查询单条数据,返回的是bool值,找到的对象会被存入传入的对象: 

```golang
good := new(Goods)
has, err := db.Id(id).Get(good)
```

+ `Exist` 查询记录是否存在,它的性能比`Get`好: 

```golang
good := new(Goods)
has, err := db.Exist(goodExist)
```

+ `Find`查询多条数据使用,第一个参数为slice的指针或Map指针,即为查询后返回的结果,如果是map的指针,那么map的key为数据的主键,这种方式无法使用复合主键;第二个参数可选,为查询的条件struct的指针: 

```golang
goods := make([]Goods,0)
err := db.Find(&goods)
```

+ `Iterate`提供逐条执行查询到的记录的方法,他所能使用的条件和Find方法完全相同,但最后一位参数为一个回调函数`func (int,interface{})error`,回调函数的第二个参数就是一条数据,但需要将对象转换成对应的类型:

```golang
err := db.Where("id >?", 1).Iterate(new(Goods), func(i int, bean interface{}) error {
		good := bean.(*Goods)
		fmt.Println("good price:", good.Price)
		return nil
	})
```

+ `Rows`Iterate方法类似,提供逐条执行查询到的记录的方法,不过Rows更加灵活好用:

```golang
good := new(Goods)
rows, err := engine.Where("id >?", 1).Rows(good)
if err != nil {
}
defer rows.Close()
for rows.Next() {
    err = rows.Scan(good)
    if err != nil{
        fmt.Println("error:",err)
    }else{
        fmt.Println("good price:",good.Price)
    }
}
```

+ `Count`统计数据数量使用,Count方法的参数为struct的指针并且成为查询条件.

```golang
count, err := db.Where("id >?", 1).Count(new(Goods))
```


+ `FindAndCount` 结合Count和Find操作

```golang
goods := make([]Goods,0)
counts, err := engine.FindAndCount(&goods)
```

+ `Sum`/`SumInt`用于求和

```golang
total, err := db.Where("id >?", 1).Sum(new(Goods), "price")
```

这些方法之间都可以插入一些条件语句用于确定查询范围,比如`Where`,`Join`,或者直接使用`Sql`直接使用sql语句筛选.主要是:


+ `SQL(string, …interface{})`直接写出要执行的sql语句

+ `Select(string)`直接写出sql语句中select部分的字符串

+ `Where(string, …interface{})`直接写出sql语句中where部分的字符串

+ `Join(string,interface{},string)`;连接两张表
+ `GroupBy(string)`聚合操作字符串
+ `Having(string)` having操作的字符串

+ `AllCols()`/`Cols(…string)`/`Omit(…string)`指定全部列/指定特定的列/指定除特定列外的所有列

+ `Table(nameOrStructPtr interface{})`指定表
+ `Alias(string)`为表取别名,用于在后面的条件中表达

+ `OrderBy(string)`以字段作为排序的条件
+ `Asc(…string)`/`Desc(…string)`结果正序/逆序排列

+ `Id()`限定主键作为查询条件

+ `And(string, …interface{})`表达并列条件
+ `Or(interface{}, …interface{})`表达或条件

+ `Limit(int, …int)`/`Top(int)`限制结果数量
+ `Distinct(…string)`去重
+ `In(string, …interface{})`字段取值所在范围


> 直接执行sql语句查询数据

xorm提供了不依赖于定义好的结构体直接查询sql语句的接口`Query(sql)`,其返回值是`[]map[string][]byte`类型我们需要按需要将其值转化成我们可以用的值

```golang
sql := "select * from goods"
results, err := db.Query(sql)
if err != nil {
    fmt.Println("error:", err)
} else {
    for _, v := range results {
        fmt.Printf("good id:%v;\nprice:%v\n", string(v["id"]), string(v["price"]))
    }

}

```

类似的还有接口

+ `QueryString`,其返回值为`[]map[string]string`
+ `QueryInterface`,其返回值为`[]map[string]interface{}`

### 简单事务支持

xorm支持简单事务操作,使用`Session`会话对象代替数据库客户端对象,并配合上事务特有接口即可

+ `Begin()`开始事务
+ `Rollback()`事务回退
+ `Commit()`提交事务

每次会话创建后需要使用`Close()`接口关闭,一个更好的方式是事务的创建和定义直接使用函数包装,借助`defer`关键字关闭会话防止业务复杂了遗漏关闭操作.
需要注意一次事务必须定义在一个协程中.

```golang
func simple_transact(db *xorm.Engine) {
	fmt.Println("--------------------simple transact---------------------")
	session := db.NewSession()
	defer session.Close()
	// add Begin() before any action
	err := session.Begin()
	good := Goods{Id: 5, Price: 70}
	_, err = session.Insert(&good)
	if err != nil {
		session.Rollback()
		fmt.Println("error:", err)
		return
	}
	good2 := Goods{Price: 40}
	_, err = session.Where("id = ?", 2).Update(&good2)
	if err != nil {
		session.Rollback()
		fmt.Println("error:", err)
		return
	}

	_, err = session.Exec("delete from goods where id = ?", 2)
	if err != nil {
		session.Rollback()
		fmt.Println("error:", err)
		return
	}

	// add Commit() after all actions
	err = session.Commit()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

```

### 对业务支持

> 一致性缓存

xorm可以使用服务内存来实现一致性缓存功能,不过默认并没有开启,要开启缓存，需要在db创建完后进行配置,

```
cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
db.SetDefaultCacher(cacher)
```
这个缓存的设置为使用服务内存创建,并且每个struc(表)缓存1000条数据.要为不同表设定不同的缓存可以使用

```golang
db.MapCacher(&good, cacher)
```

不过需要特别注意不适用缓存或者需要手动编码的地方:

1. 当使用了Distinct,Having,GroupBy方法将不会使用缓存

2. 在Get或者Find时使用了Cols,Omit方法,则在开启缓存后此方法无效,系统仍旧会取出这个表中的所有字段.

3. 在使用Exec方法执行了方法之后可能会导致缓存与数据库不一致的地方.因此如果启用缓存应尽量避免使用Exec.如果必须使用,则需要在使用了Exec之后调用ClearCache手动做缓存清除的工作.比如
```golang
db.Exec("update user set name = ? where id = ?", "xlw", 1)
db.ClearCache(new(User))
```

使用服务内存作为缓存是最简单的方式,但它不利于多节点共享缓存数据,我们可以使用官方的扩展包[github.com/go-xorm/cachestore](https://github.com/go-xorm/cachestore/blob/master/README_zh-CN.md)来利用redis实现缓存

```golang

configs := map[string]string{
    "conn":"localhost:6379",
    "key":"default", // the collection name of redis for cache adapter.
}
ccStore := cachestore.NewRedisCache(configs)
cacher := xorm.NewLRUCacher(ccStore, 99999999)
db.SetDefaultCacher(cacher)
```


> 操作钩子

xorm有两种方式定义钩子:

1. 在为结构体定义对应名字的实例方法,这种方式适合用于定义每次都需要做的操作

    可以使用的钩子名有:

    + `BeforeInsert()`在将此实例插入到数据库之前执行
    + `AfterInsert()`在将此实例成功插入到数据库之后执行
    + `BeforeUpdate()`在将此实例更新到数据库之前执行
    + `AfterUpdate()`在将此实例成功更新到数据库之后执行
    + `BeforeDelete()`在将此实例对应的条件数据从数据库删除之前执行
    + `AfterDelete()`在将此实例对应的条件数据成功从数据库删除之后执行
    + `BeforeSet(name string, cell xorm.Cell)`在 Get 或 Find 方法中，当数据已经从数据库查询出来,而在设置到结构体之前调用,name为数据库字段名称,cell为数据库中的字段值.
    + `AfterSet(name string, cell xorm.Cell)`在 Get 或 Find 方法中，当数据已经从数据库查询出来,而在设置到结构体之后调用,name为数据库字段名称,cell为数据库中的字段值.
    
    ```golang
    type Goods struct {
        Id    int `xorm:"not null pk autoincr INTEGER"`
        Price uint
    }

    func (self *Goods) BeforeInsert() {
        fmt.Println("before insert good %v", self.Id)
    }

    func (self *Goods) AfterInsert() {
        fmt.Println("after insert good %v", self.Id)
    }

    func (self *Goods) BeforeUpdate() {
        fmt.Println("before update good %v", self.Id)
    }

    func (self *Goods) AfterUpdate() {
        fmt.Println("after update good %v", self.Id)
    }

    func (self *Goods) BeforeDelete() {
        fmt.Println("after delete good %v", self.Id)
    }
    func (self *Goods) AfterDelete() {
        fmt.Println("after delete good %v", self.Id)
    }

    func (self *Goods) BeforeSet(name string, cell xorm.Cell) {
        fmt.Println("before set %v as %v", name, *cell)
    }
    func (self *Goods) AfterSet(name string, cell xorm.Cell) {
        fmt.Println("after set %v as %v", name, *cell)
    }
    ```

2. 在执行语句前使用`Before(beforeFunc interface{})`/`After(afterFunc interface{})`定义操作,其回调函数的签名为`func(bean interface{})`也就是获得对象的实例.

    ```golang
    has, err := db.After(func(bean interface{}) {
        temp := bean.(*Goods)
        fmt.Println("after get select get table instance: ", temp)
    }).Id(3).Get(good)
    ```

## 时序数据库influxdb

[influxdb](https://docs.influxdata.com/influxdb/v1.7/)本身就是golang写的,它有官方的go语言客户端[github.com/influxdata/influxdb1-client/v2](https://github.com/influxdata/influxdb1-client)

他的引用方式比较奇葩,这主要是go mod系统的锅:

```golang
import(
    _ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
    influx "github.com/influxdata/influxdb1-client/v2"
)
```

influxdb只有3个操作:

1. 连接服务

```golang
client, err := influx.NewHTTPClient(client.HTTPConfig{
    Addr: "http://localhost:8086",
})
if err != nil {
    fmt.Println("Error creating InfluxDB Client: ", err.Error())
}
defer c.Close()
```

2. 发出请求

```golang
func query_point() {
	...
	q := influx.NewQuery("SELECT count(*) FROM cpu_usage", "BumbleBeeTuna", "s")
	if response, err := client.Query(q); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
}
```

3. 写入数据

```golang
func randomWrite() {
    ...
    bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
        Database:  "BumbleBeeTuna",
        Precision: "s",
    })
    tags := map[string]string{"cpu": "cpu-total"}
    fields := map[string]interface{}{
        "idle":   10.1,
        "system": 53.3,
        "user":   46.6,
    }
    pt, err := influx.NewPoint("cpu_usage", tags, fields, time.Now())
    if err != nil {
        fmt.Println("Error: ", err.Error())
    }
    bp.AddPoint(pt)
    err = client.Write(bp)
    if err != nil {
        fmt.Println("Error: ", err.Error())
    }
}
```

本例代码在[这里](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/influxdb_test)

## 最常用的缓存Redis

实际上分布式缓存有很多选择,但恐怕最有通用性的就是redis了,我们一般使用[github.com/go-redis/redis](https://github.com/go-redis/redis)这个库来连接redis.

本部分例子在[这里](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/redis_test)

### 建立连接

通常redis有两种服务端:

1. 单机`NewClient(*redis.Options)`

```golang
opt, err  := redis.ParseURL("redis://localhost:6379/1")
if err != nil {
   panic(err)
}
client:=redis.NewClient(opt)
```

2. 集群`NewClusterClient(*redis.ClusterOptions)`

```golang
redisdb := redis.NewClusterClient(&redis.ClusterOptions{
	Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
})
```
本文将以单机为例,因为像亚马逊,阿里云这些服务提供商已经将集群封装成了单机接口.


### 使用命令

这个包已经封装了几乎全部的redis命令.

```golang
val, err := client.Get(key).Result()
if err != nil {
    if err == redis.Nil {
        fmt.Printf("key %v 不存在\n", key)
    } else {
        fmt.Println("error:", err)
    }
} else {
    fmt.Println("key", key, val)
}
```

如果不确定要用的命令是否已经有封装,可以直接使用`Do()`接口直接发出命令


```golang
val, err := client.Do("get", key).String()
```

`,`用于分隔命令中各段

### 事务

redis使用`pipline`来定义事务.

```golang
func incr_pipeline(client *redis.Client) {
	pipe := client.Pipeline()

	incr := pipe.Incr("pipeline_counter")
	pipe.Expire("pipeline_counter", time.Hour)
	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)

}
```

## 常见消息中间件

消息中间件一般就是两个功能

1. 队列(生产消费模式)
2. 广播(广播模式)

队列一般用于生产消费模式,用于将并发的任务串行化;广播则用于在范围内同时通知多个组件,常用于并行化处理数据.

接下来的例子中本文将使用下面3种工具实现一个分发随机数求平方和的功能.其中会用到这两种模式:

+ 生产消费模式: 生产者向sourceQ队列发送数据,消费者从sourceQ队列取数据,消费者计算完成平方后将结果放入队列resultQ,生产者接收resultQ队列中的结果更新累加结果并打印在标准输出中.
+ 广播模式: 生产者在收到KeyboardInterrupt错误时向频道exitCh发出消息,消费者订阅频道exitCh,当收到消息时退出.

### 使用redis做消息中间件

redis因为数双端列表和pub/sub模式,而且实时性非常好,所以在允许信息丢失的情况下经常有人用它做消息中间件,比如著名的任务队列工具celery及其衍生工具就常用redis做broker.

代码[在这里](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/broker_test/redis_broker)

> 双端列表做消息队列

对应的操作包括

+ `EXISTS key`判断key是否存在
+ `TYPE key`用于获取key对应数据结构,必须得是`list`
+ `lpush key values`将消息从左侧推入key对应的list中
+ `lpushx key values`将消息从左侧推入key对应的一个已经存在的list中
+ `rpop key`从右侧取出key对应list中第一个值

> pub/sub做广播

对应的操作

+ 发布广播到信道

```golang
err := redisdb.Publish("mychannel1", "hello").Err()
```

+ 订阅监听

```golang
pubsub := client.Subscribe("mychannel1")

ch := pubsub.Channel()

// Consume messages.
for msg := range ch {
    fmt.Println(msg.Channel, msg.Payload)
}

```

> 例子实现:

+ 共用的推送代码

```golang
package push

import (
	"fmt"

	"github.com/go-redis/redis"
)

func Push(client *redis.Client, Q string, value string) error {
	isExists, err1 := client.Exists(Q).Result()
	if err1 != nil {
		if err1 == redis.Nil {
			fmt.Printf("key %v 不存在\n", Q)
		} else {
			fmt.Println("error:", err1)
		}
		return err1
	}

	type_, err2 := client.Type(Q).Result()
	if err2 != nil {
		if err2 == redis.Nil {
			fmt.Printf("queue %v 不存在\n", Q)
		} else {
			fmt.Println("error:", err2)
		}
		return err2
	}
	if isExists > 0 && type_ == "list" {
		_, err := client.LPushX(Q, value).Result()
		if err != nil {
			if err == redis.Nil {
				fmt.Printf("queue %v 不存在\n", Q)
			} else {
				fmt.Println("error:", err)
			}
			return err
		}
		fmt.Printf("send %v to %v\n", value, Q)
		return nil
	}
	_, err := client.Del(Q).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("key %v 不存在\n", Q)
		} else {
			fmt.Println("error:", err)
		}
		return err
	}
	_, err = client.LPush(Q, value).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("key %v 不存在\n", Q)
		} else {
			fmt.Println("error:", err)
		}
		return err
	}
	fmt.Printf("send %v to %v\n", value, Q)
	return nil
}
```

将消息推送到队列中操作在生产者和消费者两端都是一样的,需要改变的只是队列的名字和要推送的值而已,因此将其抽象出来单独作为一个子模块.


+ 生产者

```golang
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"redis_test/push"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var c chan os.Signal
var c1 chan os.Signal
var wg sync.WaitGroup

func producer(client *redis.Client, sourceQ string, exitCh string) {
Loop:
	for {
		select {
		case s := <-c:
			fmt.Println()
			fmt.Println("Producer | get exit signal", s)
			client.Publish(exitCh, "Exit")
			break Loop
		default:
		}
		data := rand.Int31n(400)
		err := push.Push(client, sourceQ, strconv.Itoa(int(data)))
		if err != nil {
			fmt.Println("err:", err)
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Producer |  exit")
	wg.Done()
}

func collector(client *redis.Client, resultQ string) error {
	var sum int64 = 0
Loop:
	for {
		select {
		case s := <-c1:
			fmt.Println()
			fmt.Println("collector | get exit signal", s)
			break Loop
		default:
		}
		isExists, err1 := client.Exists(resultQ).Result()
		if err1 != nil {
			if err1 == redis.Nil {
				fmt.Printf("key %v 不存在\n", resultQ)
			} else {
				fmt.Println("collector error:", err1)
			}
			return err1
		}
		type_, err2 := client.Type(resultQ).Result()
		if err2 != nil {
			if err2 == redis.Nil {
				fmt.Printf("key %v 不存在\n", resultQ)
			} else {
				fmt.Println("collector error:", err2)
			}
			return err2
		}
		if isExists > 0 && type_ == "list" {
			data, err3 := client.RPop(resultQ).Result()
			if err3 != nil {
				if err3 == redis.Nil {
					fmt.Printf("key %v 不存在\n", resultQ)
					time.Sleep(1 * time.Second)
				} else {
					fmt.Println("collector error:", err3)
				}
				return err3
			}
			fmt.Println("collector received data: ", data)
			d, err := strconv.ParseInt(data, 10, 64)
			if err != nil {
				fmt.Println("collector err:", err)
			}
			sum += d
			fmt.Println("collector get sum ", sum)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("collector | exit")
	wg.Done()
	return nil
}

func main() {
	opt, err := redis.ParseURL("redis://localhost:6379/1")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	sourceQ := "sourceQ"
	resultQ := "resultQ"
	exitCh := "exitCh"

	c = make(chan os.Signal, 1)
	c1 = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	signal.Notify(c1, os.Interrupt, os.Kill)
	wg.Add(1)
	go collector(client, resultQ)
	wg.Add(1)
	go producer(client, sourceQ, exitCh)
	wg.Wait()
}
```

生产者部分我们构造两个协程--生产器`producer`用于每隔1s向队列中发送要处理的数据,和收集器`collector`用于收集运算结果.同时使用`sync.WaitGroup`阻塞主go协程,同时等待两个go协程执行完毕退出.于此同时主协程监听退出信号(ctr+c),利用channel向两个go协程发送退出指令.当有退出信号时,生产器会将向redis中的退出信道发出一条广播消息.

+ 消费者

```golang
package main

import (
	"fmt"
	"redis_test/push"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func get_result(client *redis.Client, sourceQ string, resultQ string) error {
	fmt.Printf("get_result running\n")

	for {
		isExists, err1 := client.Exists(sourceQ).Result()
		if err1 != nil {
			if err1 == redis.Nil {
				fmt.Printf("key %v 不存在\n", sourceQ)
			} else {
				fmt.Println("error:", err1)
			}
			return err1
		}
		type_, err2 := client.Type(sourceQ).Result()
		if err2 != nil {
			if err2 == redis.Nil {
				fmt.Printf("key %v 不存在\n", sourceQ)
			} else {
				fmt.Println("error:", err2)
			}
			return err2
		}
		if isExists > 0 && type_ == "list" {
			data, err3 := client.RPop(sourceQ).Result()
			if err3 != nil {
				if err3 == redis.Nil {
					fmt.Printf("key %v 不存在\n", resultQ)
					time.Sleep(1 * time.Second)
				} else {
					fmt.Println("error:", err3)
				}
				return err3
			}
			fmt.Println("received data: ", data)
			d, err := strconv.ParseInt(data, 10, 32)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				result := d * d
				push.Push(client, resultQ, strconv.Itoa(int(result)))
			}
		}
	}
}

func main() {
	opt, err := redis.ParseURL("redis://localhost:6379/1")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	sourceQ := "sourceQ"
	resultQ := "resultQ"
	exitCh := "exitCh"
	go get_result(client, sourceQ, resultQ)
	pubsub := client.Subscribe(exitCh)
	defer pubsub.Close()
	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		if msg.Payload == "Exit" {
			return
		}
	}
}

```

消费者端只需要定时监听资源队列,有新数据就处理,处理完后丢到结果队列即可.

### 使用kafka做消息中间件


kafka是一个追求高吞吐的分布式消息队列.和redis比较是另一个极端,几乎是为复杂而生:

+ 天生支持分布式,并且它必须依赖zonekeeper维护集群一致性.
+ 天生持久化,硬盘允许情况下保留所有消息.

kafka使用groupid来区分监听端是一次性消耗还是广播,当监听端使用不同的groupid时它相当于做广播,而使用相同的groupid时它相当于做生产消费的队列.

我们使用[gopkg.in/confluentinc/confluent-kafka-go.v1/kafka](https://github.com/confluentinc/confluent-kafka-go)包来连接kafka,注意这个包无法在windows下安装使用


代码[在这里](https://github.com/hsz1273327/TutorialForGoLang/tree/master/%E4%BD%BF%E7%94%A8Golang%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/code/broker_test/kafka_broker)
其操作是:

+ 生产端

```golang
p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

doneChan := make(chan bool)

go func() {
    defer close(doneChan)
    for e := range p.Events() {
        switch ev := e.(type) {
        case *kafka.Message:
            m := ev
            if m.TopicPartition.Error != nil {
                fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
            } else {
                fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
                    *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
            }
            return

        default:
            fmt.Printf("Ignored event: %s\n", ev)
        }
    }
}()

value := "Hello Go!"
p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}

// wait for delivery report goroutine to finish
_ = <-doneChan

p.Close()
```

生产端由于kafka有消息确认机制,因此需要一个channel来获取被确认的消息,上面例子中`p.Events()`就是这样一个channel,通常我们发送归发送,监控消息发送则是另起一个go携程来做,发送我们就是直接向`p.ProduceChannel()`中丢入一个`*kafka.Message`而已.


+ 消费端

```golang
c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		// Enable generation of PartitionEOF when the
		// end of a partition is reached.
		"enable.partition.eof": true,
		"auto.offset.reset":    "earliest"})
        
err = c.SubscribeTopics(topics, nil)


for run == true {
    select {
    case sig := <-sigchan:
        fmt.Printf("Caught signal %v: terminating\n", sig)
        run = false

    default:
        ev :=c.Poll(100)
        if ev == nil {
				continue
			}
        switch e := ev.(type) {
        case kafka.AssignedPartitions:
            fmt.Fprintf(os.Stderr, "%% %v\n", e)
            c.Assign(e.Partitions)
        case kafka.RevokedPartitions:
            fmt.Fprintf(os.Stderr, "%% %v\n", e)
            c.Unassign()
        case *kafka.Message:
            fmt.Printf("%% Message on %s:\n%s\n",
                e.TopicPartition, string(e.Value))
        case kafka.PartitionEOF:
            fmt.Printf("%% Reached %v\n", e)
        case kafka.Error:
            // Errors should generally be considered as informational, the client will try to automatically recover
            fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
        }
    }
}

fmt.Printf("Closing consumer\n")
c.Close()
```

消费端一样是监听`c.Poll(100)`,根据其类型进行不同的处理

> 例子实现:

+ 生产者

```golang
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var csend chan os.Signal
var cprod chan os.Signal
var ccoll chan os.Signal
var wg sync.WaitGroup

func secondRef(p *kafka.Producer) {
	fmt.Println("secondRef | start")
Loop:
	for {
		select {
		case s := <-csend:
			fmt.Println()
			fmt.Println("secondRef | get exit signal", s)
			break Loop
		case e := <-p.Events():
			fmt.Println("secondRef | get event")
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		default:
		}
	}
	fmt.Println("secondRef |  exit")
	wg.Done()
}

func producer(p *kafka.Producer, sourceQ string, exitCh string) {
Loop:
	for {
		select {
		case s := <-cprod:
			fmt.Println()
			fmt.Println("Producer | get exit signal", s)
			p.ProduceChannel() <- &kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &exitCh,
					Partition: kafka.PartitionAny,
				},
				Value: []byte("Exit")}
			time.Sleep(1 * time.Second)
			fmt.Println("Producer send msg exit to exitch")
			break Loop
		default:
		}
		data := rand.Int31n(400)
		p.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &sourceQ,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(strconv.Itoa(int(data)))}
		fmt.Println("Producer send msg ", data, " to sourceQ")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Producer |  exit")
	wg.Done()
}

func collector(c *kafka.Consumer, resultQ string) error {
	err := c.SubscribeTopics([]string{resultQ}, nil)
	if err != nil {
		fmt.Println("collector | err", err)
		return err
	}
	var sum int64 = 0
Loop:
	for {
		select {
		case s := <-ccoll:
			fmt.Println()
			fmt.Println("collector | get exit signal", s)
			break Loop

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				value := string(e.Value)
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition)
				d, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					fmt.Println("collector err:", err)
				}
				sum += d
				fmt.Println("collector get sum ", sum)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}
	fmt.Println("collector | exit")
	wg.Done()
	return nil
}

func main() {
	sourceQ := "sourceQ"
	resultQ := "resultQ"
	exitCh := "exitCh"
	broker := "localhost:9092"
	producerConf := kafka.ConfigMap{
		"bootstrap.servers": broker,
	}
	consumerConf := kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        resultQ,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            true,
		"auto.offset.reset":               "latest"}
	kafkaProducer, err1 := kafka.NewProducer(&producerConf)
	if err1 != nil {
		panic(err1)
	}
	kafkaConsumer, err2 := kafka.NewConsumer(&consumerConf)
	if err2 != nil {
		panic(err2)
	}
	csend = make(chan os.Signal, 1)
	cprod = make(chan os.Signal, 1)
	ccoll = make(chan os.Signal, 1)
	signal.Notify(ccoll, os.Interrupt, os.Kill)
	signal.Notify(cprod, os.Interrupt, os.Kill)
	signal.Notify(csend, os.Interrupt, os.Kill)
	wg.Add(1)
	go collector(kafkaConsumer, resultQ)
	wg.Add(1)
	go producer(kafkaProducer, sourceQ, exitCh)
	wg.Add(1)
	go secondRef(kafkaProducer)
	wg.Wait()
}
```

+ 消费者

```golang
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var csend chan os.Signal
var cexit chan os.Signal
var cget chan os.Signal
var wg sync.WaitGroup

func secondRef(p *kafka.Producer) {
	fmt.Println("secondRef | start")
Loop:
	for {
		select {
		case s := <-csend:
			fmt.Println()
			fmt.Println("secondRef | get exit signal", s)
			break Loop
		case e := <-p.Events():
			fmt.Println("secondRef | get event")
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		default:
		}
	}
	fmt.Println("secondRef |  exit")
	wg.Done()
}

func get_result(c *kafka.Consumer, p *kafka.Producer, sourceQ string, resultQ string) error {
	fmt.Printf("get_result running\n")
	err := c.SubscribeTopics([]string{sourceQ}, nil)
	if err != nil {
		fmt.Println("collector | err", err)
		return err
	}
Loop:
	for {
		select {
		case s := <-cget:
			fmt.Println()
			fmt.Println("get_result | get exit signal", s)
			break Loop

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				value := string(e.Value)
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition)
				d, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					fmt.Println("get_result err:", err)
				}
				result := d * d
				p.ProduceChannel() <- &kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &resultQ,
						Partition: kafka.PartitionAny,
					},
					Value: []byte(strconv.Itoa(int(result)))}
				fmt.Println("get_result send msg ", result)
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}
	fmt.Println("get_result |  exit")
	wg.Done()
	return nil
}
func listen_exit(c *kafka.Consumer, exitCh string) error {
	fmt.Printf("listen_exit running\n")
	err := c.SubscribeTopics([]string{exitCh}, nil)
	if err != nil {
		fmt.Println("listen_exit | err", err)
		return err
	}
Loop:
	for {
		select {
		case s := <-cexit:
			fmt.Println()
			fmt.Println("listen_exit | get exit signal", s)
			break Loop
		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				fmt.Println("listen_exit | get exit signal from remote----------------")
				value := string(e.Value)
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition)
				if value == "Exit" {
					fmt.Println("get exit")
					cget <- os.Interrupt
					csend <- os.Interrupt
					break Loop
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}
	fmt.Println("listen_exit |  exit")
	wg.Done()
	return nil
}
func main() {
	broker := "localhost:9092"
	sourceQ := "sourceQ"
	resultQ := "resultQ"
	exitCh := "exitCh"
	producerConf := kafka.ConfigMap{
		"bootstrap.servers": broker,
	}
	consumerConf := kafka.ConfigMap{
		"bootstrap.servers":               broker,
		"group.id":                        resultQ,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"enable.partition.eof":            true,
		"auto.offset.reset":               "latest"}
	kafkaProducer, err1 := kafka.NewProducer(&producerConf)
	if err1 != nil {
		panic(err1)
	}
	kafkaConsumer, err2 := kafka.NewConsumer(&consumerConf)
	if err2 != nil {
		panic(err2)
	}
	kafkaExitConsumer, err2 := kafka.NewConsumer(&consumerConf)
	if err2 != nil {
		panic(err2)
	}
	csend = make(chan os.Signal, 1)
	cget = make(chan os.Signal, 1)
	cexit = make(chan os.Signal, 1)
	signal.Notify(csend, os.Interrupt, os.Kill)
	signal.Notify(cget, os.Interrupt, os.Kill)
	signal.Notify(cexit, os.Interrupt, os.Kill)
	wg.Add(1)
	go get_result(kafkaConsumer, kafkaProducer, sourceQ, resultQ)
	wg.Add(1)
	go secondRef(kafkaProducer)
	wg.Add(1)
	go listen_exit(kafkaExitConsumer, exitCh)
	wg.Wait()
}

```

需要注意监听退出消息的消费者和监听消息的消费者不能是同一个


```go

```
