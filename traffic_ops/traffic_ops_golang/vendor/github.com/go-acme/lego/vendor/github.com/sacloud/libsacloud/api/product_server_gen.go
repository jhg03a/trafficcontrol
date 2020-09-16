package api

/************************************************
  generated by IDE. for [ProductServerAPI]
************************************************/

import (
	"github.com/sacloud/libsacloud/sacloud"
)

/************************************************
   To support fluent interface for Find()
************************************************/

// Reset 検索条件のリセット
func (api *ProductServerAPI) Reset() *ProductServerAPI {
	api.reset()
	return api
}

// Offset オフセット
func (api *ProductServerAPI) Offset(offset int) *ProductServerAPI {
	api.offset(offset)
	return api
}

// Limit リミット
func (api *ProductServerAPI) Limit(limit int) *ProductServerAPI {
	api.limit(limit)
	return api
}

// Include 取得する項目
func (api *ProductServerAPI) Include(key string) *ProductServerAPI {
	api.include(key)
	return api
}

// Exclude 除外する項目
func (api *ProductServerAPI) Exclude(key string) *ProductServerAPI {
	api.exclude(key)
	return api
}

// FilterBy 指定キーでのフィルター
func (api *ProductServerAPI) FilterBy(key string, value interface{}) *ProductServerAPI {
	api.filterBy(key, value, false)
	return api
}

// FilterMultiBy 任意項目でのフィルタ(完全一致 OR条件)
func (api *ProductServerAPI) FilterMultiBy(key string, value interface{}) *ProductServerAPI {
	api.filterBy(key, value, true)
	return api
}

// WithNameLike 名称条件
func (api *ProductServerAPI) WithNameLike(name string) *ProductServerAPI {
	return api.FilterBy("Name", name)
}

// WithTag タグ条件
func (api *ProductServerAPI) WithTag(tag string) *ProductServerAPI {
	return api.FilterBy("Tags.Name", tag)
}

// WithTags タグ(複数)条件
func (api *ProductServerAPI) WithTags(tags []string) *ProductServerAPI {
	return api.FilterBy("Tags.Name", []interface{}{tags})
}

// func (api *ProductServerAPI) WithSizeGib(size int) *ProductServerAPI {
// 	api.FilterBy("SizeMB", size*1024)
// 	return api
// }

// func (api *ProductServerAPI) WithSharedScope() *ProductServerAPI {
// 	api.FilterBy("Scope", "shared")
// 	return api
// }

// func (api *ProductServerAPI) WithUserScope() *ProductServerAPI {
// 	api.FilterBy("Scope", "user")
// 	return api
// }

// SortBy 指定キーでのソート
func (api *ProductServerAPI) SortBy(key string, reverse bool) *ProductServerAPI {
	api.sortBy(key, reverse)
	return api
}

// SortByName 名称でのソート
func (api *ProductServerAPI) SortByName(reverse bool) *ProductServerAPI {
	api.sortByName(reverse)
	return api
}

// func (api *ProductServerAPI) SortBySize(reverse bool) *ProductServerAPI {
// 	api.sortBy("SizeMB", reverse)
// 	return api
// }

/************************************************
   To support Setxxx interface for Find()
************************************************/

// SetEmpty 検索条件のリセット
func (api *ProductServerAPI) SetEmpty() {
	api.reset()
}

// SetOffset オフセット
func (api *ProductServerAPI) SetOffset(offset int) {
	api.offset(offset)
}

// SetLimit リミット
func (api *ProductServerAPI) SetLimit(limit int) {
	api.limit(limit)
}

// SetInclude 取得する項目
func (api *ProductServerAPI) SetInclude(key string) {
	api.include(key)
}

// SetExclude 除外する項目
func (api *ProductServerAPI) SetExclude(key string) {
	api.exclude(key)
}

// SetFilterBy 指定キーでのフィルター
func (api *ProductServerAPI) SetFilterBy(key string, value interface{}) {
	api.filterBy(key, value, false)
}

// SetFilterMultiBy 任意項目でのフィルタ(完全一致 OR条件)
func (api *ProductServerAPI) SetFilterMultiBy(key string, value interface{}) {
	api.filterBy(key, value, true)
}

// SetNameLike 名称条件
func (api *ProductServerAPI) SetNameLike(name string) {
	api.FilterBy("Name", name)
}

// SetTag タグ条件
func (api *ProductServerAPI) SetTag(tag string) {
	api.FilterBy("Tags.Name", tag)
}

// SetTags タグ(複数)条件
func (api *ProductServerAPI) SetTags(tags []string) {
	api.FilterBy("Tags.Name", []interface{}{tags})
}

// func (api *ProductServerAPI) SetSizeGib(size int) {
// 	api.FilterBy("SizeMB", size*1024)
// }

// func (api *ProductServerAPI) SetSharedScope() {
// 	api.FilterBy("Scope", "shared")
// }

// func (api *ProductServerAPI) SetUserScope() {
// 	api.FilterBy("Scope", "user")
// }

// SetSortBy 指定キーでのソート
func (api *ProductServerAPI) SetSortBy(key string, reverse bool) {
	api.sortBy(key, reverse)
}

// SetSortByName 名称でのソート
func (api *ProductServerAPI) SetSortByName(reverse bool) {
	api.sortByName(reverse)
}

// func (api *ProductServerAPI) SetSortBySize(reverse bool) {
// 	api.sortBy("SizeMB", reverse)
// }

/************************************************
  To support CRUD(Create/Read/Update/Delete)
************************************************/

//func (api *ProductServerAPI) New() *sacloud.ProductServer {
//	return &sacloud.ProductServer{}
//}

// func (api *ProductServerAPI) Create(value *sacloud.ProductServer) (*sacloud.ProductServer, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.create(api.createRequest(value), res)
// 	})
// }

// Read 読み取り
func (api *ProductServerAPI) Read(id int64) (*sacloud.ProductServer, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.read(id, nil, res)
	})
}

// func (api *ProductServerAPI) Update(id int64, value *sacloud.ProductServer) (*sacloud.ProductServer, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.update(id, api.createRequest(value), res)
// 	})
// }

// func (api *ProductServerAPI) Delete(id int64) (*sacloud.ProductServer, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.delete(id, nil, res)
// 	})
// }

/************************************************
  Inner functions
************************************************/

func (api *ProductServerAPI) setStateValue(setFunc func(*sacloud.Request)) *ProductServerAPI {
	api.baseAPI.setStateValue(setFunc)
	return api
}

func (api *ProductServerAPI) request(f func(*sacloud.Response) error) (*sacloud.ProductServer, error) {
	res := &sacloud.Response{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.ServerPlan, nil
}

func (api *ProductServerAPI) createRequest(value *sacloud.ProductServer) *sacloud.Request {
	req := &sacloud.Request{}
	req.ServerPlan = value
	return req
}
