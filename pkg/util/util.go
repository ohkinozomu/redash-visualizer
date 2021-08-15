package util

import (
	"strings"

	"github.com/ohkinozomu/redash-client-go/redash"
)

func JoinDataSources(ds *[]redash.DataSource) string {
	dataSources := ""
	for _, v := range *ds {
		dataSources += v.Name
		dataSources += ","
	}
	return strings.TrimRight(dataSources, ",")
}

func JoinGroups(groups *[]redash.Group) string {
	g := ""
	for _, v := range *groups {
		g += v.Name
		g += ","
	}
	return strings.TrimRight(g, ",")
}

func JoinUsers(users *redash.UserList) string {
	u := ""
	for _, v := range users.Results {
		u += v.Name
		u += ","
	}
	return strings.TrimRight(u, ",")
}
