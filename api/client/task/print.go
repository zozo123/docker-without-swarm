package task

import (
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/api/client/idresolver"
	"github.com/docker/go-units"
)

const (
	psTaskItemFmt = "%s\t%s\t%s\t%s\t%s\t%s %s ago\t%s\n"
	maxErrLength  = 30
)

func (t tasksBySlot) Len() int {
	return len(t)
}

func (t tasksBySlot) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t tasksBySlot) Less(i, j int) bool {
	// Sort by slot.
	if t[i].Slot != t[j].Slot {
		return t[i].Slot < t[j].Slot
	}

	// If same slot, sort by most recent.
	return t[j].Meta.CreatedAt.Before(t[i].CreatedAt)
}

// Print task information in a table format
func Print(dockerCli *client.DockerCli, ctx context.Context, resolver *idresolver.IDResolver) error {
	writer := tabwriter.NewWriter(dockerCli.Out(), 0, 4, 2, ' ', 0)

	// Ignore flushing errors
	defer writer.Flush()
	fmt.Fprintln(writer, strings.Join([]string{"ID", "NAME", "IMAGE", "NODE", "DESIRED STATE", "CURRENT STATE", "ERROR"}, "\t"))

	return nil
}
