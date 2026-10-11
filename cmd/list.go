package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/moon9t/svcmgr/internal/config"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var (
	groupBy    string
	showHidden bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List cosmic services",
	Run:   runList,
}

func init() {
	listCmd.Flags().StringVarP(&groupBy, "group-by", "g", "", "Group by (type|category)")
	listCmd.Flags().BoolVar(&showHidden, "show-hidden", false, "Reveal celestial secrets")
	// The command registration happens in root.go for clarity.
	AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	services, err := config.LoadServices()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load cosmic services")
	}
	// Service has no hidden indicator, so --show-hidden currently has no effect.
	if groupBy == "" {
		printServices(services)
		return
	}
	if groupBy != "type" && groupBy != "category" {
		log.Fatal().Msg("--group-by must be type or category")
	}
	groups := make(map[string][]config.Service)
	for _, service := range services {
		key := service.Type
		if groupBy == "category" {
			key = service.Category
		}
		groups[key] = append(groups[key], service)
	}
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for i, key := range keys {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s: \n", key)
		printServices(groups[key])
	}
}

func printServices(services []config.Service) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tCATEGORY\tHOST\tPORT")
	for _, service := range services {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n", service.Name, service.Type, service.Category, service.Host, service.Port)
	}
	w.Flush()
}
