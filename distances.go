package editdistance

var DefaultLevenshtein = EditOperations{
	Match{},
	Insertion{1},
	Deletion{1},
	Substitution{1},
}

var DefaultDamerau = EditOperations{
	Match{},
	Insertion{1},
	Deletion{1},
	Substitution{1},
	Transposition{1},
}

var DefaultLCS = EditOperations{
	Match{},
	Insertion{1},
	Deletion{1},
}
