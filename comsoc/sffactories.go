package comsoc

func SWFFactory(swf func (p Profile) (Count, error), tiebreak func ([]Alternative) (Alternative, error)) (func(Profile) ([]Alternative, error)) {

}

func SCFFactory(swf func (p Profile) ([]Alternative, error), tiebreak func ([]Alternative) (Alternative, error)) (func(Profile) (Alternative, error)) {
	
}
