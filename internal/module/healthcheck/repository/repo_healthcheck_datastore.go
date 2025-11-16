package healthcheckrepository

import (
	"context"
	"time"
)

func (r *repository) PingDatabase(ctx context.Context) (responseTime time.Duration, err error) {
	timeNow := time.Now().UTC()
	err = r.db.RDBMS().Ping(ctx)
	responseTime = time.Since(timeNow)
	return
}
