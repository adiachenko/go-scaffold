package rabbitmq

import (
	"fmt"
	"strings"

	"adiachenko/go-scaffold/internal/platform/conf"
)

func CompileRabbitMqServerURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		conf.RabbitMqUser,
		escapeURL(conf.RabbitMqPassword),
		conf.RabbitMqHost,
		conf.RabbitMqPort,
	)
}

func escapeURL(param string) string {
	for _, separator := range []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "?"} {
		param = strings.Replace(param, separator, "\\"+separator, -1)
	}

	return param
}
