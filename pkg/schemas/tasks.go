package schemas

type Task struct {
	Id          string
	Title       string
	Description string
	Status      string
}

type TaskSort struct {
	Id            string
	Status        string
	Sorting_order []string
}
