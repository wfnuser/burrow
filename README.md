# burrow
A Distributed Cache Library for Golang. BTW, Burrow is where gophers live! Star it and gopher will be happy.😄

This is a toy project for learning how distributed cache work. It should not be used in production.

# How to use
```
func main() {
	burrow.NewBurrow("test", 5, burrow.FuncGetter(
		func(key string) (lru.Value, bool) {
			log.Println("Fetch data from datasource by: ", key)
			if v, ok := db[key]; ok {
				return v, true
			}
			return nil, false
		}))
	servers := []string{"localhost:5001", "localhost:5002", "localhost:5003"}
	for _, serverURL := range servers {
		server := burrow.NewHTTPPoolWithServers(serverURL, servers)
		go func(serverURL string) {
			http.ListenAndServe(serverURL, server)
		}(serverURL)
	}
	select {}
}
```

目前功能写的还比较粗糙，使用时需要手动在代码中注册缓存节点，手动启动各个节点的实例。且彼此没有通信，不知道对方存活情况。希望使用者在端侧管理负载均衡策略，最好可以注册所有实例，直接也用一致性哈希访问策略访问数据缓存节点。
当然缓存节点间目前做了转发的服务，即客户端也可随机访问其中一个节点，该节点如果不是对应的缓存节点，会根据一致性哈希访问对应的节点。
调用服务采用http通信，访问路由 ${server_path}/${burrow}/${namespace}/${key} 可返回对应的value。
需要手动注册数据源。

笔者一直以来是前端工程师，第一次写这种玩具项目，收获很大。在研发过程中自然而然的想到了如何负载均衡，通信是否可以自定义协议，如何摘除不可用节点，如果节点信息注册在端测如何同步等问题。大部分考虑都是出于直觉，想到其实我的考虑一定不全面，且业界应该肯定有很多优秀的经验了，所以该项目先告一段落，想等之后学习一段时间之后继续优化。当然如何让该工具易用也是一个非常重要的问题，希望有一天可以实现一个真正好用的分布式缓存而不仅仅是玩具。 哈哈哈。

# 致谢
本项目非常多的参考了 [geecache](https://github.com/geektutu/7days-golang/tree/master/gee-cache) 和 [groupcache](https://github.com/golang/groupcache) 非常感谢他们
