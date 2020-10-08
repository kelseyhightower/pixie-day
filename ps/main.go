package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/kelseyhightower/pixie"
)

func main() {
	clusters, err := pixie.GetClusters()
	if err != nil {
		log.Fatal(err)
	}

	clusterIDs := make([]string, 0)
	for _, cluster := range clusters {
		if cluster.Status != 1 {
			continue
		}

		clusterIDs = append(clusterIDs, cluster.Name)
	}

	podstats, err := pixie.GetPodStats("kube-system", clusterIDs)
	if err != nil {
		log.Fatal(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 8, 0, '\t', 0)
	defer w.Flush()

	fmt.Fprintf(w, "\n%s\t%s\t%s\t%s\t%s\t%s\t%s\t", "C", "NS", "POD", "PID", "VSZ", "RSS", "Command")
	for _, v := range podstats {
		fmt.Fprintf(w, "\n%s\t%s\t%s\t%d\t%d\t%d\t%s\t", v.ClusterName, v.Namespace, v.PodName, v.Pid, v.VSZ, v.RSS, v.Command)
	}
}
