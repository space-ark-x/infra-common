package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestQueryBuilder 测试MongoDB查询构建器功能
func TestQueryBuilder(t *testing.T) {
	// 连接到MongoDB数据库
	// 数据库地址: 127.0.0.1 用户名: admin 密码: 123456
	uri := "mongodb://admin:123456@127.0.0.1:27017/admin"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("无法连接到MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			t.Errorf("断开MongoDB连接失败: %v", err)
		}
	}()

	// 检查连接是否有效
	err = client.Ping(ctx, nil)
	if err != nil {
		t.Skipf("MongoDB服务器不可达: %v", err)
	}

	// 创建测试数据库和集合
	db := client.Database("test_db")
	collection := db.Collection("test_collection")
	defer func() {
		// 清理测试数据
		_ = collection.Drop(ctx)
	}()

	// 插入测试数据
	testDocs := []interface{}{
		bson.M{"name": "Alice", "age": 25, "city": "New York"},
		bson.M{"name": "Bob", "age": 30, "city": "San Francisco"},
		bson.M{"name": "Charlie", "age": 35, "city": "New York"},
		bson.M{"name": "David", "age": 20, "city": "Los Angeles"},
		bson.M{"name": "Eve", "age": 28, "city": "Chicago"},
	}
	_, err = collection.InsertMany(ctx, testDocs)
	if err != nil {
		t.Fatalf("插入测试数据失败: %v", err)
	}

	t.Run("Eq查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Eq("name", "Alice").Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Alice", results[0]["name"])
	})

	t.Run("Ne查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Ne("name", "Alice").Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 4)
	})

	t.Run("Gt查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Gt("age", 30).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 1) // 只有Charlie(35)
	})

	t.Run("Gte查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Gte("age", 30).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 2) // Bob(30)和Charlie(35)
	})

	t.Run("Lt查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Lt("age", 25).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 1) // 只有David(20)
	})

	t.Run("Lte查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Lte("age", 25).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 2) // Alice(25)和David(20)
	})

	t.Run("In查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.In("name", []interface{}{"Alice", "Bob"}).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("Nin查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Nin("name", []interface{}{"Alice", "Bob"}).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("Exists查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Exists("name", true).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 5)
	})

	t.Run("组合查询测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Eq("city", "New York").Gte("age", 30).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		if len(results) > 0 {
			assert.Equal(t, "Charlie", results[0]["name"])
		}
	})

	t.Run("Limit和Skip测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Limit(3).Skip(1).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("Sort测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Sort(bson.D{{Key: "age", Value: 1}}).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 5)
		if len(results) >= 2 {
			// 第一个是年龄最小的应该David(20)
			assert.EqualValues(t, 20, results[0]["age"])
			// 第二个应该是Alice(25)
			assert.EqualValues(t, 25, results[1]["age"])
		}
	})

	t.Run("Project测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var results []bson.M
		err := queryBuilder.Project(bson.D{{Key: "name", Value: 1}, {Key: "_id", Value: 0}}).Find("test_collection", &results)
		assert.NoError(t, err)
		assert.Len(t, results, 5)
		for _, result := range results {
			assert.Contains(t, result, "name")
			assert.NotContains(t, result, "age")
			assert.NotContains(t, result, "city")
		}
	})

	t.Run("Count测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		count, err := queryBuilder.Gt("age", 25).Count("test_collection")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count) // Bob(30), Charlie(35), Eve(28)
	})

	t.Run("FindOne测试", func(t *testing.T) {
		queryBuilder := NewQueryBuilder(client, "test_db")
		var result bson.M
		err := queryBuilder.Eq("name", "Alice").FindOne("test_collection", &result)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", result["name"])
	})
}