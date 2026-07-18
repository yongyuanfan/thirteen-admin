package runtime

import (
	"fmt"
	"log"
	"thirteen-admin/internal/logging"
)

func logCoref(format string, args ...any) {
	log.Printf("[runtime] "+format, args...)
	logging.Info(fmt.Sprintf(format, args...), "component", "runtime")
}

func logServicef(name string, format string, args ...any) {
	log.Printf("[runtime] [%s] "+format, append([]any{name}, args...)...)
	logging.Info(fmt.Sprintf(format, args...), "component", "runtime", "service", name)
}
