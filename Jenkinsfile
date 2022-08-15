pipeline {
    agent any
   
    environment {
        NEXUS_CREDS = credentials('nexus-credentials')
        BRANCH_NAME = "${GIT_BRANCH.split('/').size() > 1 ? GIT_BRANCH.split('/')[1..-1].join('/') : GIT_BRANCH}"
    }
    stages {
        stage('Test') {
            agent {
                docker { 
                    image 'golang:1.17.13' 
                    reuseNode true
                }            
            }
            steps {
                //sh 'make test.all'
                sh 'GOCACHE=/tmp/.cache/go-build GOENV=/tmp/.config/go/env GOMODCACHE=/tmp/go/pkg/mod go test ./pkg/... ./app/... ./compose/... ./system/... ./federation/... ./auth/... ./automation/... ./tests/... ./store/tests/...'
            }
        }
        stage('Build Web Console') {
            agent {
                docker { 
                    image 'node:16.16.0'
                    reuseNode true 
                }
            }
            steps {
                sh 'cd ./webconsole && yarn install && yarn build'
            }
        }
        stage('Build Server') {
            agent {
                docker { 
                    image 'golang:1.17.13' 
                    reuseNode true
                } 
            }
            steps {
                sh 'ls build/'
                sh 'rm -rf build/corteza-server-${BRANCH_NAME}'
                sh 'rm -rf build/corteza-server-${BRANCH_NAME}.tar.gz'
                sh 'GOCACHE=/tmp/.cache/go-build GOENV=/tmp/.config/go/env GOMODCACHE=/tmp/go/pkg/mod go build -o build/corteza-server-${BRANCH_NAME} cmd/corteza/main.go'
                sh 'ls build/'
                sh 'curl -v --user $NEXUS_CREDS --upload-file ./build/corteza-server-${BRANCH_NAME}.tar.gz https://nexus.rabbahsoft.ma/repository/row-repo/corteza-server-${BRANCH_NAME}.tar.gz'
            }
        }
        stage('Build Docker image') {
            
            steps {
                echo 'Starting to build docker image'

                script {
                    def cortezaServerImage = docker.build("mrabbah/corteza-server:${BRANCH_NAME}", "VERSION=${BRANCH_NAME} CORTEZA_VERSION=2022.3.4 NEXUS_CREDS=${NEXUS_CREDS}")
                    cortezaServerImage.push()
                }
                
            }
        }

    }
}
