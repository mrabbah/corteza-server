pipeline {
    agent {
        kubernetes {
            yaml '''
                apiVersion: v1
                kind: Pod
                metadata:
                  namespace: ci
                  labels:
                    pod-name: corteza-server-builder
                spec:
                  containers:
                    - name: golang
                      image: golang:1.17.13
                      command: ["tail", "-f", "/dev/null"]
                      tty: true
                    - name: node
                      image: node:16.16.0
                      command: ["tail", "-f", "/dev/null"]
                      tty: true
                    - name: mc
                      image: mrabbah/mc:1.1
                      command: ["tail", "-f", "/dev/null"]
                      tty: true 
                    - name: docker
                      image: docker:20.10.21
                      command:
                        - sleep
                      args:
                        - 99d
                      env:
                        - name: DOCKER_HOST
                          value: tcp://docker-daemon:2375
                    - name: docker-daemon
                      image: docker:20.10.21-dind
                      securityContext:
                        privileged: true
                      env:
                        - name: DOCKER_TLS_CERTDIR
                          value: ""             
                '''
                //workspaceVolume persistentVolumeClaimWorkspaceVolume(claimName: 'pvc-workspace', readOnly: false)
                //workspaceVolume dynamicPVC(accessModes: 'ReadWriteOnce', requestsSize: "8Gi")
        }      
    }
   
    environment {
        DOCKERHUB_CREDS = credentials('dockerhub-credentials')
        BRANCH_NAME = "${GIT_BRANCH.split('/').size() > 1 ? GIT_BRANCH.split('/')[1..-1].join('/') : GIT_BRANCH}"
        MINIO_CREDS = credentials('minio-credentials')
        MINIO_HOST = "http://minio.data:9000"
    }
    
    stages {
        /*stage('Test') {
            steps {
                container('golang') {
                    sh 'GOCACHE=/tmp/.cache/go-build GOENV=/tmp/.config/go/env GOMODCACHE=/tmp/go/pkg/mod go test ./pkg/... ./app/... ./compose/... ./system/... ./federation/... ./auth/... ./automation/... ./tests/... ./store/tests/...'
                }              
            }
        }
        stage('Build Web Console') {
            steps {
                container('node') {
                    sh 'cd ./webconsole && yarn install && yarn build && cd ..'
                }                   
            }
        }
        stage('Build Server') {
            steps {
                container('golang') {
                    sh 'rm -rf ./build/corteza-server-${BRANCH_NAME}'
                    sh 'rm -rf ./build/corteza-server-${BRANCH_NAME}.tar.gz'
                    sh 'GOCACHE=/tmp/.cache/go-build GOENV=/tmp/.config/go/env GOMODCACHE=/tmp/go/pkg/mod GOOS=linux GOARCH=amd64 go build  -ldflags "-X github.com/mrabbah/corteza-server/pkg/version.Version=${BRANCH_NAME} " -o ./build/corteza-server-${BRANCH_NAME} cmd/corteza/main.go'
                    sh 'mkdir -p ./build/pkg/corteza-server ./build/pkg/corteza-server/bin'
                    sh 'cp ./README.md ./LICENSE ./CONTRIBUTING.md ./DCO ./.env.example ./build/pkg/corteza-server/'
                    sh 'cp -r ./provision ./build/pkg/corteza-server'
                    sh 'rm -f ./build/pkg/corteza-server/provision/README.adoc ./build/pkg/corteza-server/provision/update.sh'
                    sh 'cp ./build/corteza-server-${BRANCH_NAME} ./build/pkg/corteza-server/bin/corteza-server'
                    sh 'tar -C ./build/pkg/ -czf ./build/corteza-server-${BRANCH_NAME}.tar.gz corteza-server' 
                }  
            }
        }*/
        stage('Pushing Artifact') {
            steps {
                container('mc') {
                    sh 'mc --config-dir /tmp/.mc alias set minio $MINIO_HOST $MINIO_CREDS_USR $MINIO_CREDS_PSW'
                    //sh 'mc --config-dir /tmp/.mc cp ./build/corteza-server-${BRANCH_NAME}.tar.gz minio/corteza-artifacts'             
                }   
            }
        }
        /*stage('Build Docker image') {
            steps {
                container('docker') {
                    sh 'docker build -t mrabbah/corteza-server:${BRANCH_NAME} --build-arg VERSION=${BRANCH_NAME} . '             
                }                
            }
        }
        
        stage('Push Docker image') {
            
            steps {
                container('docker') {
                    echo 'Pushing docker image'
                    script {
                        sh 'echo $DOCKERHUB_CREDS_PSW | docker login -u $DOCKERHUB_CREDS_USR --password-stdin'    
                        sh 'docker push mrabbah/corteza-server:${BRANCH_NAME}'           
                    }         
                }            
                
            }
        }*/

    }
}
