package srv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/Characters/pkg/model"
	"github.com/ShatteredRealms/Characters/pkg/pb"
	"github.com/ShatteredRealms/Characters/pkg/service"
	"github.com/ShatteredRealms/Characters/pkg/srv"
)

var _ = Describe("Character", func() {
	var (
		serviceSrv pb.CharactersServiceServer
	)

	BeforeEach(func() {
		serviceSrv = srv.NewCharacterServiceServer(newServiceMock(), nil)
		Expect(serviceSrv).NotTo(BeNil())
	})
})

type serviceMock struct {
}

// AddPlayTime implements service.CharacterService
func (serviceMock) AddPlayTime(characterId uint64, amount uint64) (uint64, error) {
	return 1, nil
}

// Create implements service.CharacterService
func (serviceMock) Create(ownerId uint64, name string, genderId uint64, realmId uint64) (*model.Character, error) {
	return nil, nil
}

// Delete implements service.CharacterService
func (serviceMock) Delete(id uint64) error {
	return nil
}

// Edit implements service.CharacterService
func (serviceMock) Edit(character *pb.Character) (*model.Character, error) {
	return nil, nil
}

// FindAll implements service.CharacterService
func (serviceMock) FindAll() ([]*model.Character, error) {
	return nil, nil
}

// FindAllByOwner implements service.CharacterService
func (serviceMock) FindAllByOwner(ownerId uint64) ([]*model.Character, error) {
	return nil, nil
}

// FindById implements service.CharacterService
func (serviceMock) FindById(id uint64) (*model.Character, error) {
	return nil, nil
}

func newServiceMock() service.CharacterService {
	return serviceMock{}
}
