package model

// ARLevel 定义了 Level 的3种级别。
type ARLevel interface {
	ERROR() ARLevel
	WARN() ARLevel
	INFO() ARLevel

	//private()
}
