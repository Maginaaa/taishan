package logic

import (
	"github.com/gin-gonic/gin"
	"scene/internal/biz/log"
	"scene/internal/dal"
	"scene/rao"
)

func GetTagLiat(ctx *gin.Context, tp int32) (res []rao.Tag, err error) {
	tx := dal.GetQuery().Tag
	res = make([]rao.Tag, 0)
	list, err := tx.WithContext(ctx).Where(tx.Type.Eq(tp), tx.IsDelete.Not()).Find()
	if err != nil {
		log.Logger.Error("logic.tag.GetTagLiat.Find，err:", err)
		return
	}
	for _, l := range list {
		res = append(res, rao.Tag{
			ID:    l.ID,
			Label: l.Label,
		})
	}
	return
}
