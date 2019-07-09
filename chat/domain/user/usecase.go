package user

// Function signatures for user usecases
type Fetch func(userName string) (User, error)
type FetchAll func() ([]User, error)
type Kick func(userName string) (User, error)
type GetBool func(userName string) bool
