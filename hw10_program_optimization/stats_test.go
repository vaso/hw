//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty source", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetUsers(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	fixedIdData := `{"ID":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"ID":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"ID":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"ID":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"ID":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("get users with zero ID due to data source", func(t *testing.T) {
		result, err := getUsers(bytes.NewBufferString(data))
		require.NoError(t, err)
		expectedUserSlice := []User{
			{
				ID:       0,
				Name:     "Clarence Olson",
				Username: "RachelAdams",
				Email:    "RoseSmith@Browsecat.com",
				Phone:    "988-48-97",
				Password: "71kuz3gA5w",
				Address:  "Monterey Park 39",
			},
		}
		require.Equal(t, expectedUserSlice, result[2:3])
	})
	t.Run("get users with correct ID from fixed data source", func(t *testing.T) {
		result, err := getUsers(bytes.NewBufferString(fixedIdData))
		require.NoError(t, err)
		expectedUserSlice := []User{
			{
				ID:       3,
				Name:     "Clarence Olson",
				Username: "RachelAdams",
				Email:    "RoseSmith@Browsecat.com",
				Phone:    "988-48-97",
				Password: "71kuz3gA5w",
				Address:  "Monterey Park 39",
			},
		}
		require.Equal(t, expectedUserSlice, result[2:3])
	})
}

func TestCountDomains(t *testing.T) {
	var testUsers users
	testUsers[0] = User{
		ID:       3,
		Name:     "Clarence Olson",
		Username: "RachelAdams",
		Email:    "RoseSmith@Browsecat.com",
		Phone:    "988-48-97",
		Password: "71kuz3gA5w",
		Address:  "Monterey Park 39",
	}
	testUsers[1] = User{
		ID:       4,
		Name:     "Gregory Reid",
		Username: "tButler",
		Email:    "5Moore@Teklist.net",
		Phone:    "520-04-16",
		Password: "r639qLNu",
		Address:  "Sunfield Park 20",
	}

	t.Run("count domains", func(t *testing.T) {
		result, err := countDomains(testUsers, "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 1,
		}, result)
	})
	t.Run("count domains", func(t *testing.T) {
		result, err := countDomains(testUsers, "net")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"teklist.net": 1,
		}, result)
	})
}
