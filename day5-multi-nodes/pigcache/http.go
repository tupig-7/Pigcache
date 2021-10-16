package pigcache

import (
	"Pigcache/day5-multi-nodes/pigcache/consistenthash"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_pigcache/"
	defaultReplicas = 50 // 副本个数
)

// HTTPPool implements PeerPicker for a pool of HTTP peers
// HTTPPool 从一个HTTP池中获取一个节点
type HTTPPool struct {
	// this peer's base URL, e.g."https://example.net:8080"
	self string
	basePath string

	mu sync.Mutex // guards peers and httpGetters
	peers *consistenthash.Map
	httpGetters map[string]*httpGetter // keyed by e.g.. "http://10.0.0.2:8008"
}

// NewHTTPPool initializes an HTTP pool of peers
// NewHTTPPool 初始化一个HTTP节点池
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
// Log 记录服务器日志
func (p *HTTPPool) Log(format string, v ...interface{})  {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

//ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: " + groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream") // 只能提交二进制， only one binary
	w.Write(view.ByteSlice())
}

// Set updates the pool's list of peers
func (p *HTTPPool) Set(peers ...string)  {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))

}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key),
		) //QueryEscape函数对s进行转码使之可以安全的用在URL查询里
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}
var _ PeerGetter = (*httpGetter)(nil) // 检查有无实现接口