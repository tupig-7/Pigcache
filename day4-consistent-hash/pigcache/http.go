package pigcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_pigcache/"

// HTTPPool implements PeerPicker for a pool of HTTP peers
// HTTPPool 从一个HTTP池中获取一个节点
type HTTPPool struct {
	// this peer's base URL, e.g."https://example.net:8080"
	self string
	basePath string
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