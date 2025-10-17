package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// QueryBuilder MongoDB查询构建器
type QueryBuilder struct {
	client     *mongo.Client
	database   string
	filter     bson.D
	projection bson.D
	sort       bson.D
	limit      *int64
	skip       *int64
}

// NewQueryBuilder 创建一个新的查询构建器
func NewQueryBuilder(client *mongo.Client, database string) *QueryBuilder {
	return &QueryBuilder{
		client:   client,
		database: database,
		filter:   bson.D{},
	}
}

// Eq 添加等于条件
func (q *QueryBuilder) Eq(field string, value interface{}) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: value})
	return q
}

// Ne 添加不等于条件
func (q *QueryBuilder) Ne(field string, value interface{}) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$ne", Value: value}}})
	return q
}

// Gt 添加大于条件
func (q *QueryBuilder) Gt(field string, value interface{}) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$gt", Value: value}}})
	return q
}

// Gte 添加大于等于条件
func (q *QueryBuilder) Gte(field string, value interface{}) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$gte", Value: value}}})
	return q
}

// Lt 添加小于条件
func (q *QueryBuilder) Lt(field string, value interface{}) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$lt", Value: value}}})
	return q
}

// Lte 添加小于等于条件
func (q *QueryBuilder) Lte(field string, value interface{}) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$lte", Value: value}}})
	return q
}

// In 添加IN条件
func (q *QueryBuilder) In(field string, values []interface{}) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$in", Value: values}}})
	return q
}

// Nin 添加NOT IN条件
func (q *QueryBuilder) Nin(field string, values []interface{}) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$nin", Value: values}}})
	return q
}

// Exists 添加存在性检查条件
func (q *QueryBuilder) Exists(field string, exists bool) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$exists", Value: exists}}})
	return q
}

// Regex 添加正则表达式匹配条件
func (q *QueryBuilder) Regex(field string, pattern string) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: field, Value: bson.D{{Key: "$regex", Value: pattern}}})
	return q
}

// And 添加AND条件
func (q *QueryBuilder) And(conditions bson.D) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: "$and", Value: conditions})
	return q
}

// Or 添加OR条件
func (q *QueryBuilder) Or(conditions bson.D) *QueryBuilder {
	q.filter = append(q.filter, bson.E{Key: "$or", Value: conditions})
	return q
}

// Project 设置投影字段
func (q *QueryBuilder) Project(projection bson.D) *QueryBuilder {
	q.projection = projection
	return q
}

// Sort 设置排序字段
func (q *QueryBuilder) Sort(sort bson.D) *QueryBuilder {
	q.sort = sort
	return q
}

// Limit 设置限制数量
func (q *QueryBuilder) Limit(limit int64) *QueryBuilder {
	q.limit = &limit
	return q
}

// Skip 设置跳过数量
func (q *QueryBuilder) Skip(skip int64) *QueryBuilder {
	q.skip = &skip
	return q
}

// Find 执行查询并返回多个结果
func (q *QueryBuilder) Find(collection string, results interface{}) error {
	ctx := context.Background()
	coll := q.client.Database(q.database).Collection(collection)

	opts := options.Find()
	if q.projection != nil {
		opts.SetProjection(q.projection)
	}
	if q.sort != nil {
		opts.SetSort(q.sort)
	}
	if q.limit != nil {
		opts.SetLimit(*q.limit)
	}
	if q.skip != nil {
		opts.SetSkip(*q.skip)
	}

	cursor, err := coll.Find(ctx, q.filter, opts)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// FindOne 执行查询并返回单个结果
func (q *QueryBuilder) FindOne(collection string, result interface{}) error {
	ctx := context.Background()
	coll := q.client.Database(q.database).Collection(collection)

	opts := &options.FindOneOptions{}
	if q.projection != nil {
		opts.SetProjection(q.projection)
	}
	if q.sort != nil {
		opts.SetSort(q.sort)
	}

	return coll.FindOne(ctx, q.filter, opts).Decode(result)
}

// Count 执行查询并返回匹配的文档数量
func (q *QueryBuilder) Count(collection string) (int64, error) {
	ctx := context.Background()
	coll := q.client.Database(q.database).Collection(collection)

	opts := options.Count()
	if q.limit != nil {
		opts.SetLimit(*q.limit)
	}
	if q.skip != nil {
		opts.SetSkip(*q.skip)
	}

	return coll.CountDocuments(ctx, q.filter, opts)
}
