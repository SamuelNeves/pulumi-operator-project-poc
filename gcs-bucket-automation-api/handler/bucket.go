package handler

import (
	program_bucket "automation-api/program"
	"automation-api/requests"
	"automation-api/responses"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"net/http"
	"os"
	"strconv"
)

type Bucket interface {
	Get(echo.Context) error
	Post(echo.Context) error
	Delete(echo.Context) error
	List(echo.Context) error
}

type bucketHandler struct {
	project string
}

// Get godoc
// @Summary Get a specific bucket stack.
// @Description get the status of server.
// @Tags buckets
// @Accept */*
// @Produce json
// @Security BearerAuth
// @Param stack_id  path string true "Stack Name"
// @Success 200 {object} responses.BucketResponse
// @Router /bucket/{stack_id} [get]
func (b bucketHandler) Get(c echo.Context) error {
	stackName := c.Param("name")
	var program pulumi.RunFunc = nil
	ctx := context.Background()
	s, err := auto.SelectStackInlineSource(ctx, stackName, b.project, program)
	if err != nil {
		// if the stack doesn't already exist, 404
		if auto.IsSelectStack404Error(err) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("stack %q not found", stackName))

		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

	// fetch the outputs from the stack
	outs, err := s.Outputs(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

	response := &responses.BucketResponse{
		ID:   stackName,
		Name: outs["bucket-name"].Value.(string),
		URL:  outs["bucket-url"].Value.(string),
	}
	return c.JSON(http.StatusOK, response)
}

// Post godoc
// @Summary Create a Bucket Stacks.
// @Description TODO
// @Tags buckets
// @Accept */*
// @Produce json
// @Param request body requests.BucketRequest true "bucket params"
// @Security BearerAuth
// @Router /bucket [post]
func (b bucketHandler) Post(c echo.Context) error {
	var err error
	bucket := &requests.BucketRequest{UniformBucketLevelAccess: true, PublicAccessPrevention: "enforced"}
	if err = c.Bind(bucket); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(bucket); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ctx := context.Background()

	stackName := bucket.Name
	program := program_bucket.CreateBucketProgram("")

	s, err := auto.NewStackInlineSource(ctx, stackName, b.project, program)
	if err != nil {
		// if stack already exists, 409
		if auto.IsCreateStack409Error(err) {
			//return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("stack %q already exists", stackName))

		} else {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	c.Logger().Printf(fmt.Sprintf("%+v", bucket))
	s.SetConfig(ctx, "gcp:project", auto.ConfigValue{Value: bucket.Project})
	s.SetConfig(ctx, "bucket-automation-api:bucket-name", auto.ConfigValue{Value: bucket.Name})
	s.SetConfig(ctx, "location", auto.ConfigValue{Value: bucket.Location})
	s.SetConfig(ctx, "bucket-automation-api:uniform-level-access", auto.ConfigValue{Value: strconv.FormatBool(bucket.UniformBucketLevelAccess)})

	// deploy the stack
	// we'll write all of the update logs to st	out so we can watch requests get processed
	upRes, err := s.Up(ctx, optup.ProgressStreams(os.Stdout))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := &responses.BucketResponse{
		ID:   stackName,
		Name: upRes.Outputs["bucket-name"].Value.(string),
		URL:  upRes.Outputs["bucket-url"].Value.(string),
	}
	return c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Delete a specific bucket stack.
// @Description get the status of server.
// @Tags buckets
// @Accept */*
// @Produce json
// @Security BearerAuth
// @Param stack_id  path string true "Stack Name"
// @Router /bucket/{stack_id} [delete]
func (b bucketHandler) Delete(c echo.Context) error {
	stackName := c.Param("name")
	var program pulumi.RunFunc = nil
	ctx := context.Background()
	s, err := auto.SelectStackInlineSource(ctx, stackName, b.project, program)
	if err != nil {
		// if the stack doesn't already exist, 404
		if auto.IsSelectStack404Error(err) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("stack %q not found", stackName))

		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}
	_, err = s.Destroy(ctx, optdestroy.ProgressStreams(os.Stdout))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}
	s.Workspace().RemoveStack(ctx, stackName)
	return c.JSON(http.StatusOK, "stack deleted with success")

}

// List godoc
// @Summary List Buckets Stacks.
// @Description List the buckets stacks.
// @Tags buckets
// @Accept */*
// @Produce json
// @Security BearerAuth
// @Success 200 {array} responses.ListBucketsResponse
// @Router /bucket [get]
func (b bucketHandler) List(c echo.Context) error {
	ctx := context.Background()
	// set up a workspace with only enough information for the list stack operations
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:    tokens.PackageName(b.project),
		Runtime: workspace.NewProjectRuntimeInfo("go", nil),
	}))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}
	stacks, err := ws.ListStacks(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}
	var ids []string
	for _, stack := range stacks {
		ids = append(ids, stack.Name)
	}

	response := &responses.ListBucketsResponse{
		IDs: ids,
	}

	return c.JSON(http.StatusOK, response)
}

func NewBucketHandler(project string) Bucket {
	return bucketHandler{project: project}
}
