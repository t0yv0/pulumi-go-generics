package acm

import (
	pulumi "github.com/t0yv0/pulumi-go-generics/outputs"
)

type CertificateArgs struct {
	DomainName       pulumi.Output[string]
	ValidationMethod pulumi.Output[string]
}

type Certificate struct {
	DomainValidationOptions pulumi.Output[[]CertificateDomainValidationOption]
}

type CertificateDomainValidationOption struct {
	ResourceRecordValue *string
}

func NewCertificate(
	ctx *pulumi.Context,
	name string,
	args *CertificateArgs,
) (*Certificate, error) {

	v1 := "v1"
	v2 := "v2"

	return &Certificate{

		DomainValidationOptions: pulumi.Resolved(
			ctx,
			[]CertificateDomainValidationOption{
				{
					ResourceRecordValue: &v1,
				},
				{
					ResourceRecordValue: &v2,
				},
			},
		),
	}, nil
}
