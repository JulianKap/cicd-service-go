pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                container('golang:1.22.1') {
                    sh 'go mod download'
                    sh 'go mod verify'
                    sh 'go build -v -o cicd-service-go ./cmd/app'
                }
            }
        }

        stage('Test') {
            steps {
                container('golang:1.22.1') {
                    sh 'go test -v ./...'
                }
            }
        }

        stage('Build docker image') {
            steps {
                container('docker:latest') {
                    sh 'docker build -t cicd-service-go:latest .'
                    sh 'docker tag cicd-service-go:latest registry:5000/cicd-service-go:latest'
                    sh 'docker push registry:5000/cicd-service-go:latest'
                }
            }
        }

        stage('Deploy') {
            steps {
                sh './deploy.sh'
            }
        }
    }
}
