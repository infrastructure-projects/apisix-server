package appctx

import (
	"encoding/json"
	"flag"
	"github.com/apache/apisix-go-plugin-runner/internal/eventbus"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

const keyFilePath = "filePath"

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	pflag.String(keyFilePath, "/config/application.yaml", "配置文件的绝对路径")
	viper.SetConfigType("yaml")
	loadConfigs()
}

func loadConfigs() {
	filePath := viper.GetString(keyFilePath)
	refreshProperties(filePath)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("[AppContext]无法监控文件事件.error=%s", err.Error())
		return
	}
	go func(path string) {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Infof("[AppContext]检测到文件发生变化")
					refreshProperties(path)
					eventbus.Publish(EventConfigChange, event.Name)
				}
			case e := <-watcher.Errors:
				if e != nil {
					log.Infof("[AppContext]监控文件出错.error=%s", e.Error())
				}
			}
		}
	}(filePath)
	if err = watcher.Add(filePath); err != nil {
		log.Infof("[Properties]添加文件监控失败.error=%s", err.Error())
	}
}

func refreshProperties(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Errorf("[Properties]无法打开配置文件.error=%s", err.Error())
		return
	}
	if err = viper.MergeConfig(file); err != nil {
		log.Errorf("[Properties]无法加载配置.error=%s", err.Error())
	}
	settings, _ := json.Marshal(viper.AllSettings())
	log.Infof("[Properties]已刷新配置.settings:%s", string(settings))
}
