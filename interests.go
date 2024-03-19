package main

// An interest differs from a skill by being an subject or hobby. It does not represent mastery of a subject, merely thought on a subject. It includes all skills, as

// An interest is a noun or noun+verb phrase that has a wikipedia page
// this includes both careers and skills

type Interest struct {
	id             string
	name           string
	secondaryNames []string
	url            string
	InterestType   InterestType
	people         []string
}

type Term struct {
	nouns string
	verb  string
}

type InterestType struct {

	// A Career is a job someone can make a living on: eg irs occupation data
	Career bool

	// A skill is mastery of a subject, often it overlaps with career. One who knows how to do electrical work might be an electrician
	Skill bool

	// A hobby is something primarily done for enjoyment
	// it usually doesn't make a ton of money
	// and it is a lot more flexable than skills or careers: eg, old cars, weightlifting
	Hobby bool
}

// Find takes in the post text. This excludes all other forms of media other than text.
// it finds the nouns and builds an interest database
// func (i *Interest) Find(post string) []Interest
