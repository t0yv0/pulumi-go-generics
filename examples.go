// Inspired by:
//
// https://www.pulumi.com/docs/intro/concepts/inputs-outputs/

package main

import (
	"fmt"

	acm "github.com/t0yv0/pulumi-go-generics/acm"
	pulumi "github.com/t0yv0/pulumi-go-generics/outputs"
	route53 "github.com/t0yv0/pulumi-go-generics/route53"
)

func ApplyExample(ctx *pulumi.Context) {
	fmt.Printf("\n\nApplyExample\n")
	vpc := &Vpc{DnsName: pulumi.Resolved(ctx, "mycompany.com")}

	url := pulumi.Apply(ctx, func(t *pulumi.T) string {
		return "https://" + pulumi.Eval(t, vpc.DnsName)
	})

	pulumi.Debug(url)
}

type Vpc struct {
	DnsName pulumi.Output[string]
}

func AllExample(ctx *pulumi.Context) {
	fmt.Printf("\n\nAllExample\n")
	sqlServer := &Named{pulumi.Resolved(ctx, "sql-server")}
	database := &Named{pulumi.Resolved(ctx, "my-database")}

	connectionString := pulumi.Apply(ctx, func(t *pulumi.T) string {
		server := pulumi.Eval(t, sqlServer.Name)
		db := pulumi.Eval(t, database.Name)
		return fmt.Sprintf("Server=tcp:%s.database.windows.net;initial catalog=%s...", server, db)
	})

	pulumi.Debug(connectionString)
}

type Named struct {
	Name pulumi.Output[string]
}

func LiftingExample(ctx *pulumi.Context) {
	fmt.Printf("\n\nLifting Example\n")
	cert, err := acm.NewCertificate(ctx, "cert", &acm.CertificateArgs{
		DomainName:       pulumi.Resolved(ctx, "example"),
		ValidationMethod: pulumi.Resolved(ctx, "DNS"),
	})
	if err != nil {
		panic(err)
	}

	// No lifting in this approach, but Apply/Eval recovers normal
	// value that can be accessed using normal property accessors.
	record, err := route53.NewRecord(ctx, "validation", &route53.RecordArgs{
		Records: pulumi.Apply(ctx, func(t *pulumi.T) []string {
			opts := pulumi.Eval(t, cert.DomainValidationOptions)
			return []string{*opts[0].ResourceRecordValue}
		}),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("record: %v\n", record)
}

// TODO sprintf, interpolation

func Examples(ctx *pulumi.Context) {
	ApplyExample(ctx)
	AllExample(ctx)
	LiftingExample(ctx)
}
