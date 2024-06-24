package security

import (
	"fmt"
	"github.com/apache/apisix-go-plugin-runner/internal/appctx"
	"github.com/apache/apisix-go-plugin-runner/internal/eventbus"
	apisix "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"net/http"
	"strings"
)

type Filter struct {
	BlockList struct {
		Paths   []string `mapstructure:"paths"`
		Headers []string `mapstructure:"headers"`
		XSS     []string `mapstructure:"xss"`
	}
}

func New() *Filter {
	filter := &Filter{}
	filter.reloadConfigs()
	eventbus.Subscribe(appctx.EventConfigChange, func(args ...any) {
		filter.reloadConfigs()
	})
	return filter
}

func (Self *Filter) reloadConfigs() {
	if err := appctx.Configs().UnmarshalKey("application.security", Self); err != nil {
		log.Errorf("[SecurityFilter]无法加载配置信息.error=%s", err.Error())
	}
}

func (Self *Filter) RequestFilter(conf any, w http.ResponseWriter, request apisix.Request) {
	// 阻止内部API被外部调用
	if !Self.validPathBlockList(w, request) {
		return
	}
	// 阻止不被允许的请求头
	if !Self.validHeaderBlockList(w, request) {
		return
	}
	// 阻止XSS攻击代码
	if !Self.validHeaderXssBlockList(w, request) {
		return
	}
}

func (Self *Filter) ResponseFilter(conf any, w apisix.Response) {

}

func (Self *Filter) validHeaderBlockList(w http.ResponseWriter, request apisix.Request) bool {
	headers := request.Header()
	for index := range Self.BlockList.Headers {
		key := Self.BlockList.Headers[index]
		if headers.Get(key) != "" {
			w.WriteHeader(400)
			data := fmt.Sprintf(`{"code": "400", "message": "invalid value in request Headers: %s"}`, key)
			_, _ = w.Write([]byte(data))
			return false
		}
	}
	return true
}

func (Self *Filter) validPathBlockList(w http.ResponseWriter, request apisix.Request) bool {
	if blockList := Self.BlockList.Paths; len(blockList) > 0 {
		path := string(request.Path())
		for i := range blockList {
			if strings.Contains(path, blockList[i]) {
				w.WriteHeader(404)
				_, _ = w.Write([]byte(fmt.Sprintf("%s not found", path)))
				return false
			}
		}
	}
	return true
}

func (Self *Filter) validHeaderXssBlockList(w http.ResponseWriter, request apisix.Request) bool {
	headers := request.Header().View()
	for k, v := range headers {
		for _, block := range Self.BlockList.XSS {
			if strings.Contains(k, block) {
				Self.writeBadRequest(w)
				return false
			}
			for _, str := range v {
				if strings.Contains(str, block) {
					Self.writeBadRequest(w)
					return false
				}
			}
		}
	}
	return true
}

func (*Filter) writeBadRequest(w http.ResponseWriter) {
	w.WriteHeader(400)
	_, _ = w.Write([]byte("Bad Request"))
}

func (Self *Filter) ParseConf(in []byte) (conf any, err error) {
	return "", nil
}

func (Self *Filter) Name() string {
	return "SecurityFilter"
}
