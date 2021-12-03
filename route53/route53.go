package route53

import (
	pulumi "github.com/t0yv0/pulumi-go-generics/outputs"
)

type RecordArgs struct {
	Records pulumi.Output[[]string]
}

type Record struct {}

func NewRecord(ctx *pulumi.Context, name string, args *RecordArgs) (*Record, error) {
	return &Record{}, nil
}
