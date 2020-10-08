package pixie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Clusters map[string]*Cluster

type Cluster struct {
	Name   string `json:"ClusterName"`
	ID     string `json:"ID"`
	Status int64  `json:"Status"`
}

type PodStats struct {
	Pid         int64  `json:"pid"`
	Container   string `json:"container"`
	Pod         string `json:"pod"`
	RSS         int64  `json:"rss"`
	VSZ         int64  `json:"vsz"`
	Node        string `json:"node"`
	Namespace   string `json:"namespace"`
	Command     string `json:"cmd"`
	CPUTime     int64  `json:"time"`
	ClusterName string
	PodName     string
}

func GetClusters() (Clusters, error) {
	clusters := make(Clusters)

	cmd := exec.Command("px", "get", "clusters", "-o", "json")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(&out)
	for decoder.More() {
		var c Cluster
		err := decoder.Decode(&c)
		if err != nil {
			return nil, err
		}

		clusters[c.Name] = &c
	}

	return clusters, nil
}

func GetPodStats(namespace string, clusters []string) ([]PodStats, error) {
	if namespace == "" {
		namespace = "default"
	}

	podstats := make([]PodStats, 0)

	for _, cluster := range clusters {
		c := cluster
		fmt.Printf("Getting data for %s\n", c)

		cmd := exec.Command("px", "run", "-o", "json", "-f", "-", "-c", c, "--", "--namespace", namespace)

		var out bytes.Buffer
		cmd.Stdout = &out

		var errOut bytes.Buffer
		cmd.Stderr = &errOut

		cmd.Stdin = strings.NewReader(podStatsQuery)

		err := cmd.Run()
		if err != nil {
			fmt.Println(string(errOut.Bytes()))
			return nil, err
		}

		decoder := json.NewDecoder(&out)
		for decoder.More() {
			var stats PodStats
			err := decoder.Decode(&stats)
			if err != nil {
				return nil, err
			}
			s := strings.Split(stats.Pod, "/")

			stats.ClusterName = c
			stats.PodName = s[1]
			podstats = append(podstats, stats)
		}
	}

	return podstats, nil
}
