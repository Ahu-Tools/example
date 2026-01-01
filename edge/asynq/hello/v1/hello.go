package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	//@ahum: imports
)

type WorldPayload struct {
}

//@ahum: payloads

func HandleWorld(ctx context.Context, t *asynq.Task) error {
	var p WorldPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("World executed!")
	// World code ...
	return nil
}

//@ahum: handlers
