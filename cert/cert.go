package cert

import (
	"fmt"
	"path/filepath"

	"github.com/yaoapp/gou/ssl"
	"github.com/yaoapp/kun/log"
	"github.com/yaoapp/yao/config"
	"github.com/yaoapp/yao/share"
)

// Load 加载API
func Load(cfg config.Config) error {
	var root = filepath.Join(cfg.Root, "certs")
	return LoadFrom(root, "")
}

// LoadFrom 从特定目录加载
func LoadFrom(dir string, prefix string) error {

	if share.DirNotExists(dir) {
		return fmt.Errorf("%s does not exists", dir)
	}

	err := share.Walk(dir, ".pem", func(root, filename string) {
		name := prefix + share.SpecName(root, filename)
		_, err := ssl.LoadCertificateFrom(filename, name)
		if err != nil {
			log.With(log.F{"root": root, "file": filename}).Error(err.Error())
		}
	})

	return err
}
