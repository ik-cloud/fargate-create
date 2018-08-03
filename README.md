fargate-create
==============

A CLI tool for scaffolding out new [AWS ECS/Fargate](https://aws.amazon.com/fargate/) applications based on [terraform-ecs-fargate](https://github.com/turnerlabs/terraform-ecs-fargate) and [Fargate CLI](https://github.com/turnerlabs/fargate).


### Why?

The main design goal of this project is to create an easy and maintainable Fargate experience by separating `infrastructure` related concerns and `application` related concerns using tools that are optimized for each.


### Usage

Assuming you have a project with a [Dockerfile]()...

Specify your template's input parameters in [terraform.tfvars](https://www.terraform.io/docs/configuration/variables.html) (or terraform.json).  The [default web application template's](https://github.com/turnerlabs/terraform-ecs-fargate) input looks something like this.

```hcl
region = "us-east-1"
aws_profile = "default"
app = "my-app"
environment = "dev"
internal = "true"
container_name = "app"
container_port = "8080"
lb_port = "80"
lb_protocol = "HTTP"
replicas = "1"
health_check = "/health"
health_check_interval = "10"
health_check_timeout = "5"
health_check_matcher = "200-299"
vpc = "xyz"
private_subnets = "xyz,abc"
public_subnets = "def,ghi"
saml_role = "devops"
tags = {
  app = "my-app"
  env = "dev"
}
```

```shell
$ fargate-create
scaffolding my-app dev
Looking up AWS Account ID using profile: default
downloading terraform template https://github.com/turnerlabs/terraform-ecs-fargate/archive/v0.2.0.zip
installing terraform template

done
```

Now you have all the files you need to spin up something in Fargate.

Infrastructure:  provision using Terraform
```shell
cd iac/base
terraform init && terraform apply
cd ../env/dev
terraform init && terraform apply
```

Application:  build/push using Docker and deploy using Fargate CLI
```shell
docker-compose build
login=$(aws ecr get-login --no-include-email) && eval "$login"
docker-compose push
fargate service deploy -f docker-compose.yml
```

To scaffold out additional environnments, simply change the `environment` input parameter in `terraform.tfvars` and re-run
```shell
$ fargate-create
scaffolding my-app prod
Looking up AWS Account ID using profile: default
downloading terraform template https://github.com/turnerlabs/terraform-ecs-fargate/archive/v0.2.0.zip
installing terraform template
iac/base already exists, ignoring

done
```

And then bring up the new environment (no need to apply base again since it's shared):
```shell
cd ../prod
terraform init && terraform apply
```

You'll end up with a directory structure that looks something like this:
```
.
|____iac
| |____base
| |____env
| | |____dev
| | |____prod
```


### CI/CD (coming soon)

Using this technique, it's easy to codegen CI/CD pipelines for many popular build tools.  The `build` command support this. For example:

```shell
$ fargate-create build circleciv2
```

### Extensibility

works with any terraform template repo (`-t`) that has:
- `base` and `env/dev` directory structure 
- a `env/dev/main.tf` with an s3 remote state backend
- `app` and `environment` input variables

For example (coming soon):
```shell
$ fargate-create -f my-scheduledtask.tfvars -t https://github.com/example/terraform-scheduledtask/archive/v0.1.0.zip
```