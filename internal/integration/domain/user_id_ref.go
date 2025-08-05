package domain

type UserIDRef string

func (u UserIDRef) String() string {
	return string(u)
}
