package fin

import "github.com/clstb/phi/pkg/pb"

type Accounts struct {
	Data   []Account
	byId   map[string]int32
	byName map[string]int32
}

func AccountsFromPB(pb *pb.Accounts) Accounts {
	var data []Account
	byId := make(map[string]int32)
	byName := make(map[string]int32)
	var i int32
	for _, v := range pb.Data {
		account := AccountFromPB(v)
		data = append(data, account)
		byId[account.Id] = i
		byName[account.Name] = i
		i++
	}

	return Accounts{
		Data:   data,
		byId:   byId,
		byName: byName,
	}
}

func (a Accounts) PB() *pb.Accounts {
	var data []*pb.Account
	byId := make(map[string]int32)
	byName := make(map[string]int32)
	var i int32
	for _, account := range a.Data {
		pb := account.PB()
		data = append(data, pb)
		byId[account.Id] = i
		byName[account.Name] = i
		i++
	}

	return &pb.Accounts{
		Data:   data,
		ById:   byId,
		ByName: byName,
	}
}

func (a Accounts) ById(id string) (Account, bool) {
	i, ok := a.byId[id]
	if !ok {
		return Account{}, false
	}
	return a.Data[i], true
}

func (a Accounts) ByName(name string) (Account, bool) {
	i, ok := a.byName[name]
	if !ok {
		return Account{}, false
	}
	return a.Data[i], true
}
