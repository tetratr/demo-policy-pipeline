node {
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
            
            docker.withRegistry('', "docker-hub-remiphilippe") {
                dir('src/cmd/project/') {
                    sh "git rev-parse --short HEAD > .git/commit-id"
                    def commit_id = readFile(".git/commit-id").trim()
                    println commit_id

                    stage('build and push image') {
                        def app = docker.build "remiphilippe/demo-policy-pipeline"
                        app.push('latest')
                        app.push("${commit_id}")
                    }

                    stage('Remove Unused docker image') {
                        sh "docker rmi remiphilippe/demo-policy-pipeline:${commit_id}"
                    }
                }
            }

            stage("deploy kubernetes") {
                if (fileExists ("$GOPATH/src/cmd/project/deploy-dev.yml")) {
                    echo "Deploying Policies"
                    sh """kubectl apply -f $GOPATH/src/cmd/project/deploy-dev.yml"""
                }
            }

            stage("deploy policies") {
                if (fileExists ("$GOPATH/src/cmd/project/policy-dev.yml")) {
                    echo "Deploying Policies"
                    sh """/usr/local/bin/kubepol -file=$GOPATH/src/cmd/project/policy-dev.yml"""
                }
            }
        }
    }
}
