package cli

import (
	"fmt"
	"os"

	"github.com/j4ng5y/bulma/pkg/parser/puml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	ShortDescription = "A diagram-driven infrastructure-as-code tool."
)

var (
	Version = "dev"
)

func rootCMD() *cobra.Command {
	var (
		pumlFile, outDir string
		verbosity        int
	)
	var rootcmd = &cobra.Command{
		Use:     "bulma",
		Version: Version,
		Short:   ShortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			switch verbosity {
			case 1:
				log.Info().Msgf("setting log level to %s", zerolog.DebugLevel.String())
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			case 2:
				log.Info().Msgf("setting log level to %s", zerolog.TraceLevel.String())
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			default:
				log.Info().Msgf("setting log level to %s", zerolog.InfoLevel.String())
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}

			if pumlFile == "" {
				log.Fatal().Err(fmt.Errorf("provided puml file must not be blank")).Send()
			}

			log.Info().Msgf("attempting to parse %s", pumlFile)
			log.Debug().Msgf("output directory set to %s", outDir)

			p, err := puml.NewParser(pumlFile)
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			if err := p.Parse(); err != nil {
				log.Fatal().Err(err).Send()
			}
		},
	}

	rootcmd.PersistentFlags().StringVarP(&pumlFile, "puml.file", "f", "", "The diagram file to parse.")
	rootcmd.MarkFlagFilename("puml.file", "puml")
	rootcmd.PersistentFlags().StringVarP(&outDir, "out.dir", "o", "./generated/", "The directory to output generated code into.")
	rootcmd.MarkFlagDirname("out")
	rootcmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "The verbosity of the logs.")

	return rootcmd
}

func Run() {
	if err := rootCMD().Execute(); err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}
}
