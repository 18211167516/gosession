package gosession

type StoreInterface interface {
	//删除某个key的值
	Remove(Sid string, key string) error
	//获取某个key的值
	Get(Sid string, key string) (interface{}, error)
	//设置某个key的值
	Set(Sid string, key string, value interface{}, ttl int) error
}
