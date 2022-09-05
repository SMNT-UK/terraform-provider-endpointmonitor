pipeline {
    agent any

    envirionment {
        EPM_API_KEY = credentials('terraform-epm-api-key')
        AWS_REGION="eu-central-1"
    }

    stages {
        stage("Build") {
            agent {
                docker {
                    image "golang:1.18-alpine"
                }
            }

            steps {
                wsCleanup()

                sh "go mod tidy"
                sh "go build -o terraform-provider-endpointmonitor"

                stash "terraform-provider-endpointmonitor"
            }
        }

        stage("Test") {
            agent {
                docker {
                    image "hashicorp/terraform"
                    args "--entrypoint=''"
                }
            }

            steps {
                dir("tests/integration") {
                    unstash "terraform-provider-endpointmonitor"

                    withCredentials([file(credentialsId: 'reading-internal-ca-cert', variable: 'FILE')]) {
                        sh "cp $FILE /usr/local/share/ca-certificates/"
                        sh "update-ca-certificates"
                    }

                    withCredentials([[$class: 'AmazonWebServicesCredentialsBinding',
                            credentialsId: 'aws-jenkins',
                            accessKeyVariable: 'AWS_ACCESS_KEY_ID',
                            secretKeyVariable: 'AWS_SECRET_ACCESS_KEY']]) {
                        sh "mkdir -p ~/.terraform.d/registry.terraform.io/smnt/endpointmonitor"
                        sh "mv ./terraform-provider-endpointmonitor ~/.terraform.d/registry.terraform.io/smnt/endpointmonitor/"
                        sh "terraform init"
                        sh "terraform plan -out terraform.plan"
                        sh "terraform apply terraform.plan"
                    }
                }
            }
        }
    }
}