package errmsg

const (
	BindFailed = iota + 10001
	ParamWrong
	EncryptedPwdFailed
	RegisterFailed
	RegisterSuccess
	PhoneIsEmpty
	PhoneAlreadyExists
	PhoneDoesNotExists
	PhoneCodeIsEmpty
	NameIsEmpty
	PwdIsEmpty
	AgeIllegal
	EmailIsEmpty
	EmailAlreadyExists
	EmailFormatWrong
	SetCodeRedisFailed
	GetCodeRedisFailed
	NameAlreadyExists
	SendCodeSuccess
	SendSmsFailed
	SendEmailFailed
	PhoneCodeWrong
	PhoneCodeExpired
	MySQLQueryFailed
	MySQLCreateAUserFailed
	RecordNotFound
	DuplicateEntry
	CreateImgCodeFailed
)
const (
	LoginFailed = iota + 11001
	LoginSuccess
	PwdWrong
	MySQLUpdateFailed
)
const (
	BindFailedMsg             = "参数绑定失败"
	ParamWrongMsg             = "部分参数为空"
	EncryptedPwdFailedMsg     = "密码加密失败"
	RegisterFailedMsg         = "注册失败"
	RegisterSuccessMsg        = "注册成功"
	PhoneIsEmptyMsg           = "手机号为空"
	PhoneAlreadyExistsMsg     = "手机号已经存在"
	PhoneDoesNotExistsMsg     = "手机号不存在"
	PhoneCodeIsEmptyMsg       = "短信验证码为空"
	NameIsEmptyMsg            = "姓名为空"
	PwdIsEmptyMsg             = "密码为空"
	AgeIllegalMsg             = "年龄不合法"
	EmailIsEmptyMsg           = "邮箱为空"
	EmailAlreadyExistsMsg     = "邮箱已经存在"
	EmailFormatWrongMsg       = "邮箱格式错误"
	SetCodeRedisFailedMsg     = "设置验证码缓存失败"
	GetCodeRedisFailedMsg     = "获取验证码缓存失败"
	SendCodeSuccessMsg        = "发送验证码成功"
	SendSmsFailedMsg          = "发送短信验证码失败"
	SendEmailFailedMsg        = "发送邮箱验证码失败"
	PhoneCodeWrongMsg         = "短信验证码错误"
	PhoneCodeExpiredMsg       = "短信验证码过期"
	MySQLQueryFailedMsg       = "数据库查询出错"
	MySQLCreateAUserFailedMsg = "数据库创建用户失败"
	RecordNotFoundMsg         = "记录没有找到"
	DuplicateEntryMsg         = "重复记录"
	CreateImgCodeFailedMsg    = "生成图片验证码失败"
	NameAlreadyExistsMsg      = "用户名已存在"

	LoginFailedMsg       = "登录失败"
	LoginSuccessMsg      = "登录成功"
	PwdWrongMsg          = "密码错误"
	MySQLUpdateFailedMsg = "数据库更新失败"
)
