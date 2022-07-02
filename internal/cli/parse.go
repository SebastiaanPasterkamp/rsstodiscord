package cli

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/SebastiaanPasterkamp/go-cache"
	"github.com/alexflint/go-arg"
)

// Parse returns a Configuration object populated with the CLI options.
func Parse(flags []string) (*Configuration, error) {
	args := Arguments{}

	p, err := arg.NewParser(arg.Config{Program: flags[0]}, &args)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsingFailed, err)
	}

	if err = p.Parse(flags[1:]); err != nil {
		switch {
		case errors.Is(err, arg.ErrHelp):
			p.WriteHelp(os.Stdout)
		case errors.Is(err, arg.ErrVersion):
			log.Println(args.Version())
		default:
			p.WriteUsage(os.Stderr)
		}

		return nil, err
	}

	if p.Subcommand() == nil {
		p.WriteUsage(os.Stderr)
		return nil, arg.ErrHelp
	}

	mem := cache.Configuration{}

	if args.Address != "" {
		mem.RedisSettings = &args.RedisSettings
	} else {
		mem.InMemorySettings = &args.InMemorySettings
	}

	cfg := Configuration{
		Base:  args.Base,
		Cache: mem,
	}

	return &cfg, nil
}
