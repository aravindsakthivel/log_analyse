package controller

type SDBCtrl struct{}

func (c SDBCtrl) Users() (*SUsersCL, error) {
	cl := &SUsersCL{}
	err := cl.setCollection()
	if err != nil {
		return nil, err
	}
	return cl, nil
}

func (f SDBCtrl) Files() (*SFilesCL, error) {
	cl := &SFilesCL{}
	err := cl.setCollection()
	if err != nil {
		return nil, err
	}
	return cl, nil
}
