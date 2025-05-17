package rao

// SceneCase TODO: Extend后期需要抽象成接口
type SceneCase struct {
	CaseID   int32       `json:"case_id"`             // case id
	Title    string      `json:"title"`               // case名
	ParentID int32       `json:"parent_id,omitempty"` // 父节点id
	Type     int32       `json:"type"`                // 最后修改人
	SceneID  int32       `json:"scene_id"`            // 最后修改时间
	Sort     int32       `json:"sort"`                // 排序
	Disabled bool        `json:"disabled"`            // 是否禁用
	Extend   interface{} `json:"extend"`
}

type ImportCase struct {
	CaseID  int32 `json:"case_id"`
	SceneID int32 `json:"scene_id"`
}

type CaseSortReq struct {
	Before   SceneCase `json:"before"`
	After    SceneCase `json:"after"`
	Position int32     `json:"position"`
}

type SceneCaseTree struct {
	*SceneCase
	Children []*SceneCaseTree `json:"children,omitempty"`
}

type CurlParseReq struct {
	Curl string `json:"curl" binding:"required"`
}

type CurlParseResp struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Params  map[string]string `json:"params"`
}

type SQLExecutionRequest struct {
	ID  int32  `json:"db_id"`
	SQL string `json:"sql"`
}

type SQLExecutionResponse struct {
	IsSuccess bool                     `json:"is_success"` // 操作是否成功
	SqlResult []map[string]interface{} `json:"sql_result"` // 数据库查询记录
}

const (
	Before = iota
	After
	Inner
)

const (
	HttpCaseType int32 = iota + 1
	_
	_
	_
	_
	_
	_
	_
	_
	_
	LogicControlType
)

func ConvertTreeSliceToList(trees []*SceneCaseTree) []*SceneCaseTree {
	var result []*SceneCaseTree
	for _, tree := range trees {
		traverseTree(tree, &result)
	}
	return result
}

func traverseTree(node *SceneCaseTree, result *[]*SceneCaseTree) {
	if node == nil {
		return
	}
	*result = append(*result, node)
	for _, child := range node.Children {
		traverseTree(child, result)
	}
}
