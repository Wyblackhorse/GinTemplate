/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package rule

type LoginRequest struct {
	//Username string `json:"username" binding:"required"`
	Username   string `form:"username"  binding:"required,min=5,max=10"` //账户名字
	Password   string `form:"password"  binding:"required,min=5,max=10"` //密码
	GoogleCode string `form:"code"  binding:"omitempty,min=6,max=6"`     //谷歌验证码
}
