package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/polundrra/shortlink/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LinkServiceSuit struct {
	suite.Suite
	mockCtrl *gomock.Controller
	repoMock *repo.MockLinkRepo
	SUT Service
}

func TestService(t *testing.T) {
	suite.Run(t, new(LinkServiceSuit))
}

func (s *LinkServiceSuit) SetupTest() {
	ctrl := gomock.NewController(s.T())

	s.mockCtrl = ctrl
	s.repoMock = repo.NewMockLinkRepo(ctrl)

	s.SUT = New(repo.Opts{
		Timeout: 1,
	}, s.repoMock)
}

func (s *LinkServiceSuit) TestCreateShortLink_CustomEnd() {
	defer s.mockCtrl.Finish()
	a := assert.New(s.T())

	ctx := context.Background()
	url := "foo12345"
	customEnd := "foo"

	s.repoMock.EXPECT().IsCodeExists(ctx, customEnd).Times(1).Return(false, nil)
	s.repoMock.EXPECT().SetLink(ctx, url, customEnd, true).Times(1).Return(nil)

	res, err := s.SUT.CreateShortLink(ctx, url, customEnd)
	a.NoError(err)
	a.Equal(customEnd, res)
}

func (s *LinkServiceSuit) TestCreateShortLink_CodeConflictErr() {
	defer s.mockCtrl.Finish()
	a := assert.New(s.T())

	ctx := context.Background()
	url := "foo12345"
	customEnd := "foo"
	exUrl := "bar12345"

	s.repoMock.EXPECT().IsCodeExists(ctx, customEnd).Times(1).Return(true, nil)
	s.repoMock.EXPECT().GetLongLinkByCode(ctx, customEnd).Times(1).Return(exUrl, nil)

	res, err := s.SUT.CreateShortLink(ctx, url, customEnd)
	a.Error(ErrCodeConflict, err)
	a.Equal("", res)
}

func (s *LinkServiceSuit) TestCreateShortLink_CustomEnd_SameURL() {
	defer s.mockCtrl.Finish()
	a := assert.New(s.T())

	ctx := context.Background()
	url := "foo12345"
	customEnd := "foo"
	exUrl := "foo12345"

	s.repoMock.EXPECT().IsCodeExists(ctx, customEnd).Times(1).Return(true, nil)
	s.repoMock.EXPECT().GetLongLinkByCode(ctx, customEnd).Times(1).Return(exUrl, nil)

	res2, err := s.SUT.CreateShortLink(ctx, url, customEnd)
	a.NoError(err)
	a.Equal(customEnd, res2)
}

func (s *LinkServiceSuit) TestCreateShortLink_CodeExists() {
	defer s.mockCtrl.Finish()
	a := assert.New(s.T())

	ctx := context.Background()
	url := "foo12345"
	exCode := "f"

	s.repoMock.EXPECT().GetCodeByLongLink(ctx, url).Times(1).Return(exCode, nil)

	res, err := s.SUT.CreateShortLink(ctx, url, "")
	a.NoError(err)
	a.Equal(exCode, res)
}

/*func (s *LinkServiceSuit) TestCreateShortLink() {
	defer s.mockCtrl.Finish()

	ctx := context.Background()
	url := "foo12345"
	code :=

	s.repoMock.EXPECT().GetCodeByLongLink(ctx, url).Times(1).Return("", nil)
	s.repoMock.EXPECT().SetLink(ctx, url, code, false).Times(1).Return(nil)


}*/

func (s *LinkServiceSuit) TestCreateShortLink_IsCodeExistsErr() {
	defer s.mockCtrl.Finish()
	a := assert.New(s.T())

	ctx := context.Background()
	url := "foo12345"
	customEnd := "foo"
	expected := errors.New("any")

	s.repoMock.EXPECT().IsCodeExists(ctx, customEnd).Times(1).Return(false, expected)

	res, err := s.SUT.CreateShortLink(ctx, url, customEnd)
	a.Error(expected, err)
	a.Equal("", res)
}

func (s *LinkServiceSuit) TestCreateShortLink_GetLongLinkByCodeErr() {
	defer s.mockCtrl.Finish()
	a := assert.New(s.T())

	ctx := context.Background()
	url := "foo12345"
	customEnd := "foo"
	expected := errors.New("any")

	s.repoMock.EXPECT().IsCodeExists(ctx, customEnd).Times(1).Return(true, nil)
	s.repoMock.EXPECT().GetLongLinkByCode(ctx, customEnd).Times(1).Return("", expected)

	res, err := s.SUT.CreateShortLink(ctx, url, customEnd)
	a.Error(expected, err)
	a.Equal("", res)
}

func (s *LinkServiceSuit) TestCreateShortLink_GetCodeByLongLinkErr() {
	defer s.mockCtrl.Finish()
	a := assert.New(s.T())

	ctx := context.Background()
	url := "foo12345"
	customEnd := ""
	expected := errors.New("any")

	s.repoMock.EXPECT().GetCodeByLongLink(ctx, url).Times(1).Return("", expected)

	res, err := s.SUT.CreateShortLink(ctx, url, customEnd)
	a.Error(expected, err)
	a.Equal("", res)
}

func TestToBase62(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		n uint64
		expected string
	}{
		{0,  "0"},
		{1, "1"},
		{10, "a"},
		{36, "A"},
		{100, "1C"},
		{1234567891, "1ly7vl"},
	}

	for _, tc := range testCases {
		res := toBase62(tc.n)
		a.Equal(tc.expected, res)
	}
}

