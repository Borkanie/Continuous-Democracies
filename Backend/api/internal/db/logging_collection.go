package db

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoCollection is the small subset of *mongo.Collection's API that any
// repository in this codebase actually calls (confirmed by grep: Find,
// FindOne, Aggregate, Distinct are the only four). Repositories depend on
// this interface instead of *mongo.Collection directly so
// GetLoggingCollection can transparently wrap every query with timing and
// tracing, with zero changes to repository method bodies.
type MongoCollection interface {
	Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) (*mongo.Cursor, error)
	Distinct(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult
}

// loggingCollection wraps a real *mongo.Collection, emitting one debug-level
// "mongo query" log line per call (collection, operation, duration, and the
// error where the driver surfaces it synchronously) before returning the
// real result unchanged. This is the one decorator that gives every
// repository automatic, zero-maintenance DB-layer tracing.
type loggingCollection struct {
	name string
	next *mongo.Collection
}

// GetLoggingCollection returns database's collectionName collection wrapped
// for automatic query logging. Repositories should call this instead of
// database.Collection directly.
func GetLoggingCollection(database *mongo.Database, collectionName string) MongoCollection {
	return &loggingCollection{name: collectionName, next: database.Collection(collectionName)}
}

func (collection *loggingCollection) Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error) {
	start := time.Now()
	cursor, err := collection.next.Find(ctx, filter, opts...)
	logMongoQuery(ctx, collection.name, "Find", start, err)
	return cursor, err
}

func (collection *loggingCollection) FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
	start := time.Now()
	result := collection.next.FindOne(ctx, filter, opts...)
	logMongoQuery(ctx, collection.name, "FindOne", start, result.Err())
	return result
}

func (collection *loggingCollection) Aggregate(ctx context.Context, pipeline any, opts ...options.Lister[options.AggregateOptions]) (*mongo.Cursor, error) {
	start := time.Now()
	cursor, err := collection.next.Aggregate(ctx, pipeline, opts...)
	logMongoQuery(ctx, collection.name, "Aggregate", start, err)
	return cursor, err
}

func (collection *loggingCollection) Distinct(ctx context.Context, fieldName string, filter any, opts ...options.Lister[options.DistinctOptions]) *mongo.DistinctResult {
	start := time.Now()
	result := collection.next.Distinct(ctx, fieldName, filter, opts...)
	logMongoQuery(ctx, collection.name, "Distinct", start, result.Err())
	return result
}

func logMongoQuery(ctx context.Context, collectionName string, operation string, start time.Time, err error) {
	slog.DebugContext(ctx, "mongo query",
		"collection", collectionName,
		"operation", operation,
		"durationMs", time.Since(start).Milliseconds(),
		"error", err,
	)
}
