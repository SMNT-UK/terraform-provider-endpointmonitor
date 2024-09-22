pipeline {
    agent any

    environment {
        EPM_API_KEY = credentials('terraform-epm-api-key')
        AWS_REGION="eu-central-1"
    }

    stages {
        stage("Build") {
            agent {
                docker {
                    image "golang:1.22-alpine"
                    args "-u 0:0"
                }
            }

            steps {
                sh "go mod tidy"
                sh "go build -o terraform-provider-endpointmonitor"

                stash "terraform-provider-endpointmonitor"
            }
        }

        stage("Test") {
            agent {
                docker {
                    image "hashicorp/terraform"
                    args "-u 0:0 --entrypoint=''"
                }
            }

            steps {
                dir("tests/integration") {
                    unstash "terraform-provider-endpointmonitor"

                    withCredentials([file(credentialsId: 'reading-internal-ca-cert', variable: 'FILE')]) {
                        sh "cp $FILE /usr/local/share/ca-certificates/"
                        sh "update-ca-certificates"
                    }

                    sh "rm -rf .terraform"
                    sh "rm -rf .terraform.d"
                    sh "rm -f .terraform.lock.hcl"

                    withCredentials([[$class: 'AmazonWebServicesCredentialsBinding',
                            credentialsId: 'aws-jenkins',
                            accessKeyVariable: 'AWS_ACCESS_KEY_ID',
                            secretKeyVariable: 'AWS_SECRET_ACCESS_KEY']]) {
                        sh "mkdir -p ~/.terraform.d/plugins/registry.terraform.io/smnt/endpointmonitor/0.1/linux_amd64"
                        sh "mv ./terraform-provider-endpointmonitor ~/.terraform.d/plugins/registry.terraform.io/smnt/endpointmonitor/0.1/linux_amd64/"
                        sh "terraform init"
                        sh "terraform plan -out terraform.plan"
                        sh "terraform apply terraform.plan"
                    }
                }
            }
        }
    }
    post {
        always {
            cleanWs()
        }
    }
}