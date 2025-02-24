package signer

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/signer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_signer_signing_job")
func DataSourceSigningJob() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceSigningJobRead,

		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"completed_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_invoker": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"requested_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"revocation_record": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"revoked_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"revoked_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"signature_expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signed_object": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"s3": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"s3": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSigningJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SignerConn()
	jobId := d.Get("job_id").(string)

	describeSigningJobOutput, err := conn.DescribeSigningJobWithContext(ctx, &signer.DescribeSigningJobInput{
		JobId: aws.String(jobId),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Signer signing job (%s): %s", d.Id(), err)
	}

	if err := d.Set("completed_at", aws.TimeValue(describeSigningJobOutput.CompletedAt).Format(time.RFC3339)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job completed at: %s", err)
	}

	if err := d.Set("created_at", aws.TimeValue(describeSigningJobOutput.CreatedAt).Format(time.RFC3339)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job created at: %s", err)
	}

	if err := d.Set("job_invoker", describeSigningJobOutput.JobInvoker); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job invoker: %s", err)
	}

	if err := d.Set("job_owner", describeSigningJobOutput.JobOwner); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job owner: %s", err)
	}

	if err := d.Set("platform_display_name", describeSigningJobOutput.PlatformDisplayName); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job platform display name: %s", err)
	}

	if err := d.Set("platform_id", describeSigningJobOutput.PlatformId); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job platform id: %s", err)
	}

	if err := d.Set("profile_name", describeSigningJobOutput.ProfileName); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job profile name: %s", err)
	}

	if err := d.Set("profile_version", describeSigningJobOutput.ProfileVersion); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job profile version: %s", err)
	}

	if err := d.Set("requested_by", describeSigningJobOutput.RequestedBy); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job requested by: %s", err)
	}

	if err := d.Set("revocation_record", flattenSigningJobRevocationRecord(describeSigningJobOutput.RevocationRecord)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job revocation record: %s", err)
	}

	signatureExpiresAt := ""
	if describeSigningJobOutput.SignatureExpiresAt != nil {
		signatureExpiresAt = aws.TimeValue(describeSigningJobOutput.SignatureExpiresAt).Format(time.RFC3339)
	}
	if err := d.Set("signature_expires_at", signatureExpiresAt); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job requested by: %s", err)
	}

	if err := d.Set("signed_object", flattenSigningJobSignedObject(describeSigningJobOutput.SignedObject)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job signed object: %s", err)
	}

	if err := d.Set("source", flattenSigningJobSource(describeSigningJobOutput.Source)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job source: %s", err)
	}

	if err := d.Set("status", describeSigningJobOutput.Status); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job status: %s", err)
	}

	if err := d.Set("status_reason", describeSigningJobOutput.StatusReason); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting signer signing job status reason: %s", err)
	}

	d.SetId(aws.StringValue(describeSigningJobOutput.JobId))

	return diags
}
