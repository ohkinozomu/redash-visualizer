package cmd

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/snowplow-devops/redash-client-go/redash"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().String("host", "", "host")
	runCmd.PersistentFlags().String("api-key", "", "api-key")
	runCmd.PersistentFlags().String("file", "graph.png", "file name")
}

func addGroupsNodes(graph *cgraph.Graph, groups *[]redash.Group) error {
	for _, v := range *groups {
		_, err := graph.CreateNode(v.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func addDataSourcesNodes(c *redash.Client, graph *cgraph.Graph, ds *[]redash.DataSource, groups *[]redash.Group) error {
	for _, v := range *ds {
		dsn, err := graph.CreateNode(v.Name)
		if err != nil {
			return err
		}

		ads, err := c.GetDataSource(v.ID)
		if err != nil {
			return err
		}
		for i := range ads.Groups {
			for _, g := range *groups {
				if g.ID == i {
					gn, err := graph.Node(g.Name)
					if err != nil {
						log.Fatal(err)
					}
					_, err = graph.CreateEdge("", gn, dsn)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
	return nil
}

func addUsersNodes(graph *cgraph.Graph, users *redash.UserList) error {
	for _, v := range users.Results {
		un, err := graph.CreateNode(v.Name)
		if err != nil {
			return err
		}

		for _, g := range v.Groups {
			gn, err := graph.Node(g.Name)
			if err != nil {
				log.Fatal(err)
			}
			_, err = graph.CreateEdge("", un, gn)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}

func renderGraph(c *redash.Client, ds *[]redash.DataSource, groups *[]redash.Group, users *redash.UserList, fileName string) error {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	err = addGroupsNodes(graph, groups)
	if err != nil {
		return err
	}

	err = addDataSourcesNodes(c, graph, ds, groups)
	if err != nil {
		return err
	}

	err = addUsersNodes(graph, users)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
		return err
	}

	if err := g.RenderFilename(graph, graphviz.PNG, fileName); err != nil {
		return err
	}

	return nil
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run redash-visualizer",
	Long:  `run redash-visualizer`,
	Run: func(cmd *cobra.Command, args []string) {
		host, err := cmd.PersistentFlags().GetString("host")
		if err != nil {
			panic(err)
		}

		apiKey, err := cmd.PersistentFlags().GetString("api-key")
		if err != nil {
			panic(err)
		}

		fileName, err := cmd.PersistentFlags().GetString("file")
		if err != nil {
			panic(err)
		}

		if !(strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://")) {
			host = "http://" + host
		}

		config := redash.Config{
			RedashURI: host,
			APIKey:    apiKey,
		}

		c, err := redash.NewClient(&config)
		if err != nil {
			log.Fatal(fmt.Errorf("loading client error: %q", err))
			return
		}

		ds, err := c.GetDataSources()
		if err != nil {
			log.Fatal(err)
		}

		groups, err := c.GetGroups()
		if err != nil {
			log.Fatal(err)
		}

		users, err := c.GetUsers()
		if err != nil {
			log.Fatal(err)
		}

		err = renderGraph(c, ds, groups, users, fileName)
		if err != nil {
			log.Fatal(err)
		}
	},
}
