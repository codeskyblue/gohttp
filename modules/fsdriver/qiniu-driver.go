package fsdriver

import (
	qiniudriv "github.com/goftp/qiniu-driver"
)

func NewQiniuDriver(accessKey, secretKey, bucket string) (FSDriver, error) {
	factory := qiniudriv.NewQiniuDriverFactory(accessKey, secretKey, bucket)
	driver, err := factory.NewDriver()
	if err != nil {
		return nil, err
	}
	return &GoftpDriverAdapter{driver}, nil
}
