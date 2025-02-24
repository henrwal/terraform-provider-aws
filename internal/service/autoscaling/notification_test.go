package autoscaling_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccAutoScalingNotification_ASG_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var asgn autoscaling.DescribeNotificationConfigurationsOutput

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, autoscaling.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckASGNDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGNotificationExists(ctx, "aws_autoscaling_notification.example", []string{"foobar1-terraform-test-" + rName}, &asgn),
					testAccCheckASGNotificationAttributes("aws_autoscaling_notification.example", &asgn),
				),
			},
		},
	})
}

func TestAccAutoScalingNotification_ASG_update(t *testing.T) {
	ctx := acctest.Context(t)
	var asgn autoscaling.DescribeNotificationConfigurationsOutput

	rName := sdkacctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, autoscaling.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckASGNDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGNotificationExists(ctx, "aws_autoscaling_notification.example", []string{"foobar1-terraform-test-" + rName}, &asgn),
					testAccCheckASGNotificationAttributes("aws_autoscaling_notification.example", &asgn),
				),
			},

			{
				Config: testAccNotificationConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGNotificationExists(ctx, "aws_autoscaling_notification.example", []string{"foobar1-terraform-test-" + rName, "barfoo-terraform-test-" + rName}, &asgn),
					testAccCheckASGNotificationAttributes("aws_autoscaling_notification.example", &asgn),
				),
			},
		},
	})
}

func TestAccAutoScalingNotification_ASG_pagination(t *testing.T) {
	ctx := acctest.Context(t)
	var asgn autoscaling.DescribeNotificationConfigurationsOutput

	resourceName := "aws_autoscaling_notification.example"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, autoscaling.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckASGNDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationConfig_pagination(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGNotificationExists(ctx, resourceName,
						[]string{
							"foobar3-terraform-test-0",
							"foobar3-terraform-test-1",
							"foobar3-terraform-test-2",
							"foobar3-terraform-test-3",
							"foobar3-terraform-test-4",
							"foobar3-terraform-test-5",
							"foobar3-terraform-test-6",
							"foobar3-terraform-test-7",
							"foobar3-terraform-test-8",
							"foobar3-terraform-test-9",
							"foobar3-terraform-test-10",
							"foobar3-terraform-test-11",
							"foobar3-terraform-test-12",
							"foobar3-terraform-test-13",
							"foobar3-terraform-test-14",
							"foobar3-terraform-test-15",
							"foobar3-terraform-test-16",
							"foobar3-terraform-test-17",
							"foobar3-terraform-test-18",
							"foobar3-terraform-test-19",
						}, &asgn),
					testAccCheckASGNotificationAttributes(resourceName, &asgn),
				),
			},
		},
	})
}

func testAccCheckASGNotificationExists(ctx context.Context, n string, groups []string, asgn *autoscaling.DescribeNotificationConfigurationsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ASG Notification ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AutoScalingConn()
		opts := &autoscaling.DescribeNotificationConfigurationsInput{
			AutoScalingGroupNames: aws.StringSlice(groups),
			MaxRecords:            aws.Int64(100),
		}

		resp, err := conn.DescribeNotificationConfigurationsWithContext(ctx, opts)
		if err != nil {
			return fmt.Errorf("Error describing notifications: %s", err)
		}

		*asgn = *resp

		return nil
	}
}

func testAccCheckASGNDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_autoscaling_notification" {
				continue
			}

			groups := []*string{aws.String("foobar1-terraform-test")}
			conn := acctest.Provider.Meta().(*conns.AWSClient).AutoScalingConn()
			opts := &autoscaling.DescribeNotificationConfigurationsInput{
				AutoScalingGroupNames: groups,
			}

			resp, err := conn.DescribeNotificationConfigurationsWithContext(ctx, opts)
			if err != nil {
				return fmt.Errorf("Error describing notifications")
			}

			if len(resp.NotificationConfigurations) != 0 {
				return fmt.Errorf("Error finding notification descriptions")
			}
		}
		return nil
	}
}

func testAccCheckASGNotificationAttributes(n string, asgn *autoscaling.DescribeNotificationConfigurationsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ASG Notification ID is set")
		}

		if len(asgn.NotificationConfigurations) == 0 {
			return fmt.Errorf("Error: no ASG Notifications found")
		}

		// build a unique list of groups, notification types
		gRaw := make(map[string]bool)
		nRaw := make(map[string]bool)

		for _, n := range asgn.NotificationConfigurations {
			if *n.TopicARN == rs.Primary.Attributes["topic_arn"] {
				gRaw[*n.AutoScalingGroupName] = true
				nRaw[*n.NotificationType] = true
			}
		}

		// Grab the keys here as the list of Groups
		var gList []string
		for k := range gRaw {
			gList = append(gList, k)
		}

		// Grab the keys here as the list of Types
		var nList []string
		for k := range nRaw {
			nList = append(nList, k)
		}

		typeCount, _ := strconv.Atoi(rs.Primary.Attributes["notifications.#"])

		if len(nList) != typeCount {
			return fmt.Errorf("Error: Bad ASG Notification count, expected (%d), got (%d)", typeCount, len(nList))
		}

		groupCount, _ := strconv.Atoi(rs.Primary.Attributes["group_names.#"])

		if len(gList) != groupCount {
			return fmt.Errorf("Error: Bad ASG Group count, expected (%d), got (%d)", typeCount, len(gList))
		}

		return nil
	}
}

