package handle

import "github.com/space-ark-z/infra-common/pkg/log"

func Handle(err error) {
	if err != nil {
		log.Error(map[string]any{
			"error": err,
		})
	}
}

func Fatal(err error) {
	if err != nil {
		log.Fatal(map[string]any{
			"error": err,
		})
	}
}
