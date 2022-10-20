package simulated

import (
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/ipfs/go-cid"
)

func (s *FvmSimulator) SetAccount(actorID abi.ActorID, addr address.Address, actor migration.Actor) {
	s.actorLk.Lock()
	defer s.actorLk.Unlock()

	s.actorsMap[actorID] = actor
	s.addressMap[addr] = actorID
}

func (s *FvmSimulator) ResolveAddress(addr address.Address) (abi.ActorID, error) {
	s.actorLk.Lock()
	defer s.actorLk.Unlock()
	id, ok := s.addressMap[addr]
	if !ok {
		return 0, ErrorNotFound
	}
	return id, nil
}

func (s *FvmSimulator) NewActorAddress() (address.Address, error) {
	seed := time.Now().String()
	return address.NewActorAddress([]byte(seed))
}

func (s *FvmSimulator) CreateActor(actorID abi.ActorID, codeCid cid.Cid) error {
	s.SetAccount(actorID, address.Address{}, migration.Actor{Code: codeCid})
	return nil
}

func (s *FvmSimulator) GetActorCodeCid(addr address.Address) (*cid.Cid, error) {
	acstat, err := s.getActorWithAddress(addr)
	if err != nil {
		return nil, err
	}
	return &acstat.Code, nil
}
