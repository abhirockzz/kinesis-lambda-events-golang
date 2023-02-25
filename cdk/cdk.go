package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awskinesis"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambdaeventsources"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const functionDir = "../function"

type KinesisLambdaGolangStackProps struct {
	awscdk.StackProps
}

func NewKinesisLambdaGolangStack(scope constructs.Construct, id string, props *KinesisLambdaGolangStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("dynamodb-table"),
		&awsdynamodb.TableProps{
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String("email"),
				Type: awsdynamodb.AttributeType_STRING},
		})

	table.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY)

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("kinesis-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime:     awslambda.Runtime_GO_1_X(),
			Entry:       jsii.String(functionDir),
			Environment: &map[string]*string{"TABLE_NAME": table.TableName()},
		})

	table.GrantWriteData(function)

	kinesisStream := awskinesis.NewStream(stack, jsii.String("lambda-test-stream"), nil)

	function.AddEventSource(awslambdaeventsources.NewKinesisEventSource(kinesisStream, &awslambdaeventsources.KinesisEventSourceProps{
		StartingPosition: awslambda.StartingPosition_LATEST,
	}))

	awscdk.NewCfnOutput(stack, jsii.String("kinesis-stream-name"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("kinesis-stream-name"),
			Value:      kinesisStream.StreamName()})

	awscdk.NewCfnOutput(stack, jsii.String("dynamodb-table-name"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("dynamodb-table-name"),
			Value:      table.TableName()})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewKinesisLambdaGolangStack(app, "KinesisLambdaGolangStack", &KinesisLambdaGolangStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
