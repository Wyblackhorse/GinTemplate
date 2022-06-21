/**
 * @Author $
 * @Description //TODO $  自动采集 YouTube
 * @Date $ $
 * @Param $
 * @return $
 **/
package collect

import (
	"github.com/jinzhu/gorm"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"time"
)

func Collection(db *gorm.DB) {
	for true {
		cc := make([]model.Collection, 0)
		db.Where("status=?", 1).Find(&cc)

		if len(cc) > 0 {
			for _, v := range cc {
				if v.Kinds == 1 {
					//YouTube
					urlArray := tools.FormYouTuBeUrl(v.TaskUrl)
					if len(urlArray) > 0 {
						for _, i2 := range urlArray {
							add := model.Task{
								ApplyId:    v.Kinds,
								Remark:     v.Remark,
								TaskType:   v.TaskType,
								TaskUrl:    i2,
								EndTime:    time.Now().Unix() + int64(v.Expiry*3600),
								Price:      v.Price,
								TaskLevel:  v.TaskLevel,
								TaskNum:    v.TaskNum,
								Created:    time.Now().Unix(),
								DemandSide: v.DemandSide,
								Status:     1,
							}
							db.Save(&add)
						}
					}
					db.Model(&model.Collection{}).Where("id=?", v.ID).Update(&model.Collection{Status: 2})
				}
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
