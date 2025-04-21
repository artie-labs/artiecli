package internal

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

type Command interface {
	Execute(ctx context.Context, client ArtieClient) error
	ParseFlags(fs *flag.FlagSet, args []string) error
}

type ListDeploymentsCommand struct{}

func (ListDeploymentsCommand) Execute(ctx context.Context, client ArtieClient) error {
	return client.ListDeployments(ctx)
}

func (ListDeploymentsCommand) ParseFlags(_ *flag.FlagSet, _ []string) error {
	return nil
}

type GetDeploymentByUUIDCommand struct {
	DeploymentUUID string
}

func (g GetDeploymentByUUIDCommand) Execute(ctx context.Context, client ArtieClient) error {
	return client.GetDeploymentByUUID(ctx, g.DeploymentUUID)
}

func (g *GetDeploymentByUUIDCommand) ParseFlags(fs *flag.FlagSet, args []string) error {
	var deploymentUUID string
	fs.StringVar(&deploymentUUID, "deployment-uuid", "", "UUID of the deployment to get")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if deploymentUUID == "" {
		return fmt.Errorf("--deployment-uuid is required")
	}

	g.DeploymentUUID = deploymentUUID
	return nil
}

type CancelDeploymentBackfillCommand struct {
	DeploymentUUID string
	TableUUIDs     []string
}

func (c CancelDeploymentBackfillCommand) Execute(ctx context.Context, client ArtieClient) error {
	if err := client.CancelDeploymentBackfill(ctx, c.DeploymentUUID, c.TableUUIDs); err != nil {
		return fmt.Errorf("failed to cancel deployment backfill: %w", err)
	}

	slog.Info("Deployment backfill cancelled",
		slog.String("deploymentUUID", c.DeploymentUUID),
		slog.String("tableUUIDs", strings.Join(c.TableUUIDs, ",")))
	return nil
}

func (c *CancelDeploymentBackfillCommand) ParseFlags(fs *flag.FlagSet, args []string) error {
	var deploymentUUID string
	var tableUUIDsString string
	fs.StringVar(&deploymentUUID, "deployment-uuid", "", "UUID of the deployment to cancel backfill for")
	fs.StringVar(&tableUUIDsString, "table-uuids", "", "Comma-separated list of table UUIDs to cancel backfill for")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if deploymentUUID == "" {
		return fmt.Errorf("--deployment-uuid is required")
	}
	if tableUUIDsString == "" {
		return fmt.Errorf("--table-uuids is required")
	}

	c.DeploymentUUID = deploymentUUID
	c.TableUUIDs = strings.Split(tableUUIDsString, ",")

	if len(c.TableUUIDs) == 0 {
		return fmt.Errorf("--table-uuids must contain at least one table UUID")
	}

	return nil
}

type SourceReaderDeployCommand struct {
	SourceReaderUUID uuid.UUID
}

func (s SourceReaderDeployCommand) Execute(ctx context.Context, client ArtieClient) error {
	return client.DeploySourceReader(ctx, s.SourceReaderUUID)
}

func (s *SourceReaderDeployCommand) ParseFlags(fs *flag.FlagSet, args []string) error {
	var sourceReaderUUID string
	fs.StringVar(&sourceReaderUUID, "source-reader-uuid", "", "UUID of the source reader to deploy")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	parsedUUID, err := uuid.Parse(sourceReaderUUID)
	if err != nil {
		return fmt.Errorf("failed to parse source reader UUID: %w", err)
	}

	s.SourceReaderUUID = parsedUUID
	return nil
}

type DeployDeploymentCommand struct {
	deploymentUUID uuid.UUID
}

func (d DeployDeploymentCommand) Execute(ctx context.Context, client ArtieClient) error {
	return client.DeployDeployment(ctx, d.deploymentUUID)
}

func (d *DeployDeploymentCommand) ParseFlags(fs *flag.FlagSet, args []string) error {
	var deploymentUUID string
	fs.StringVar(&deploymentUUID, "deployment-uuid", "", "UUID of the deployment to deploy")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	parsedUUID, err := uuid.Parse(deploymentUUID)
	if err != nil {
		return fmt.Errorf("failed to parse deployment UUID: %w", err)
	}

	d.deploymentUUID = parsedUUID
	return nil
}

func ParseCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("no command provided")
	}

	var cmd Command
	switch args[1] {
	case "list-deployments":
		cmd = &ListDeploymentsCommand{}
	case "get-deployment":
		cmd = &GetDeploymentByUUIDCommand{}
	case "cancel-deployment-backfill":
		cmd = &CancelDeploymentBackfillCommand{}
	case "deploy-deployment":
		cmd = &DeployDeploymentCommand{}
	case "deploy-source-reader":
		cmd = &SourceReaderDeployCommand{}
	default:
		return nil, fmt.Errorf("unknown command: %q", args[1])
	}

	if err := cmd.ParseFlags(flag.NewFlagSet("", flag.ExitOnError), args[2:]); err != nil {
		return nil, err
	}

	return cmd, nil
}
