node {
    environment {
        registry = "remiphilipppe/demo-policy-pipeline"
        registryCredential = 'dockerhub'
    }

    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
        withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
            env.PATH="${GOPATH}/bin:$PATH"
            
            stage('Checkout'){
                echo 'Checking out SCM'
                sh 'mkdir -p src/cmd/project/'
                dir('src/cmd/project/') {
                    checkout scm
                }
            }
            
            stage('Pre Test'){
                echo 'Pulling Dependencies'
        
                sh 'go version'
                sh 'go get -u github.com/golang/dep/cmd/dep'
                sh 'go get -u golang.org/x/lint/golint'
                sh 'go get github.com/tebeka/go2xunit'
                
                //or -update
                sh 'cd ${GOPATH}/src/cmd/project/ && dep ensure' 
            }
    
            // stage('Test'){
                
            //     //List all our project files with 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'
            //     //Push our project files relative to ./src
            //     sh 'cd $GOPATH && go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org > projectPaths'
                
            //     //Print them with 'awk '$0="./src/"$0' projectPaths' in order to get full relative path to $GOPATH
            //     def paths = sh returnStdout: true, script: """awk '\$0="./src/"\$0' projectPaths"""
                
            //     echo 'Vetting'

            //     sh """cd $GOPATH && go vet ${paths}"""

            //     echo 'Linting'
            //     sh """cd $GOPATH && golint ${paths}"""
                
            //     echo 'Testing'
            //     sh """cd $GOPATH && go test -race -cover ${paths}"""
            // }
        
            stage('Build'){
                echo 'Building Executable'
            
                //Produced binary is $GOPATH/src/cmd/project/project
                sh """cd $GOPATH/src/cmd/project/ && make build"""
            }
            
            stage('Building image') {
                dockerImage = docker.build registry + ":${BUILD_NUMBER}"
            }

            docker.withRegistry('', registryCredential) {
                sh "git rev-parse HEAD > .git/commit-id"
                def commit_id = readFile('.git/commit-id').trim()
                println commit_id

                stage('build image') {
                    def app = docker.build registry
                }
                

                stage("publish image") {
                    app.push 'master'
                    app.push "${commit_id}"
                }

                stage('Remove Unused docker image') {
                    sh "docker rmi ${registry}:${commit_id}"
                }
            }
        }
    }
}
