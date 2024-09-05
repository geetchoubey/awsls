package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2Instance struct {
	svc      *ec2.EC2
	instance *ec2.Instance
}

func init() {
	register("EC2Instance", ListEC2Instances)
}

func ListEC2Instances(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)
	params := &ec2.DescribeInstancesInput{}
	resources := make([]Resource, 0)
	for {
		if resp, err := svc.DescribeInstances(params); err != nil {
			return nil, err
		} else {
			for _, reservation := range resp.Reservations {
				for _, instance := range reservation.Instances {
					resources = append(resources, &EC2Instance{
						svc:      svc,
						instance: instance,
					})
				}
			}

			if resp.NextToken == nil {
				break
			}

			params = &ec2.DescribeInstancesInput{
				NextToken: resp.NextToken,
			}
		}

	}
	return resources, nil
}

func (i *EC2Instance) Describe() (interface{}, error) {
	return i.instance, nil
}
