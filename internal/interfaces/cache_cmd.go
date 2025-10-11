package interfaces

type CacheStatusCmd interface {
	Result() (string, error)
}

type CacheStringCmd interface {
	Result() (string, error)
}

type CacheIntCmd interface {
	Result() (int64, error)
}
