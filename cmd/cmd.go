package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/hsmtkk/aukabucomgo/base"
	"github.com/hsmtkk/aukabucomgo/info/boardget"
	"github.com/hsmtkk/aukabucomgo/info/positionsget"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use: "selloptiondelta",
	Run: run,
}

var test bool

func init() {
	RootCommand.Flags().BoolVar(&test, "test", false, "test")
}

func run(cmd *cobra.Command, args []string) {
	apiPassword := os.Getenv("API_PASSWORD")
	if apiPassword == "" {
		log.Fatal("env var API_PASSWORD is not defined")
	}
	positions, err := getOptionPositions(apiPassword)
	totalDelta := 0.0
	if err != nil {
		log.Fatal(err)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Quarter", "Delta"})
	for _, pos := range positions {
		t.AppendRow([]interface{}{pos.symbolName, pos.leavesQty, pos.delta})
		totalDelta += float64(pos.leavesQty) * pos.delta
	}
	t.Render()
	fmt.Printf("Total delta: %f\n", totalDelta)
}

type optionPosition struct {
	symbolName string
	leavesQty  int
	delta      float64
}

func getOptionPositions(apiPassword string) ([]optionPosition, error) {
	var env base.Environment
	if test {
		env = base.TEST
	} else {
		env = base.PRODUCTION
	}
	baseClient, err := base.New(env, apiPassword)
	if err != nil {
		return nil, err
	}
	posClient := positionsget.New(baseClient)
	boardClient := boardget.New(baseClient)

	results := []optionPosition{}
	posResp, err := posClient.PositionsGet(positionsget.OPTION, positionsget.SELL)
	if err != nil {
		return nil, err
	}
	for _, pos := range posResp {
		boardResp, err := boardClient.BoardGet(pos.Symbol, boardget.ALL_DAY)
		if err != nil {
			return nil, err
		}
		results = append(results, optionPosition{symbolName: pos.SymbolName, leavesQty: int(pos.LeavesQty), delta: boardResp.Delta})
	}
	return results, nil
}