func testAccNotificationConfig_basic(rName string) string {
	return acctest.ConfigLatestAmazonLinuxHVMEBSAMI() + fmt.Sprintf(`
resource "aws_sns_topic" "topic_example" {
  name = "user-updates-topic-%s"
}

resource "aws_launch_configuration" "foobar" {
  name          = "foobarautoscaling-terraform-test-%s"
  image_id      = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_autoscaling_group" "bar" {
  availability_zones        = [data.aws_availability_zones.available.names[1]]
  name                      = "foobar1-terraform-test-%s"
  max_size                  = 1
  min_size                  = 1
  health_check_grace_period = 100
  health_check_type         = "ELB"
  desired_capacity          = 1
  force_delete              = true
  termination_policies      = ["OldestInstance"]
  launch_configuration      = aws_launch_configuration.foobar.name
}

resource "aws_autoscaling_notification" "example" {
  group_names = [aws_autoscaling_group.bar.name]

  notifications = [
    "autoscaling:EC2_INSTANCE_LAUNCH",
    "autoscaling:EC2_INSTANCE_TERMINATE",
  ]

  topic_arn = aws_sns_topic.topic_example.arn
}
`, rName, rName, rName)
}

func testAccNotificationConfig_update(rName string) string {
	return acctest.ConfigLatestAmazonLinuxHVMEBSAMI() + fmt.Sprintf(`
resource "aws_sns_topic" "topic_example" {
  name = "user-updates-topic-%s"
}

resource "aws_launch_configuration" "foobar" {
  name          = "foobarautoscaling-terraform-test-%s"
  image_id      = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_autoscaling_group" "bar" {
  availability_zones        = [data.aws_availability_zones.available.names[1]]
  name                      = "foobar1-terraform-test-%s"
  max_size                  = 1
  min_size                  = 1
  health_check_grace_period = 100
  health_check_type         = "ELB"
  desired_capacity          = 1
  force_delete              = true
  termination_policies      = ["OldestInstance"]
  launch_configuration      = aws_launch_configuration.foobar.name
}

resource "aws_autoscaling_group" "foo" {
  availability_zones        = [data.aws_availability_zones.available.names[2]]
  name                      = "barfoo-terraform-test-%s"
  max_size                  = 1
  min_size                  = 1
  health_check_grace_period = 200
  health_check_type         = "ELB"
  desired_capacity          = 1
  force_delete              = true
  termination_policies      = ["OldestInstance"]
  launch_configuration      = aws_launch_configuration.foobar.name
}

resource "aws_autoscaling_notification" "example" {
  group_names = [
    aws_autoscaling_group.bar.name,
    aws_autoscaling_group.foo.name,
  ]

  notifications = [
    "autoscaling:EC2_INSTANCE_LAUNCH",
    "autoscaling:EC2_INSTANCE_TERMINATE",
    "autoscaling:EC2_INSTANCE_LAUNCH_ERROR",
  ]

  topic_arn = aws_sns_topic.topic_example.arn
}
`, rName, rName, rName, rName)
}

func testAccNotificationConfig_pagination() string {
	return acctest.ConfigCompose(acctest.ConfigLatestAmazonLinuxHVMEBSAMI(), `
resource "aws_sns_topic" "user_updates" {
  name = "user-updates-topic"
}

resource "aws_launch_configuration" "foobar" {
  image_id      = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
}

data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_autoscaling_group" "bar" {
  availability_zones        = [data.aws_availability_zones.available.names[1]]
  count                     = 20
  name                      = "foobar3-terraform-test-${count.index}"
  max_size                  = 1
  min_size                  = 0
  health_check_grace_period = 300
  health_check_type         = "ELB"
  desired_capacity          = 0
  force_delete              = true
  termination_policies      = ["OldestInstance"]
  launch_configuration      = aws_launch_configuration.foobar.name
}

resource "aws_autoscaling_notification" "example" {
  group_names = aws_autoscaling_group.bar[*].name

  notifications = [
    "autoscaling:EC2_INSTANCE_LAUNCH",
    "autoscaling:EC2_INSTANCE_TERMINATE",
    "autoscaling:TEST_NOTIFICATION"
  ]
  topic_arn = aws_sns_topic.user_updates.arn
}`)
}
