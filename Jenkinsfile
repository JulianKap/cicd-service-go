pipeline {
    agent any

    environment {
        DOCKER_IMAGE = "registry.local:5000/cicd-service-go:jenkins-test"
    }

    stages {
        stage('Build') {
           // agent {
          //      docker {
          //          image 'golang:1.22.1'
         //       }
          //  }
            steps {
                script {
                    sh 'go mod download'
                    sh 'go mod verify'
                    sh 'go build -v -o cicd-service-go ./cmd/app'
                }
            }
        }

        stage('Test') {
            //agent {
             //   docker {
            //        image 'golang:1.22.1'
            //    }
           // }
            steps {
                script {
                    sh 'go test -v ./...'
                }
            }
        }

        stage('Build Docker Image') {
           // agent {
          //      docker {
          //          image 'docker:26.1.1-dind-alpine3.19'
           //         args '-v /var/run/docker.sock:/var/run/docker.sock'
          //      }
          //  }
            steps {
                script {
                    sh 'docker build -t ${DOCKER_IMAGE} .'
                    sh 'docker push ${DOCKER_IMAGE}'
                }
            }
        }

        stage('Deploy') {
          //  agent {
          //      docker {
         //           image 'registry.local:5000/cicd-ansible:latest'
         //       }
         //   }
            steps {
                script {
                    sh 'cd ./ansible/'
                    sh 'ansible-playbook --inventory inventories/hosts-service.ini playbooks/deploy.yml'
                }
            }
        }
    }
}
