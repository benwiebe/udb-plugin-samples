package datasources

import (
	"context"
	"time"
)

// CurrentTimeDatasource Is a trivial dummy datasource to demonstrate how to
// implement a UDB Datasource. It will return the current time as its data.
// This may be a useful datasource for testing or just as an example.
type CurrentTimeDatasource struct{}

func (c *CurrentTimeDatasource) GetId() string {
	return "UdbSamplePlugin/CurrentTime"
}

func (c *CurrentTimeDatasource) GetName() string {
	return "Current Time"
}

func (c *CurrentTimeDatasource) GetType() string {
	return "UdbSamplePlugin/CurrentTime"
}

func (c *CurrentTimeDatasource) GetData() any {
	return time.Now()
}

func (c *CurrentTimeDatasource) Start(_ context.Context) error { return nil }

// DataChanged returns nil because time.Now() is always fresh — no push notification needed.
func (c *CurrentTimeDatasource) DataChanged() <-chan struct{} { return nil }
