package dto

import (
	"net/url"
	"strconv"
	"strings"
)

// Operator 操作符类型
type Operator string

const (
	Equal          Operator = "eq"  // 等于
	NotEqual       Operator = "ne"  // 不等于
	Greater        Operator = "gt"  // 大于
	Less           Operator = "lt"  // 小于
	GreaterOrEqual Operator = "gte" // 大于等于
	LessOrEqual    Operator = "lte" // 小于等于
	In             Operator = "in"  // IN操作符
	Or             Operator = "or"  // OR操作符
)

// OrderDirection 排序方向
type OrderDirection string

const (
	Ascending  OrderDirection = "asc"  // 升序
	Descending OrderDirection = "desc" // 降序
)

// Condition 查询条件
type Condition struct {
	Key   string   `json:"key"`   // 字段名
	Value any      `json:"value"` // 比较值
	Op    Operator `json:"op"`    // 操作符
}

// Order 排序条件
type Order struct {
	Key       string         `json:"key"`       // 排序字段
	Direction OrderDirection `json:"direction"` // 排序方向
}

// QueryRequest 查询请求
type QueryRequest struct {
	Page      int         `json:"page"`      // 页码
	PageSize  int         `json:"page_size"` // 每页数量
	Condition []Condition `json:"condition"` // 查询条件
	Order     []Order     `json:"order"`     // 排序条件
}

// NewQueryRequestFromURL 从URL查询参数创建QueryRequest
func NewQueryRequestFromURL(rawQuery string) *QueryRequest {
	// 解析查询字符串但保留顺序
	query, _ := url.ParseQuery(rawQuery)

	qr := &QueryRequest{
		Page:      1,
		PageSize:  10,
		Condition: make([]Condition, 0),
		Order:     make([]Order, 0),
	}

	// 解析页码
	parsePage(query, qr)

	// 解析每页数量
	parsePageSize(query, qr)

	// 解析查询条件
	parseConditions(rawQuery, qr)

	// 解析排序条件
	parseSort(query, qr)

	return qr
}

// parsePage 解析页码
func parsePage(query url.Values, qr *QueryRequest) {
	if pageStr := query.Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			qr.Page = page
		}
	}
}

// parsePageSize 解析每页数量
func parsePageSize(query url.Values, qr *QueryRequest) {
	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			qr.PageSize = pageSize
		}
	}
}

// parseConditions 解析查询条件
func parseConditions(rawQuery string, qr *QueryRequest) {
	// 手动解析查询字符串以维持参数顺序
	// 分割查询字符串为键值对
	pairs := strings.Split(rawQuery, "&")
	for _, pair := range pairs {
		// 跳过空对
		if pair == "" {
			continue
		}

		// 分割键和值
		kv := strings.SplitN(pair, "=", 2)
		key := kv[0]

		// 跳过特殊参数
		if key == "page" || key == "page_size" || key == "sort" {
			continue
		}

		// 检查是否为条件查询字段
		if condition := parseCondition(key); condition != nil {
			// 如果有值，则解码并赋值
			if len(kv) == 2 {
				value, err := url.QueryUnescape(kv[1])
				if err != nil {
					// 解码失败则使用原值
					value = kv[1]
				}
				condition.Value = value

				// 对于IN操作符，需要特殊处理值为数组的情况
				if condition.Op == In {
					// 处理逗号分隔的值作为数组
					if str, ok := condition.Value.(string); ok && strings.Contains(str, ",") {
						condition.Value = strings.Split(str, ",")
					} else {
						// 单个值也转换为数组
						condition.Value = []string{condition.Value.(string)}
					}
				}
			}
			qr.Condition = append(qr.Condition, *condition)
		}
	}
}

// parseSort 解析排序条件
func parseSort(query url.Values, qr *QueryRequest) {
	// 解析排序条件
	// 支持格式: sort=key1.asc,key2.desc
	if sortStr := query.Get("sort"); sortStr != "" {
		qr.Order = parseOrder(sortStr)
	}
}

// parseCondition 解析条件字段
func parseCondition(key string) *Condition {
	// 查找操作符
	operators := []Operator{GreaterOrEqual, LessOrEqual, Equal, NotEqual, Greater, Less, In, Or}

	for _, op := range operators {
		suffix := "_" + string(op)
		if len(key) > len(suffix) && key[len(key)-len(suffix):] == suffix {
			fieldName := key[:len(key)-len(suffix)]
			return &Condition{
				Key: fieldName,
				Op:  op,
			}
		}
	}

	// 默认为相等条件
	return &Condition{
		Key: key,
		Op:  Equal,
	}
}

// parseOrder 解析排序字段
func parseOrder(sortStr string) []Order {
	orders := make([]Order, 0)

	// 支持格式: key1.asc,key2.desc
	pairs := strings.Split(sortStr, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ".")
		if len(parts) == 2 {
			order := Order{
				Key: parts[0],
			}

			switch parts[1] {
			case "asc":
				order.Direction = Ascending
			case "desc":
				order.Direction = Descending
			default:
				continue // 无效的排序方向
			}

			orders = append(orders, order)
		}
	}

	return orders
}
