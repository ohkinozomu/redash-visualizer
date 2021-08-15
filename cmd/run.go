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
	runCmd.PersistentFlags().String("layout", "dot", "layout")
	runCmd.PersistentFlags().String("format", "png", "file format")
}

func addGroupsNodes(graph *cgraph.Graph, groups *[]redash.Group) error {
	for _, v := range *groups {
		gn, err := graph.CreateNode(v.Name)
		if err != nil {
			return err
		}
		gn.SetFillColor("lightblue1")
		gn.SetStyle(cgraph.FilledNodeStyle)
	}
	return nil
}

func addDataSourcesNodes(c *redash.Client, graph *cgraph.Graph, ds *[]redash.DataSource, groups *[]redash.Group) error {
	for _, v := range *ds {
		dsn, err := graph.CreateNode(v.Name)
		if err != nil {
			return err
		}

		dsn.SetFillColor("lightgrey")
		dsn.SetStyle(cgraph.FilledNodeStyle)

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

		un.SetFillColor("lightpink")
		un.SetStyle(cgraph.FilledNodeStyle)

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

func renderGraph(c *redash.Client, ds *[]redash.DataSource, groups *[]redash.Group, users *redash.UserList, fileName, layout, fileFormat string) error {
	g := graphviz.New()
	g.SetLayout(graphviz.Layout(layout))

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
	if err := g.Render(graph, graphviz.Format(fileFormat), &buf); err != nil {
		return err
	}

	if err := g.RenderFilename(graph, graphviz.Format(fileFormat), fileName); err != nil {
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

		layout, err := cmd.PersistentFlags().GetString("layout")
		if err != nil {
			panic(err)
		}

		fileFormat, err := cmd.PersistentFlags().GetString("format")
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

		err = renderGraph(c, ds, groups, users, fileName, layout, fileFormat)
		if err != nil {
			log.Fatal(err)
		}
	},
}
