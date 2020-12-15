module HomeLab/RSSProxy

go 1.14

require github.com/go-resty/resty/v2 v2.3.0

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.7.1
	github.com/t-tomalak/logrus-easy-formatter v0.0.0-20190827215021-c074f06c5816
)
