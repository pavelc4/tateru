package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pavelc4/tateru/pkg/config"
	"github.com/pavelc4/tateru/pkg/dag"
	"github.com/pavelc4/tateru/pkg/executor"
	"github.com/pavelc4/tateru/pkg/reporter"
	"github.com/pavelc4/tateru/pkg/stages"
	"github.com/pavelc4/tateru/pkg/toolchain"
)

type Options struct {
	Target   string
	DryRun   bool
	Jobs     int
	Verbose  bool
	CacheDir string
}

func Run(args []string) error {
	fs := flag.NewFlagSet("tateru", flag.ContinueOnError)
	opts := Options{}

	fs.BoolVar(&opts.DryRun, "dry-run", false, "print stages without executing")
	fs.IntVar(&opts.Jobs, "j", 4, "parallel jobs")
	fs.BoolVar(&opts.Verbose, "v", false, "verbose output")
	fs.StringVar(&opts.CacheDir, "cache-dir", ".tateru-cache", "cache directory")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "tateru — Android/GKI kernel build tool\n\n")
		fmt.Fprintf(os.Stderr, "usage:\n")
		fmt.Fprintf(os.Stderr, "  tateru [flags] <target>\n\n")
		fmt.Fprintf(os.Stderr, "examples:\n")
		fmt.Fprintf(os.Stderr, "  tateru marble\n")
		fmt.Fprintf(os.Stderr, "  tateru --dry-run marble\n")
		fmt.Fprintf(os.Stderr, "  tateru -j 8 marble\n\n")
		fmt.Fprintf(os.Stderr, "flags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil // --help sudah print usage, exit clean
		}
		return err
	}

	if fs.NArg() < 1 {
		fs.Usage()
		return fmt.Errorf("target required")
	}
	opts.Target = fs.Arg(0)

	return runBuild(opts)
}

func runBuild(opts Options) error {
	cwd, _ := os.Getwd()
	wsRoot, err := config.FindWorkspaceRoot(cwd)
	if err != nil {
		return err
	}

	cfgPath, _, err := config.ResolveTarget(wsRoot, opts.Target)
	if err != nil {
		return err
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	tc, err := toolchain.DetectClang(cfg.Toolchain.Clang)
	if err != nil {
		return fmt.Errorf("toolchain: %w", err)
	}
	env := tc.Env()

	fmt.Printf("tateru  target=%s  clang=%s  dry-run=%v\n\n",
		opts.Target, tc.ClangBin, opts.DryRun)

	g := dag.New()
	g.Add(stages.NewDefconfig(cfg, env))
	g.Add(stages.NewKernel(cfg, env))

	rep := reporter.NewPlain()
	ex := executor.New(g, executor.Options{
		DryRun:   opts.DryRun,
		Jobs:     opts.Jobs,
		CacheDir: opts.CacheDir,
	})

	start := time.Now()
	runErr := ex.Run(context.Background())
	elapsed := time.Since(start).Round(time.Millisecond).String()

	rep.Summary(ex.Results(), elapsed)

	return runErr
}
