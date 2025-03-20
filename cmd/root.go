package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/karthikkalarikal/Qube/models"
	"github.com/karthikkalarikal/Qube/pkg/util"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "separation",
		Short: "Connect your favorite celebrities.",
		Long:  `A cli application that let's you figure out the connection between actors by separation`,
		Run:   generate,
	}
	actor1     string
	actor2     string
	separation uint
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&actor1, "actor1", "a", "", "name an actor/celeb")
	rootCmd.PersistentFlags().StringVarP(&actor2, "actor2", "b", "", "name an actor/celeb")
	rootCmd.PersistentFlags().UintVarP(&separation, "separation", "s", 3, "separation between the actors/celebs")

}

func generate(_ *cobra.Command, args []string) {
	config := models.Config{
		Actor1:     actor1,
		Actor2:     actor2,
		Separation: separation,
	}
	var tar models.Actor
	err := util.GetByURL(config.Actor1, &tar)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(tar)
}
