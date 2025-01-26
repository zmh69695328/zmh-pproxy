package internal

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vulcand/oxy/v2/forward"
	"github.com/vulcand/oxy/v2/roundrobin"
	"github.com/vulcand/oxy/v2/testutils"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

var (
	// 目标服务器列表
	targetServers []string
	// 用于保护 targetServers 的读写锁
	mu sync.RWMutex
	// 创建一个转发器
	forwarder *httputil.ReverseProxy
	// 负载均衡器
	lb *roundrobin.RoundRobin
)

// 检查目标服务器健康状态的函数
func checkServerHealth(server string) bool {
	// 向目标服务器发送请求，查看状态码
	resp, err := http.Get(server + "/health") // 假设目标服务器提供 /health 路径来进行健康检查
	if err != nil {
		// 如果请求失败，认为该服务器不可用
		return false
	}
	defer resp.Body.Close()

	// 如果返回 404 状态码，认为是健康的
	if resp.StatusCode == http.StatusNotFound {
		return true
	}

	// 其他状态码（如 5xx、4xx 等）认为服务器不健康
	return false
}

// 动态更新负载均衡器的目标服务器
func updateLoadBalancer() {
	mu.Lock()
	defer mu.Unlock()

	// 创建负载均衡器
	forwarder := forward.New(false)
	newLB, err := roundrobin.New(forwarder)
	if err != nil {
		log.Fatal("Failed to create load balancer:", err)
	}

	// 遍历目标服务器，添加健康的服务器
	for _, server := range targetServers {
		if checkServerHealth(server) {
			err := newLB.UpsertServer(testutils.MustParseRequestURI(server))
			if err != nil {
				log.Printf("Failed to add server %s: %v", server, err)
			}
		} else {
			log.Printf("Server %s is unhealthy (status code other than 404), not adding to load balancer", server)
		}
	}

	// 更新负载均衡器
	lb = newLB
}

// 获取当前目标服务器列表
func getTargetServers() []string {
	mu.RLock()
	defer mu.RUnlock()
	return targetServers
}

func app() {
	// 初始化目标服务器列表
	targetServers = []string{
		"https://baidu.com",
		"https://google.com",
	}

	// 创建转发器并应用负载均衡器
	forwarder = forward.New(false)
	// 初始化负载均衡器
	lb, err := roundrobin.New(forwarder)
	if err != nil {
		log.Fatal("创建负载均衡器失败:", err)
	}
	for _, targetServer := range targetServers {
		err := lb.UpsertServer(testutils.MustParseRequestURI(targetServer))
		if err != nil {
			log.Fatal("使用负载均衡器失败:", err)
		}
	}
	// 创建一个路由器
	r := chi.NewRouter()
	r.Use(middleware.CleanPath)
	// 代理请求并转发到负载均衡的目标服务器
	r.Get("/proxy/*", func(w http.ResponseWriter, r *http.Request) {
		prefix := "/proxy"
		// 重写请求路径
		if strings.HasPrefix(r.URL.Path, prefix) {
			r.RequestURI = r.URL.Path[len(prefix):]
		}
		lb.ServeHTTP(w, r)
	})
	// 启动 HTTP 服务，监听 8080 端口
	fmt.Println("负载均衡代理服务器启动，监听 8080 端口...")

	// 启动代理服务器
	log.Fatal(http.ListenAndServe(":8080", r))
}
