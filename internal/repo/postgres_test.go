package repo

// Integration test
// Be aware that SetupSuite runs pg docker container
// It has 30 seconds lifetime and is being killed by TearDownSuite
// TearDownSuite is called deferred by testify and any call to Fatal/Fatalf will cancel docker purge call,
// so don't call fatal or any other os.Exit in tests
import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"strconv"
	"testing"
)

type PostgresSuite struct {
	suite.Suite
	pool *dockertest.Pool
	postgres *dockertest.Resource
	db *sql.DB
	SUT *postgres
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}

func (s *PostgresSuite) SetupSuite() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		s.T().Logf("error creating testing env: %v", err)
	}

	resource, err := pool.Run("postgres", "12.3", []string{"POSTGRES_PASSWORD=secret","POSTGRES_USER=test_user","POSTGRES_DB=short"})
	if err != nil {
		s.T().Logf("couldn't start test postgres container: %v", err)
	}

	var db *sql.DB
	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://test_user:secret@localhost:%s/%s?sslmode=disable", resource.GetPort("5432/tcp"), "short"))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		s.T().Logf("couldn't connect to postgres: %v", err)
	}
	s.db = db

	init, err := ioutil.ReadFile("../../scripts/sql/init.sql")
	if err != nil {
		s.T().Logf("couldn't open init ddl: %v", err)
	}

	if _, err := db.Exec(string(init)); err != nil {
		s.T().Logf("couldn't init postgres schema: %v", err)
	}

	if err := resource.Expire(30); err != nil {
		s.T().Logf("couldn't init postgres schema: %v", err)
	}

	s.pool = pool
	s.postgres = resource

	port, err := strconv.Atoi(resource.GetPort("5432/tcp"))
	sut, err := New(Opts{
		Host:     "localhost",
		Port:     uint16(port),
		Database: "short",
		User:     "test_user",
		Password: "secret",
		Timeout:  1,
	})
	if err != nil {
		s.T().Logf("couldn't init SUT: %v", err)
	}

	s.SUT = sut.(*postgres)
}

func (s *PostgresSuite) TearDownSuite() {
	if err := s.pool.Purge(s.postgres); err != nil {
		s.T().Logf("couldn't kill test postgres container: %v", err)
	}
}

type row struct {
	Url sql.NullString
	Code sql.NullString
	IsCustom sql.NullBool
}

func (s *PostgresSuite) SetLink() {
	longLink := "foo"
	code := "f"
	isCustom := false
	ctx := context.Background()

	rows, err := s.db.Query("select * from link where code=$1", code)
	s.Assert().NoError(err)
	s.Assert().False(rows.Next())

	err = s.SUT.SetLink(ctx, longLink, code, isCustom)
	s.Assert().NoError(err)

	rows, err = s.db.Query("select url, code, is_custom from link where code=$1", code)
	s.Assert().NoError(err)

	res := row{}
	c := 0
	for rows.Next() {
		c++
		err = rows.Scan(&res.Url, &res.Code, &res.IsCustom)
		s.Assert().NoError(err)
	}
	s.Assert().Equal(1, c)
	res1, err := res.Url.Value()
	s.Assert().NoError(err)
	s.Assert().Equal("foo", res1)

	res2, err := res.Code.Value()
	s.Assert().NoError(err)
	s.Assert().Equal("f", res2)

	res3, err := res.IsCustom.Value()
	s.Assert().NoError(err)
	s.Assert().Equal(false, res3)
}

func (s *PostgresSuite) GetLongLinkByCode() {
	ctx := context.Background()
	code1 := "f"
	url1 := "foo"
	url2 := "bar"
	code2 := "b"

	_, err := s.db.Exec("insert into link(url, code) values ($1, $2), ($3, $4)", url1, code1, url2, code2)
	s.Assert().NoError(err)

	res, err := s.SUT.GetLongLinkByCode(ctx, code1)
	s.Assert().NoError(err)
	s.Assert().Equal("foo", res)
}

func (s *PostgresSuite) GetCodeByLongLink() {
	ctx := context.Background()
	longLink := "foo"

	res, err := s.SUT.GetCodeByLongLink(ctx, longLink)
	s.Assert().NoError(err)
	s.Assert().Equal("", res)

	_, err = s.db.Exec("insert into link(url, code, is_custom) values ($1, $2, $3)", longLink, "1", true)
	s.Assert().NoError(err)

	res1, err := s.SUT.GetCodeByLongLink(ctx, longLink)
	s.Assert().NoError(err)
	s.Assert().Equal("", res1)

	_, err = s.db.Exec("insert into link(url, code, is_custom) values ($1, $2, $3)", longLink, "2", false)
	s.Assert().NoError(err)

	res2, err := s.SUT.GetCodeByLongLink(ctx, longLink)
	s.Assert().NoError(err)
	s.Assert().Equal("2", res2)
}

func (s *PostgresSuite) GetNextSeq() {
	ctx := context.Background()

	res, err := s.SUT.GetNextSeq(ctx)
	s.Assert().NoError(err)
	s.Assert().Equal(1, res)

	for i := 0; i < 5; i++ {
		res, err = s.SUT.GetNextSeq(ctx)
		s.Assert().NoError(err)
	}
	s.Assert().Equal(6, res)
}

func (s *PostgresSuite) isCodeExists() {
	ctx := context.Background()
	code := "f"

	res1, err := s.SUT.IsCodeExists(ctx, code)
	s.Assert().NoError(err)
	s.Assert().Equal(false, res1)

	_, err = s.db.Exec("insert into link(url, code, is_custom) values ($1, $2, $3)", "foo", code, false)
	s.Assert().NoError(err)

	res2, err := s.SUT.IsCodeExists(ctx, code)
	s.Assert().Equal(true, res2)
}

