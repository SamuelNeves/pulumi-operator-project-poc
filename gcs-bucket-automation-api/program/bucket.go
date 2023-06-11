package program_bucket

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateBucketProgram(content string) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		bucketName := cfg.Require("bucket-name")
		location := cfg.Require("location")
		uniformBucketLevelAccess := cfg.RequireBool("uniform-level-access")
		// Create a GCP resource (Storage Bucket)
		bucket, err := storage.NewBucket(ctx, bucketName, &storage.BucketArgs{
			Location:                 pulumi.String(location),
			UniformBucketLevelAccess: pulumi.Bool(uniformBucketLevelAccess),
		})
		if err != nil {
			return err
		}

		// Export the DNS name of the bucket
		ctx.Export("bucket-name", bucket.Name)
		ctx.Export("bucket-url", bucket.Url)
		return nil
	}
}
